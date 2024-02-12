package codegen

import (
	"bytes"
	"fmt"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/pterm/pterm"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func GenerateProto(pb *PbPackage) error {
	goPbFilePath := filepath.Join(filepath.Dir(filepath.Dir(pb.protoFilePath)), pb.GoPackage, pb.PackageName+".pb.go")

	areas, err := InjectTagParseFile(goPbFilePath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if len(areas) > 0 {
		err = InjectTagWriteFile(goPbFilePath, areas)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	pterm.Success.Printfln("generate pb file %s", goPbFilePath)

	return nil
}

type tagItem struct {
	key   string
	value string
}

type tagItems []tagItem

func (ti tagItems) format() string {
	var tags []string
	for _, item := range ti {
		tags = append(tags, fmt.Sprintf(`%s:"%s"`, item.key, item.value))
	}
	return strings.Join(tags, " ")
}

func (ti tagItems) writeTo(b *bytes.Buffer) {
	for _, item := range ti {
		b.WriteString(item.key)
		b.WriteString(`:"`)
		b.WriteString(item.value)
		b.WriteString(`"`)
		b.WriteString(" ")
	}

	b.Truncate(b.Len() - 1)
}

func (ti tagItems) override() tagItems {
	m := make(map[string][]string, len(ti))
	for _, item := range ti {
		if _, ok := m[item.key]; !ok {
			m[item.key] = []string{}
		}

		m[item.key] = append(m[item.key], item.value)
	}

	//for _, item := range nti {
	//	if _, ok := m[item.key]; !ok {
	//		m[item.key] = []string{}
	//	}
	//
	//	m[item.key] = append(m[item.key], item.value)
	//}

	var overrided tagItems
	for k, v := range m {
		v = candy.Filter(v, func(s string) bool {
			return s != ""
		})

		switch k {
		case "gorm":
			overrided = append(overrided, tagItem{
				key:   k,
				value: strings.Join(v, ";"),
			})
		default:
			overrided = append(overrided, tagItem{
				key:   k,
				value: strings.Join(v, ","),
			})
		}
	}

	return overrided
}

func injectTag(contents []byte, area textArea) (injected []byte) {
	var b bytes.Buffer
	b.Write(contents[area.Start-1 : area.End-1])
	//b.Truncate(bytes.LastIndex(contents[area.Start-1:area.End-1], []byte("`")))
	b.Truncate(area.End - area.Start - 1)

	pterm.Info.Printfln("append custom tag %s to expression %s",
		pterm.BgYellow.Sprint(b.String()), pterm.BgGreen.Sprint(area.InjectTag.format()))

	b.WriteString(" ")
	area.InjectTag.writeTo(&b)
	b.WriteString("`")

	injected = append(injected, contents[:area.Start-1]...)
	injected = append(injected, b.Bytes()...)
	injected = append(injected, contents[area.End-1:]...)

	return
}

func InjectTagWriteFile(inputPath string, areas []textArea) error {
	stat, err := os.Stat(inputPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	contents, err := os.ReadFile(inputPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	areas = candy.SortUsing(areas, func(a, b textArea) bool {
		return a.Start > b.Start
	})

	for _, area := range areas {
		contents = injectTag(contents, area)
	}

	err = os.WriteFile(inputPath, contents, stat.Mode())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func InjectTagParseFile(inputPath string) (areas []textArea, err error) {
	f, err := parser.ParseFile(token.NewFileSet(), inputPath, nil, parser.ParseComments)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		var typeSpec *ast.TypeSpec
		for _, spec := range genDecl.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok {
				typeSpec = ts
				break
			}
		}

		// skip if can't get type spec
		if typeSpec == nil {
			continue
		}

		// not a struct, skip
		structDecl, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		for _, field := range structDecl.Fields.List {
			if field.Doc == nil {
				continue
			}

			var injectTags tagItems
			var lastTag string
			for _, v := range field.Doc.List {
				line := strings.ReplaceAll(strings.ReplaceAll(strings.TrimPrefix(v.Text, "//"), " ", ""), "\t", "")
				if !strings.HasPrefix(line, "@") {
					if lastTag != "" {
						injectTags = append(injectTags, tagItem{
							key:   lastTag,
							value: removeStrQuote(line),
						})
					}
					continue
				} else {
					idx := strings.Index(line, ":")
					if idx < 0 {
						continue
					}

					tag := line[1:idx]
					injectTags = append(injectTags, tagItem{
						key:   tag,
						value: removeStrQuote(line[idx+1:]),
					})
					lastTag = tag
				}
			}
			if len(injectTags) == 0 {
				continue
			}

			currentTag := field.Tag.Value
			currentTag = strings.TrimSuffix(strings.TrimPrefix(currentTag, "`"), "`")
			area := textArea{
				Start:     int(field.Pos()),
				End:       int(field.End()),
				InjectTag: injectTags.override(),
			}
			areas = append(areas, area)
		}
	}

	log.Infof("number of fields to inject custom tags: %d", len(areas))
	return
}

type textArea struct {
	Start     int
	End       int
	InjectTag tagItems
}

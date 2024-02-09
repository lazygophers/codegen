package codegen

import (
	"fmt"
	"github.com/lazygophers/log"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
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
		tags = append(tags, fmt.Sprintf(`%s:%s`, item.key, item.value))
	}
	return strings.Join(tags, " ")
}

func (ti tagItems) override(nti tagItems) tagItems {
	var overrided []tagItem
	for i := range ti {
		var dup = -1
		for j := range nti {
			if ti[i].key == nti[j].key {
				dup = j
				break
			}
		}
		if dup == -1 {
			overrided = append(overrided, ti[i])
		} else {
			overrided = append(overrided, nti[dup])
			nti = append(nti[:dup], nti[dup+1:]...)
		}
	}
	return append(overrided, nti...)
}

func newTagItems(tag string) tagItems {
	var items []tagItem
	splitted := rTags.FindAllString(tag, -1)

	for _, t := range splitted {
		sepPos := strings.Index(t, ":")
		items = append(items, tagItem{
			key:   t[:sepPos],
			value: t[sepPos+1:],
		})
	}
	return items
}

func injectTag(contents []byte, area textArea) (injected []byte) {
	expr := make([]byte, area.End-area.Start)
	copy(expr, contents[area.Start-1:area.End-1])

	cti := area.CurrentTag
	ti := cti.override(area.InjectTag)

	expr = rInject.ReplaceAll(expr, []byte(fmt.Sprintf("`%s`", ti.format())))

	injected = append(injected, contents[:area.Start-1]...)
	injected = append(injected, expr...)
	injected = append(injected, contents[area.End-1:]...)
	
	return
}

func InjectTagWriteFile(inputPath string, areas []textArea) (err error) {
	contents, err := os.ReadFile(inputPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// inject custom tags from tail of file first to preserve order
	for i := range areas {
		area := areas[len(areas)-i-1]
		log.Infof("inject custom tag %q to expression %q", area.InjectTag, string(contents[area.Start-1:area.End-1]))
		contents = injectTag(contents, area)
	}

	err = os.WriteFile(inputPath, contents, 0644)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return
}

func InjectTagParseFile(inputPath string) (areas []textArea, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inputPath, nil, parser.ParseComments)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	for _, decl := range f.Decls {
		// check if is generic declaration
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		var typeSpec *ast.TypeSpec
		for _, spec := range genDecl.Specs {
			if ts, tsOK := spec.(*ast.TypeSpec); tsOK {
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

			var commentLines []string
			for _, v := range field.Doc.List {
				commentLines = append(commentLines, v.Text)
			}

			injectTags := injectTagCommentParse(commentLines)
			if len(injectTags) == 0 {
				continue
			}

			currentTag := field.Tag.Value
			area := textArea{
				Start:      int(field.Pos()),
				End:        int(field.End()),
				CurrentTag: newTagItems(currentTag[1 : len(currentTag)-1]),
				InjectTag:  injectTags,
			}
			areas = append(areas, area)
		}
	}

	log.Infof("number of fields to inject custom tags: %d", len(areas))
	return
}

type textArea struct {
	Start      int
	End        int
	CurrentTag tagItems
	InjectTag  tagItems
}

var (
	rInject = regexp.MustCompile("`.+`$")
	rTags   = regexp.MustCompile(`[\w_]+:"[^"]+"`)
)

func injectTagCommentParse(lines []string) tagItems {
	var resList tagItems
	for _, line := range lines {
		line = strings.ReplaceAll(line, " ", "")
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "//@") {
			continue
		}

		line = line[3:]

		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}

		tag := line[:idx]
		value := line[idx+1:]

		resList = append(resList, tagItem{
			key:   tag,
			value: removeStrQuote(value),
		})
	}

	return resList
}

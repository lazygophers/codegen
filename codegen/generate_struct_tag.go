package codegen

import (
	"bytes"
	"github.com/lazygophers/codegen/state"
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
	goPbFilePath := filepath.Join(pb.ProjectRoot(), pb.PackageName+".pb.go")

	areas, err := InjectTagParseFile(goPbFilePath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = InjectTagWriteFile(goPbFilePath, areas)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	pterm.Success.Printfln("generate pb file %s", goPbFilePath)

	return nil
}

type tagItem struct {
	key   string
	value string
}

type tagItems []tagItem

func (p tagItems) writeTo(b *bytes.Buffer) {
	for _, item := range p {
		b.WriteString(item.key)
		b.WriteString(`:"`)
		b.WriteString(item.value)
		b.WriteString(`"`)
		b.WriteString(" ")
	}

	b.Truncate(b.Len() - 1)
}

func (p tagItems) override() tagItems {
	m := make(map[string][]string, len(p))
	for _, item := range p {
		if _, ok := m[item.key]; !ok {
			m[item.key] = []string{}
		}

		m[item.key] = append(m[item.key], item.value)
	}

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

func (p tagItems) get(key string) string {
	for _, item := range p {
		if item.key == key {
			return item.value
		}
	}

	return ""
}

func InjectTagWriteFile(inputPath string, areas []textArea) error {
	if len(areas) == 0 {
		return nil
	}

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
		return a.Start < b.Start
	})

	var b bytes.Buffer
	var lastEnd int
	for _, area := range areas {
		endIdx := bytes.LastIndex(contents[area.Start-1:area.End-1], []byte("`")) + area.Start - 1

		log.Infof("append custom tags to %s at %d", contents[area.Start-1:endIdx], endIdx)
		pterm.Info.Printfln("append custom tags to %s at %s", pterm.FgWhite.Sprint(pterm.BgCyan.Sprintf("%s", contents[area.Start-1:endIdx])), pterm.FgMagenta.Sprint(endIdx))

		b.Write(contents[lastEnd:endIdx])
		b.WriteString(" ")
		area.InjectTag.writeTo(&b)
		lastEnd = endIdx
	}

	b.Write(contents[lastEnd:])

	err = os.WriteFile(inputPath, b.Bytes(), stat.Mode())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func gormTagStr2Map(s string) map[string]string {
	m := make(map[string]string)
	for _, v := range strings.Split(s, ";") {
		idx := strings.Index(v, ":")
		if idx < 0 {
			m[v] = ""
		} else {
			m[v[:idx]] = v[idx+1:]
		}
	}

	return m
}

func InjectTagParseFile(inputPath string) ([]textArea, error) {
	f, err := parser.ParseFile(token.NewFileSet(), inputPath, nil, parser.ParseComments)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	var areas []textArea
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

		isModelStruct := strings.HasPrefix(typeSpec.Name.Name, "Model") && strings.Contains(typeSpec.Name.Name, "_")

		log.Infof("find type %s", typeSpec.Name.Name)

		for _, field := range structDecl.Fields.List {
			var injectTags tagItems

			// 私有字段不处理
			if len(field.Names) == 0 {
				continue
			}

			if !ast.IsExported(field.Names[0].Name) {
				continue
			}

			if field.Doc != nil {
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
			}

			log.Info(field.Names[0].Name)

			var fieldName string
			if x, ok := field.Names[0].Obj.Decl.(*ast.Field); ok {
				if x.Tag == nil {
					log.Warnf("no tag for %s", field.Names[0].Name)
					continue
				}

				idx := strings.Index(x.Tag.Value, `protobuf:"`)
				if idx < 0 {
					log.Warnf("unknown tag %s", x.Tag.Value)
					continue
				}

				fieldName = x.Tag.Value[idx+10:]
				idx = strings.Index(fieldName, `name=`)
				if idx < 0 {
					log.Warnf("unknown tag %s", x.Tag.Value)
					continue
				}
				fieldName = fieldName[idx+5:]

				idx = strings.Index(fieldName, `,`)
				if idx < 0 {
					log.Warnf("unknown tag %s", x.Tag.Value)
					continue
				}

				fieldName = fieldName[:idx]

			} else {
				log.Warnf("unknown type %T", x)
				continue
			}

			injectTags = injectTags.override()

			addGorm := func(m map[string]string, values string) {
				for _, value := range strings.Split(values, ";") {
					idx := strings.Index(value, ":")
					if idx < 0 {
						delete(m, value)
					} else {
						delete(m, value[:idx])
					}
				}

				for k, v := range m {
					if v == "" {
						injectTags = append(injectTags, tagItem{
							key:   "gorm",
							value: k,
						})
					} else {
						injectTags = append(injectTags, tagItem{
							key:   "gorm",
							value: k + ":" + v,
						})
					}
				}
			}

			if isModelStruct {
				if len(field.Names) > 0 {
					switch field.Names[0].Name {
					case "Id":
						if state.Config.Tables.DisableFieldId {
							break
						}

						// 添加gorm标签
						addGorm(gormTagStr2Map(state.Config.Tables.DefaultGormFieldId), injectTags.get("gorm"))
					case "CreatedAt":
						if state.Config.Tables.DisableFieldCreatedAt {
							break
						}

						// 添加gorm标签
						addGorm(gormTagStr2Map(state.Config.Tables.DefaultGormFieldCreatedAt), injectTags.get("gorm"))
					case "UpdatedAt":
						if state.Config.Tables.DisableFieldUpdatedAt {
							break
						}

						// 添加gorm标签
						addGorm(gormTagStr2Map(state.Config.Tables.DefaultGormFieldUpdatedAt), injectTags.get("gorm"))

					case "DeletedAt":
						if state.Config.Tables.DisableFieldDeletedAt {
							break
						}

						// 添加gorm标签
						addGorm(gormTagStr2Map(state.Config.Tables.DefaultGormFieldDeletedAt), injectTags.get("gorm"))

					}
				}

				injectTags = injectTags.override()
			}

			// TODO: 对 Table 的子类型的字段进行处理，暂时先按照如果如要就全体都加
			if !state.Config.Tables.DisableGormTagColumn {
				addGorm(map[string]string{
					"column": fieldName,
				}, injectTags.get("gorm"))

				injectTags = injectTags.override()
			}

			if len(injectTags) == 0 {
				continue
			}

			areas = append(areas, textArea{
				Start:     int(field.Pos()),
				End:       int(field.End()),
				InjectTag: injectTags,
			})
		}
	}

	log.Infof("number of fields to inject custom tags: %d", len(areas))

	return areas, nil
}

type textArea struct {
	Start     int
	End       int
	InjectTag tagItems
}

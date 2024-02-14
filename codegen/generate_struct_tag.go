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
		pterm.Info.Printfln("append custom tags to %s at %s", pterm.BgMagenta.Sprintf("%s", contents[area.Start-1:endIdx]), pterm.FgMagenta.Sprint(endIdx))

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

		isModelStruct := strings.HasPrefix(typeSpec.Name.Name, "Model")

		for _, field := range structDecl.Fields.List {
			var injectTags tagItems

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

			injectTags = injectTags.override()

			//Id        uint64 `gorm:"column:id;primaryKey;autoIncrement;not null" json:"id,omitempty" role:"admin"`
			//CreatedAt int64  `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
			//UpdatedAt int64  `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null" role:"admin"`
			//
			//DeleteAt int64 `json:"delete_at,omitempty" gorm:"column:delete_at;type:bigint(20);not null;default:0;index:idx_username,unique" role:"admin"`
			if isModelStruct {
				if len(field.Names) > 0 {
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

					switch field.Names[0].Name {
					case "Id":
						if state.Config.Tables.DisableAutoId {
							break
						}

						// 添加gorm标签
						addGorm(map[string]string{
							"column":        "id",
							"primaryKey":    "",
							"autoIncrement": "",
							"not null":      "",
						}, injectTags.get("gorm"))
					case "CreatedAt":
						if state.Config.Tables.DisableAutoCreatedAt {
							break
						}

						// 添加gorm标签
						addGorm(map[string]string{
							"autoCreateTime": "",
							"<-":             "create",
							"column":         "created_at",
							"not null":       "",
						}, injectTags.get("gorm"))

					case "UpdatedAt":
						if state.Config.Tables.DisableAutoUpdatedAt {
							break
						}

						// 添加gorm标签
						addGorm(map[string]string{
							"autoUpdateTime": "",
							"<-":             "",
							"column":         "updated_at",
							"not null":       "",
						}, injectTags.get("gorm"))

					case "DeletedAt":
						if state.Config.Tables.DisableAutoDeletedAt {
							break
						}

						// 添加gorm标签
						addGorm(map[string]string{
							"column":   "deleted_at",
							"not null": "",
						}, injectTags.get("gorm"))

					}
				}

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
	return
}

type textArea struct {
	Start     int
	End       int
	InjectTag tagItems
}

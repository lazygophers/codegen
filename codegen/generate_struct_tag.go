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
	"strings"
)

func GenerateProto(pb *PbPackage) error {
	goPbFilePath := GetPath(PathTypePbGo, pb)

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
		v = candy.Unique(v)

		// 校验如果存在 - ，则跳过
		if candy.Contains(v, "-") {
			overrided = append(overrided, tagItem{
				key:   k,
				value: "-",
			})
			continue
		}

		switch k {
		case "gorm":
			if candy.Contains(v, "autoIncrement") {
				// 自动的，不添加 default
				v = candy.FilterNot(v, func(s string) bool {
					return strings.Contains(s, "default:")
				})
			}

			// 解析一下是否存在 type
			for _, line := range v {
				if !strings.HasPrefix(line, "type:") {
					continue
				}

				// 针对特定类型，去掉默认值
				switch strings.ToLower(strings.TrimPrefix(line, "type:")) {
				case "text", "blob", "geometry", "json":
					v = candy.FilterNot(v, func(s string) bool {
						return strings.Contains(s, "default:")
					})

				}

				break
			}

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
		pterm.Info.Printfln("append custom tags to %s at %s", pterm.FgBlack.Sprint(pterm.BgWhite.Sprintf("%s", contents[area.Start-1:endIdx])), pterm.FgMagenta.Sprint(endIdx))

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

func gormTagStr2Map(items []string) map[string]string {
	m := make(map[string]string)

	for _, item := range items {
		for _, v := range strings.Split(item, ";") {
			idx := strings.Index(v, ":")
			if idx < 0 {
				m[v] = ""
			} else {
				m[v[:idx]] = v[idx+1:]
			}
		}
	}

	return m
}

func tagStr2Map(items []string) map[string]string {
	m := make(map[string]string)
	for _, item := range items {
		for _, v := range strings.Split(item, ",") {
			idx := strings.Index(v, "=")
			if idx < 0 {
				m[v] = ""
			} else {
				m[v[:idx]] = v[idx+1:]
			}
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

		//isModelStruct := strings.HasPrefix(typeSpec.Name.Name, "Model") && !strings.Contains(typeSpec.Name.Name, "_")

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

						// NOTE: 插入一段特殊逻辑，后面在加别名配置
						if tag == "v" {
							tag = "validate"
						}

						injectTags = append(injectTags, tagItem{
							key:   tag,
							value: removeStrQuote(line[idx+1:]),
						})
						lastTag = tag
					}
				}
			}

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

			addTag := func(key string, m map[string]string, values string) {
				seq := ","
				connect := "="
				switch key {
				case "gorm":
					seq = ";"
					connect = ":"
				}

				for _, value := range strings.Split(values, seq) {
					before, _, found := strings.Cut(value, connect)
					if !found {
						delete(m, value)
					} else {
						delete(m, before)
					}
				}

				for k, v := range m {
					if v == "" {
						injectTags = append(injectTags, tagItem{
							key:   key,
							value: k,
						})
					} else {
						injectTags = append(injectTags, tagItem{
							key:   key,
							value: k + connect + v,
						})
					}
				}
			}

			// 先按照类型获取一下
			var getFieldType func(xx ast.Expr) string
			var getObjType func(xx *ast.Object) string

			getObjType = func(xx *ast.Object) string {
				if xx == nil {
					return ""
				}

				if xx.Decl != nil {
					switch x := xx.Decl.(type) {
					case *ast.TypeSpec:
						return getFieldType(x.Type)

					default:
						log.Panicf("unknown type %T", x)
					}
				}

				return ""
			}

			getFieldType = func(xx ast.Expr) (name string) {
				switch x := xx.(type) {
				case *ast.Ident:
					name = getObjType(x.Obj)
					if name != "" {
						return name
					}

					return x.Name

				case *ast.StarExpr:
					name = getFieldType(x.X)
					if name == "" {
						return ""
					}

					return "*" + name

				case *ast.ArrayType:
					return "array"

				case *ast.MapType:
					return "map"

				case *ast.SelectorExpr:
					// TODO
					return ""

				case *ast.StructType:
					return "object"

				default:
					log.Panicf("unknown type %T", x)
				}

				return ""
			}

			fieldType := getFieldType(field.Type)

			log.Infof("%s field type: %s", fieldName, fieldType)

			tagsMap := map[string][]string{
				"yaml": {
					CoverageStyledBase(state.Config.Style.Yaml, fieldName),
					"omitempty",
				},
			}

			// 按照字段名的 gorm 默认
			if tm := state.Config.DefaultTag[fieldName]; tm != nil {
				for k, v := range tm {
					tagsMap[k] = append(tagsMap[k], v)
				}
			}

			// 按照字段类型的 gorm 默认
			switch fieldType {
			case "":
				fallthrough
			case "int32", "int64", "uint32", "uint64", "sint32", "sint64":
				fallthrough
			case "float", "double", "float32", "float64":
				fallthrough
			case "string", "bytes":
				fallthrough
			case "bool":
				fallthrough
			case "array":
				fallthrough
			case "map":
				fallthrough
			case "object":
				if tm := state.Config.DefaultTag["@"+fieldType]; tm != nil {
					for k, v := range tm {
						tagsMap[k] = append(tagsMap[k], v)
					}
				}

			default:
				if tm := state.Config.DefaultTag["@object"]; tm != nil {
					for k, v := range tm {
						tagsMap[k] = append(tagsMap[k], v)
					}
				}
			}

			if !state.Config.Tables.DisableGormTagColumn {
				tagsMap["gorm"] = append(tagsMap["gorm"], "column:"+fieldName)
			}

			// 处理用户填写的
			for key, value := range tagsMap {
				if key == "gorm" {
					addTag("gorm", gormTagStr2Map(value), injectTags.get("gorm"))
				} else {
					addTag(key, tagStr2Map(value), injectTags.get(key))
				}
			}

			injectTags = injectTags.override()

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

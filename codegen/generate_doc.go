package codegen

import (
	"os"
	"path/filepath"
	"time"

	"strings"

	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/i18n"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/osx"
	"github.com/lazygophers/utils/stringx"
	"github.com/pterm/pterm"
)

// DocField 文档字段结构
type DocField struct {
	Name        string `json:"name"`        // 字段名
	Type        string `json:"type"`        // 字段类型
	Column      string `json:"column"`      // 数据库列名
	Description string `json:"description"` // 字段描述
	IsPrimary   bool   `json:"is_primary"`  // 是否为主键
}

// DocModel 文档模型结构
type DocModel struct {
	Name        string      `json:"name"`        // 模型名
	TableName   string      `json:"table_name"`  // 表名
	Description string      `json:"description"` // 模型描述
	Fields      []*DocField `json:"fields"`      // 字段列表
	Indexes     []*DocIndex `json:"indexes"`     // 索引列表
}

// DocIndex 文档索引结构
type DocIndex struct {
	Name      string   `json:"name"`       // 索引名
	Type      string   `json:"type"`       // 索引类型：primary, unique, index
	Fields    []string `json:"fields"`     // 索引字段名列表
	Columns   []string `json:"columns"`    // 索引列名列表
	IsUnique  bool     `json:"is_unique"`  // 是否为唯一索引
	IsPrimary bool     `json:"is_primary"` // 是否为主键
}

// DocData 文档数据结构
type DocData struct {
	PackageName   string     `json:"package_name"`    // 包名
	ProtoFile     string     `json:"proto_file"`      // Proto文件名
	GoPackage     string     `json:"go_package"`      // Go包路径
	GoPackageName string     `json:"go_package_name"` // Go包名
	GeneratedAt   string     `json:"generated_at"`    // 生成时间
	Models        []DocModel `json:"models"`          // 模型列表
}

// convertPbMessageToDocModel 将 PbMessage 转换为 DocModel
func convertPbMessageToDocModel(pbMsg *PbMessage) DocModel {
	model := DocModel{
		Name:      pbMsg.FullName,
		TableName: stringx.ToSnake(pbMsg.FullName),
		Fields:    make([]*DocField, 0),
		Indexes:   make([]*DocIndex, 0),
	}

	indexMap := make(map[string]*DocIndex)

	// 转换普通字段
	for _, field := range pbMsg.NormalFields() {
		gormTag := field.GormTag()
		docField := DocField{
			Name:        field.FieldName(),
			Type:        gormTag.Type(),
			Column:      gormTag.Column(),
			Description: field.Desc(),
			IsPrimary:   field.FieldName() == "id", // 简单判断，可以后续优化
		}

		model.Fields = append(model.Fields, &docField)

		// 生成索引信息
		if docField.IsPrimary {
			model.Indexes = append(model.Indexes, &DocIndex{
				Name:      "PRIMARY",
				Type:      "primary",
				Fields:    []string{docField.Name},
				Columns:   []string{docField.Column},
				IsUnique:  true,
				IsPrimary: true,
			})
		}

		// 处理唯一索引
		for _, indexName := range gormTag.UniqueIndex() {
			if _, ok := indexMap[indexName]; !ok {
				indexMap[indexName] = &DocIndex{
					Name:      indexName,
					Type:      "unique",
					Fields:    []string{},
					Columns:   []string{},
					IsUnique:  true,
					IsPrimary: false,
				}
			}

			indexMap[indexName].Fields = append(indexMap[indexName].Fields, docField.Name)
		}

		// 处理普通索引
		for _, indexName := range gormTag.Index() {
			if _, ok := indexMap[indexName]; !ok {
				indexMap[indexName] = &DocIndex{
					Name:      indexName,
					Type:      "index",
					Fields:    []string{},
					Columns:   []string{},
					IsUnique:  false,
					IsPrimary: false,
				}
			}

			indexMap[indexName].Fields = append(indexMap[indexName].Fields, docField.Name)
		}
	}

	model.Indexes = append(model.Indexes, candy.MapValues(indexMap)...)

	log.Infof("model:%s,index:%d", model.Name, len(model.Indexes))

	return model
}

func GenerateDoc(pb *PbPackage) (err error) {
	pterm.Info.Println("try generate documentation")

	// 获取文档输出目录
	docsDir := filepath.Join(pb.ProjectRoot(), "docs")
	if !osx.IsDir(docsDir) {
		err = os.MkdirAll(docsDir, 0755)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	modelMessages := candy.Filter(pb.Messages(), func(message *PbMessage) bool {
		return message.IsTable()
	})

	if len(modelMessages) == 0 {
		pterm.Warning.Println("no Model-prefixed messages found")
		return nil
	}

	// 转换为自定义文档结构
	docModels := make([]DocModel, 0, len(modelMessages))
	for _, pbMsg := range modelMessages {
		docModels = append(docModels, convertPbMessageToDocModel(pbMsg))
	}

	// 构建文档数据
	docData := DocData{
		PackageName:   pb.PackageName(),
		ProtoFile:     pb.ProtoFileName(),
		GoPackage:     pb.GoPackage(),
		GoPackageName: pb.GoPackageName(),
		GeneratedAt:   time.Now().Format("2006-01-02 15:04:05"),
		Models:        docModels,
	}

	tpl, err := GetTemplate(TemplateTypeDocsDatabase, i18n.DefaultLanguage())
	if err != nil {
		log.Errorf("failed to get template: %v", err)
		return err
	}

	file, err := os.OpenFile(filepath.Join(docsDir, i18n.Localize(state.I18nTagCliDocFilenameDatabaseDesign)+".md"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermDefault)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	err = tpl.Execute(file, docData)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	pterm.Success.Println("documentation generated successfully")
	return nil
}

// generateIndexName 生成索引名称
// prefix: 索引前缀 (idx, uk 等)
// modelName: 模型名称
// fieldName: 字段名称
func generateIndexName(prefix, modelName, fieldName string) string {
	// 移除 Model 前缀并转换为小写蛇形命名
	tableName := stringx.ToSnake(modelName)
	if strings.HasPrefix(tableName, "model_") {
		tableName = tableName[6:] // 移除 "model_" 前缀
	}

	// 生成索引名: prefix_tablename_fieldname
	return strings.ToLower(prefix + "_" + tableName + "_" + stringx.ToSnake(fieldName))
}

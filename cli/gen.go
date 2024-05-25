package cli

import (
	"errors"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/lazygophers/utils/runtime"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var pb *codegen.PbPackage

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// NOTE: 合并输入参数与配置文件
		mergeGenCmdFlags(cmd)

		err = codegen.CheckEnvironments()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		protoFile := filepath.Join(runtime.Pwd(), getString("input", cmd))
		pterm.Info.Println("proto file:", protoFile)

		// NOTE: 检查proto文件是否存在
		{
			if !osx.IsFile(protoFile) {
				log.Errorf("file not found:%s", protoFile)
				pterm.Error.Printfln("proto file not found:%s", protoFile)
				return errors.New("proto file not found")
			}
		}

		pb, err = codegen.ParseProto(protoFile)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	},
}

func mergeGenCmdFlags(cmd *cobra.Command) {

	if cmd.Flag("protoc").Changed {
		state.Config.ProtocPath = getString("protoc", cmd)
	}

	if cmd.Flag("protoc-gen-go").Changed {
		state.Config.ProtoGenGoPath = getString("protoc-gen-go", cmd)
	}

	if cmd.Flag("go-module-prefix").Changed {
		state.Config.GoModulePrefix = strings.TrimSuffix(getString("go-module-prefix", cmd), "/")
	}

	if cmd.Flag("output-path").Changed {
		state.Config.OutputPath = getString("output-path", cmd)
	}

	if cmd.Flag("editorconfig-path").Changed {
		state.Config.Template.Editorconfig = getString("editorconfig-path", cmd)
	}

	if cmd.Flag("overwrite").Changed {
		state.Config.Overwrite = getBool("overwrite", cmd)
	}

	if cmd.Flag("proto-files").Changed {
		state.Config.ProtoFiles = getStringSlice("proto-files", cmd)
	}

	if cmd.Flag("add-proto-files").Changed {
		state.Config.ProtoFiles = append(state.Config.ProtoFiles, getStringSlice("proto-files", cmd)...)
	}

	// NOTE: tables
	if cmd.Flag("tables-enable_field_id").Changed {
		state.Config.Tables.DisableFieldId = !getBool("tables-enable_field_id", cmd)
	}

	if cmd.Flag("tables-enable_field_created_at").Changed {
		state.Config.Tables.DisableFieldId = !getBool("tables-enable_field_created_at", cmd)
	}

	if cmd.Flag("tables-enable_field_updated_at").Changed {
		state.Config.Tables.DisableFieldUpdatedAt = !getBool("tables-enable_field_updated_at", cmd)
	}

	if cmd.Flag("tables-enable_field_deleted_at").Changed {
		state.Config.Tables.DisableFieldDeletedAt = !getBool("tables-enable_field_deleted_at", cmd)
	}

	if cmd.Flag("tables-enable_gorm_tag_column").Changed {
		state.Config.Tables.DisableGormTagColumn = !getBool("tables-enable_gorm_tag_column", cmd)
	}
}

func init() {
	genCmd.PersistentFlags().String("protoc", "", "protoc path")
	genCmd.PersistentFlags().String("protoc-gen-go", "", "protoc-gen-go path")
	genCmd.PersistentFlags().String("output-path", "", "output path")

	genCmd.PersistentFlags().String("go-module-prefix", "", "go module prefix")
	genCmd.PersistentFlags().StringSlice("proto-files", nil, "import other .proto files, if their not in the .proto file")
	genCmd.PersistentFlags().StringSlice("add-proto-files", nil, "append import other .proto files, if their not in the .proto file")

	genCmd.PersistentFlags().Bool("overwrite", false, "overwrite configuration")
	genCmd.PersistentFlags().String("editorconfig-path", "", ".editorconfig path")

	genCmd.PersistentFlags().Bool("tables-enable_field_id", false, "enable field gen tags for id field")
	genCmd.PersistentFlags().Bool("tables-enable_field_created_at", false, "enable field gen tags for created_at field")
	genCmd.PersistentFlags().Bool("tables-enable_field_updated_at", false, "enable field gen tags for updated_at field")
	genCmd.PersistentFlags().Bool("tables-enable_field_deleted_at", false, "enable field gen tags for deleted_at field")
	genCmd.PersistentFlags().Bool("tables-enable_gorm_tag_column", false, "enable gen tags for column tag for each")

	genCmd.PersistentFlags().StringP("input", "i", "", "input protobuf file")

	rootCmd.AddCommand(genCmd)
}

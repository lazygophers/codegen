package gen

import (
	"errors"
	"github.com/lazygophers/codegen/cli/utils"
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
	Use:     "gen",
	Aliases: []string{"g", "generate"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// NOTE: 合并输入参数与配置文件
		mergeGenCmdFlags(cmd)

		err = codegen.CheckEnvironments()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		protoFile := filepath.Join(runtime.Pwd(), utils.GetString("input", cmd))
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
		state.Config.ProtocPath = utils.GetString("protoc", cmd)
	}

	if cmd.Flag("protoc-gen-go").Changed {
		state.Config.ProtoGenGoPath = utils.GetString("protoc-gen-go", cmd)
	}

	if cmd.Flag("go-module-prefix").Changed {
		state.Config.GoModulePrefix = strings.TrimSuffix(utils.GetString("go-module-prefix", cmd), "/")
	}

	if cmd.Flag("output-path").Changed {
		state.Config.OutputPath = utils.GetString("output-path", cmd)
	}

	if cmd.Flag("overwrite").Changed {
		state.Config.Overwrite = utils.GetBool("overwrite", cmd)
	}

	if cmd.Flag("proto-files").Changed {
		state.Config.ProtoFiles = utils.GetStringSlice("proto-files", cmd)
	}

	if cmd.Flag("add-proto-files").Changed {
		state.Config.ProtoFiles = append(state.Config.ProtoFiles, utils.GetStringSlice("proto-files", cmd)...)
	}

	// NOTE: tables
	if cmd.Flag("tables-enable_field_id").Changed {
		state.Config.Tables.DisableFieldId = !utils.GetBool("tables-enable_field_id", cmd)
	}

	if cmd.Flag("tables-enable_field_created_at").Changed {
		state.Config.Tables.DisableFieldId = !utils.GetBool("tables-enable_field_created_at", cmd)
	}

	if cmd.Flag("tables-enable_field_updated_at").Changed {
		state.Config.Tables.DisableFieldUpdatedAt = !utils.GetBool("tables-enable_field_updated_at", cmd)
	}

	if cmd.Flag("tables-enable_field_deleted_at").Changed {
		state.Config.Tables.DisableFieldDeletedAt = !utils.GetBool("tables-enable_field_deleted_at", cmd)
	}

	if cmd.Flag("tables-enable_gorm_tag_column").Changed {
		state.Config.Tables.DisableGormTagColumn = !utils.GetBool("tables-enable_gorm_tag_column", cmd)
	}
}

func Load(rootCmd *cobra.Command) {
	rootCmd.Short = state.Localize(state.I18nTagCliGenShort)
	rootCmd.Long = state.Localize(state.I18nTagCliGenLong)

	genCmd.PersistentFlags().String("protoc", state.Config.ProtocPath, state.Localize(state.I18nTagCliGenFlagsProtoc))
	genCmd.PersistentFlags().String("protoc-gen-go", state.Config.ProtoGenGoPath, state.Localize(state.I18nTagCliGenFlagsProtocGenGo))
	genCmd.PersistentFlags().String("output-path", state.Config.OutputPath, state.Localize(state.I18nTagCliGenFlagsOutputPath))

	genCmd.PersistentFlags().String("go-module-prefix", state.Config.GoModulePrefix, state.Localize(state.I18nTagCliGenFlagsGoModulePrefix))
	genCmd.PersistentFlags().StringSlice("proto-files", state.Config.ProtoFiles, state.Localize(state.I18nTagCliGenFlagsProtoFiles))
	genCmd.PersistentFlags().StringSlice("add-proto-files", nil, state.Localize(state.I18nTagCliGenFlagsAddProtoFiles))

	genCmd.PersistentFlags().Bool("overwrite", false, state.Localize(state.I18nTagCliGenFlagsOverwrite))

	genCmd.PersistentFlags().Bool("tables-enable_field_id", state.Config.Tables.DisableFieldId, state.Localize(state.I18nTagCliGenFlagsTablesEnableFieldId))
	genCmd.PersistentFlags().Bool("tables-enable_field_created_at", state.Config.Tables.DisableFieldCreatedAt, state.Localize(state.I18nTagCliGenFlagsTablesEnableFieldCreatedAt))
	genCmd.PersistentFlags().Bool("tables-enable_field_updated_at", state.Config.Tables.DisableFieldUpdatedAt, state.Localize(state.I18nTagCliGenFlagsTablesEnableFieldUpdatedAt))
	genCmd.PersistentFlags().Bool("tables-enable_field_deleted_at", state.Config.Tables.DisableFieldDeletedAt, state.Localize(state.I18nTagCliGenFlagsTablesEnableFieldDeletedAt))
	genCmd.PersistentFlags().Bool("tables-enable_gorm_tag_column", state.Config.Tables.DisableGormTagColumn, state.Localize(state.I18nTagCliGenFlagsTablesEnableGormTagColumn))

	genCmd.PersistentFlags().StringP("input", "i", "", state.Localize(state.I18nTagCliGenFlagsInput))

	initAddRpc()
	initAll()
	initCache()
	initCmd()

	rootCmd.AddCommand(genCmd)
}

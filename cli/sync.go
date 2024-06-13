package cli

import (
	"github.com/lazygophers/codegen/state"
	"github.com/spf13/cobra"
)

// 加载远端的配置文件
var syncCmd = &cobra.Command{
	Use:         "sync",
	Aliases:     []string{"remote", "sync-remote", "update", "update-remote"},
	Example:     "",
	Args:        cobra.MaximumNArgs(1),
	Annotations: nil,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return nil
	},
}

func initSync() {
	syncCmd.Short = state.Localize(state.I18nTagCliSyncShort)
	syncCmd.Long = state.Localize(state.I18nTagCliSyncLong)

	syncCmd.Flags().String("cache-template-path", "", state.Localize(state.I18nTagCliSyncFlagsTemplatePath))

	rootCmd.AddCommand(syncCmd)
}

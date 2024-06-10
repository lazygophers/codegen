package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use: "cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		state.Config.State.Cache = true
		return runGenCache(cmd, args)
	},
}

func runGenCache(c *cobra.Command, args []string) (err error) {
	err = codegen.GenerateCache(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initCache() {
	cacheCmd.Short = state.Localize(state.I18nTagCliGenCacheShort)
	cacheCmd.Long = state.Localize(state.I18nTagCliGenCacheLong)

	genCmd.AddCommand(cacheCmd)
}

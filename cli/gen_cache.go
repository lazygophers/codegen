package cli

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var genCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "generate cache file",
	RunE:  ranGenCache,
}

func ranGenCache(c *cobra.Command, args []string) (err error) {
	err = codegen.GenerateCache(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(genCacheCmd)
}

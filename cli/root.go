package cli

import (
	"github.com/lazygophers/codegen/cli/gen"
	"github.com/lazygophers/codegen/cli/i18n"
	"github.com/lazygophers/codegen/cli/utils"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "codegen",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if utils.GetBool("debug", cmd) {
			pterm.EnableDebugMessages()
			log.SetLevel(log.DebugLevel)
			log.SetOutput(os.Stdout)

			pterm.Debug.Println("enable debug")
		} else {
			log.SetLevel(log.InfoLevel)
			log.SetOutput(io.Discard)
		}

		return nil
	},
}

func Run() {
	var err error

	cobra.EnablePrefixMatching = true
	cobra.EnableCommandSorting = true
	cobra.EnableTraverseRunHooks = true

	// 加载配置文件
	err = state.Load()
	if err != nil {
		pterm.Error.Printfln("load state error\n%s", err.Error())
		log.Errorf("err:%v", err)
		return
	}

	rootCmd.Short = state.Localize(state.I18nTagCliShort)
	rootCmd.Long = state.Localize(state.I18nTagCliLong)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, state.Localize(state.I18nTagCliFlagsDebug))

	gen.Load(rootCmd)
	i18n.Load(rootCmd)
	initSync()

	initDefault()

	err = rootCmd.Execute()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}

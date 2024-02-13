package cli

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "codegen",
	Short: "codegen",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if getBool("debug", cmd) {
			pterm.EnableDebugMessages()
			log.SetLevel(log.DebugLevel)
			log.SetOutput(os.Stdout)

			pterm.Debug.Println("enable debug")
		} else {
			log.SetLevel(log.InfoLevel)
			log.SetOutput(io.Discard)
		}

		// 加载配置文件
		err = state.Load()
		if err != nil {
			pterm.Error.Printfln("load state error\n%s", err.Error())
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	},
}

func Run() {
	var err error

	cobra.EnablePrefixMatching = true
	cobra.EnableCommandSorting = true
	cobra.EnableTraverseRunHooks = true

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug output")

	err = rootCmd.Execute()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}

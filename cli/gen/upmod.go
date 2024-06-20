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
)

var upModCmd = &cobra.Command{
	Use:     "update-mod",
	Aliases: []string{"upmod", "um", "up-mod"},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var projectDir string

		protooFile, err := getPbFile(cmd)
		if err != nil {
			log.Warnf("err:%v", err)

			// 试图寻找最近的 go.mod 文件
			find := func() (string, error) {
				path := runtime.Pwd()

				for filepath.Dir(path) != path {
					if osx.IsFile(filepath.Join(path, "go.mod")) {
						return path, nil
					}

					path = filepath.Dir(path)
				}

				return "", errors.New("not found go.mod file")
			}

			projectDir, err = find()
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		} else {
			pb, err = codegen.ParseProto(protooFile)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			projectDir = pb.ProjectRoot()
		}

		// 加载配置
		if osx.IsFile(filepath.Join(projectDir, ".lazygophers")) {
			err = state.LoadLazeConfig(filepath.Join(projectDir, ".lazygophers"))
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		} else {
			log.Warnf("not found %s", filepath.Join(projectDir, ".lazygophers"))
			pterm.Warning.Printfln("not found .lazygophers in product package, use defulat")
			state.LazyConfig.Apply()
		}

		if cmd.Flag("goproxy").Value.String() != "" {
			state.LazyConfig.GoMod.Goproxy = utils.GetString("goproxy", cmd)
		}

		err = codegen.UpdateGoMod(filepath.Join(projectDir, "go.mod"))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	},
}

func initUpmod() {
	upModCmd.Short = state.Localize(state.I18nTagCliUpModShort)
	upModCmd.Long = state.Localize(state.I18nTagCliUpModLong)

	upModCmd.Flags().String("goproxy", "", state.Localize(state.I18nTagCliUpModFlagsGoproxy))
	upModCmd.PersistentFlags().StringP("input", "i", "", state.Localize(state.I18nTagCliGenFlagsInput))
}

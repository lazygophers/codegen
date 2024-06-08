package gen

import (
	"errors"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/stringx"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"strings"
)

var addRpcCmd = &cobra.Command{
	Use:   "add-rpc",
	Short: "add rpc to proto with model",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		opt := codegen.NewAddRpcOption()

		if v, err := cmd.Flags().GetString("default-role"); err == nil {
			opt.DefaultRole = v
		}

		if v, err := cmd.Flags().GetString("gen-to"); err == nil {
			opt.GenTo = v
		}

		var msg *codegen.PbMessage
		if v, err := cmd.Flags().GetString("model"); err == nil && v != "" {
			// 先看一下加了 Model 的 时候存在
			{
				model := stringx.ToSnake(v)
				model = strings.TrimPrefix(model, "model_")
				model = "model_" + model
				model = stringx.ToCamel(model)

				msg = pb.GetMessage(model)
				if msg != nil {
					opt.Model = model
				}
			}

			// 直接找名字的
			if msg == nil {
				msg = pb.GetMessage(v)
				if msg != nil {
					opt.Model = v
				}
			}

			if msg == nil {
				log.Errorf("not found model:%v", v)
				pterm.Error.Printfln("not found model:%v", v)
				return errors.New("not found model")
			}
		} else {
			log.Errorf("missed argument `model`")
			pterm.Error.Printfln("missing argument `model`")
			return errors.New("missing argument `model`")
		}

		if v, err := cmd.Flags().GetString("action"); err == nil && v != "" {
			opt.ParseActions(v)
		} else {
			log.Errorf("missed argument `action`")
			pterm.Error.Printfln("missing argument `action`")
			return errors.New("missing argument `action`")
		}

		if v, err := cmd.Flags().GetString("list-option"); err == nil && v != "" {
			opt.ParseListOption(v, msg)
		}

		err = codegen.GenerateAddRpc(pb, msg, opt)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return err
	},
}

func init() {
	addRpcCmd.Flags().StringP("model", "m", "", "model name，must be specified")
	addRpcCmd.Flags().StringP("gen-to", "t", "", "generate go source code path, will be used in gen go source code")

	addRpcCmd.Flags().StringP("action", "a", "", "action for adding, segmente by ';',\nsupports: add/set/get/list/del/update")

	addRpcCmd.Flags().StringP("list-option", "l", "", "list options, segmente by ';'")

	addRpcCmd.Flags().String("default-role", "", "default role")

	genCmd.AddCommand(addRpcCmd)
}

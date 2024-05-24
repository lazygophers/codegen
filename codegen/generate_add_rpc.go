package codegen

import (
	"bytes"
	"fmt"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/pterm/pterm"
	"os"
	"strings"
)

type AddRpcOptionAction struct {
	Roles []string
}

type AddRpcOption struct {
	Model string

	Action map[string]*AddRpcOptionAction

	ListOptions map[string]string
}

func (p *AddRpcOption) ParseActions(s string) {
	candy.Each(strings.Split(s, ";"), func(item string) {
		idx := strings.Index(item, ":")
		if idx < 0 {
			p.Action[item] = &AddRpcOptionAction{}
			return
		}

		action := item[:idx]
		p.Action[action] = &AddRpcOptionAction{}

		candy.Each(strings.Split(item[idx+1:], ","), func(item string) {
			// TODO: 更多的格式解析
			p.Action[action].Roles = append(p.Action[action].Roles, item)
		})
	})
}

func (p *AddRpcOption) ParseListOption(s string, msg *PbMessage) {
	candy.Each(strings.Split(s, ";"), func(item string) {
		var option string
		var optionType string

		idx := strings.Index(item, ":")
		if idx < 0 {
			option = item

			// 如果没填，会按照字段类型填充默认数据
			if field, ok := msg.normalFields[option]; ok {
				option = field.Type()
			}

		} else {
			option = item[:idx]
			optionType = item[idx+1:]
		}

		p.ListOptions[option] = optionType
	})
}

func NewAddRpcOption() *AddRpcOption {
	return &AddRpcOption{
		Action:      map[string]*AddRpcOptionAction{},
		ListOptions: map[string]string{},
	}
}

func GenerateAddRpc(pb *PbPackage, opt *AddRpcOption) (err error) {
	// 重新加载一下 pb 文件
	pb, err = ParseProto(pb.protoFilePath, true)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 找到最后一个 service 里面 rpc 的位置
	var lastRpcPos int
	{
		rpc := candy.Last(pb.RPCs())
		if rpc != nil {
			// 找到这一个 rpc 的结尾
			s := rpc.RPC().Position.Offset
			e := s

			for e <= len(pb.ProtoBuffer) {
				for e < len(pb.ProtoBuffer) &&
					pb.ProtoBuffer[e] != '\n' {
					e++
				}

				line := pb.ProtoBuffer[s:e]

				line = strings.ReplaceAll(line, " ", "")
				line = strings.ReplaceAll(line, "\r", "")
				line = strings.ReplaceAll(line, "\t", "")
				if line == "};" || line == "}" {
					lastRpcPos = e
					break
				}

				s = e + 1
				e = s
			}
		}
	}

	if lastRpcPos == 0 {
		service := candy.Last(pb.Services())
		if service != nil {
			// 找到这一个 rpc 的结尾
			s := service.Service().Position.Offset
			e := s

			for e <= len(pb.ProtoBuffer) {
				for e < len(pb.ProtoBuffer) &&
					pb.ProtoBuffer[e] != '\n' {
					e++
				}

				line := pb.ProtoBuffer[s:e]

				line = strings.ReplaceAll(line, " ", "")
				line = strings.ReplaceAll(line, "\r", "")
				line = strings.ReplaceAll(line, "\t", "")
				log.Info(line)
				if line == "};" || line == "}" {
					lastRpcPos = s
					break
				}

				s = e + 1
				e = s
			}
		}
	}

	if lastRpcPos == 0 {
		pterm.Fatal.Printfln("not found service in %s, please add it", pb.ProtoFileName())
		return fmt.Errorf("not found service in %s", pb.ProtoFileName())
	}

	for action, actionOpt := range opt.Action {
		for _, role := range actionOpt.Roles {
			args := map[string]interface{}{
				"PB":     pb,
				"Model":  opt.Model,
				"Role":   role,
				"Action": action,
			}

			var rpcName string
			// 先获取一下 rpc 的名字
			{
				tpl, err := GetTemplate(TemplateTypeProtoRpcName, action)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				b := log.GetBuffer()

				if role == "default" || role == "def" {
					role = ""
				}

				err = tpl.Execute(b, map[string]interface{}{
					"PB":    pb,
					"Model": opt.Model,
					"Role":  role,
				})

				rpcName = strings.ReplaceAll(b.String(), " ", "")
				rpcName = strings.ReplaceAll(rpcName, "\n", "")
				rpcName = strings.ReplaceAll(rpcName, "\r", "")
				rpcName = strings.ReplaceAll(rpcName, "\t", "")
			}

			log.Infof("try add rpc %s", rpcName)

			args["RpcName"] = rpcName
			args["RequestType"] = rpcName + "Req"
			args["ResponseType"] = rpcName + "Rsp"

			// 处理 server.rpc
			if rpc := pb.GetRPC(rpcName); rpc == nil {
				b := bytes.NewBufferString(pb.ProtoBuffer[:lastRpcPos])
				b.WriteByte('\n')
				b.WriteByte('\n')

				b.WriteString("中文占位\n")

				//tpl, err := GetTemplate(TemplateTypeProtoService)
				//if err != nil {
				//	log.Errorf("err:%v", err)
				//	return err
				//}
				//
				//err = tpl.Execute(b, args)
				//if err != nil {
				//	log.Errorf("err:%v", err)
				//	return err
				//}

				log.Infof("move rpc position from %d to %d", lastRpcPos, b.Len())
				newPos := b.Len()

				b.WriteString(pb.ProtoBuffer[lastRpcPos+1:])
				lastRpcPos = newPos

				pb.ProtoBuffer = b.String()
			}
		}
	}

	err = os.WriteFile(pb.protoFilePath, []byte(pb.ProtoBuffer), 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

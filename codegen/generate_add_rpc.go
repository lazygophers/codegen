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

	DefaultRole string

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

	for _, action := range p.Action {
		action.Roles = candy.Map(action.Roles, func(s string) string {
			if s == "" {
				return p.DefaultRole
			}

			return s
		})

		if len(action.Roles) == 0 {
			action.Roles = append(action.Roles, p.DefaultRole)
		}

		action.Roles = candy.Unique(action.Roles)
	}
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
				optionType = field.Type()
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

func GenerateAddRpc(pb *PbPackage, msg *PbMessage, opt *AddRpcOption) (err error) {
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
					lastRpcPos = e // 当前行的换行符
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
				if line == "}" {
					lastRpcPos = s - 1
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

	// NOTE: 寻找主键
	var pkField *PbNormalField
	{
		for _, field := range msg.normalFields {
			// 先简单粗暴用 id 当作主键，后面再改
			if field.Name == "id" {
				pkField = field
				break
			}
		}
	}

	var rpcBlock string

	for action, actionOpt := range opt.Action {
		for _, role := range actionOpt.Roles {
			args := map[string]interface{}{
				"PB":          pb,
				"Model":       opt.Model,
				"Role":        role,
				"Action":      action,
				"ListOptions": opt.ListOptions,
			}

			if pkField != nil {
				args["PprimaryKey"] = pkField.Name
				args["PprimaryKeyType"] = pkField.Type()
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

				//rpcName = b.String()
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
				tpl, err := GetTemplate(TemplateTypeProtoService)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				var b bytes.Buffer
				err = tpl.Execute(&b, args)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				rpcBlock += b.String()
			}

			//  处理 request
			if req := pb.GetMessage(rpcName + "Req"); req == nil {
				tpl, err := GetTemplate(TemplateTypeProtoRpcReq, action)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				log.Info(args["ListOptions"])

				var b bytes.Buffer
				err = tpl.Execute(&b, args)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				pb.ProtoBuffer += "\n"
				pb.ProtoBuffer += b.String()
			}

			//  处理 response
			if rsp := pb.GetMessage(rpcName + "Rsp"); rsp == nil {
				tpl, err := GetTemplate(TemplateTypeProtoRpcResp, action)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				var b bytes.Buffer
				err = tpl.Execute(&b, args)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				pb.ProtoBuffer += "\n"
				pb.ProtoBuffer += b.String()
			}
		}
	}

	if rpcBlock != "" {
		b := bytes.NewBufferString(pb.ProtoBuffer[:lastRpcPos])

		b.WriteByte('\n')

		b.WriteString(rpcBlock)

		b.WriteString(pb.ProtoBuffer[lastRpcPos:])

		pb.ProtoBuffer = b.String()
	}

	err = os.WriteFile(pb.protoFilePath, []byte(pb.ProtoBuffer), 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

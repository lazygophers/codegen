package cli_test

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"reflect"
	"testing"
)

func TestSyncMerge(t *testing.T) {
	c := &state.Cfg{
		Template: &state.CfgTemplate{
			Editorconfig: "222",
		},
	}

	//nc := state.Cfg{
	//	Template: &state.CfgTemplate{
	//		Editorconfig: "1",
	//		Orm:          "2",
	//		TableName:    "3",
	//		TableField:   "4",
	//		Proto: &state.CfgProto{
	//			Action: map[string]*state.CfgProtoAction{
	//				"1": {
	//					Name: "1",
	//					Req:  "2",
	//					Resp: "3",
	//				},
	//			},
	//			Rpc: "4",
	//		},
	//		Impl: &state.CfgImpl{
	//			Action: nil,
	//			Impl:   "4",
	//			Path:   "5",
	//			Route:  "54",
	//		},
	//		Main: "1",
	//		Conf: "3",
	//		State: &state.CfgTemplateState{
	//			Table: "1",
	//			Conf:  "2",
	//			Cache: "2",
	//			State: "e",
	//			I18n:  "d",
	//		},
	//		Goreleaser:   "s",
	//		Makefile:     "w",
	//		Golangci:     "w",
	//		Gitignore:    "q",
	//		Dockerignore: "s",
	//		I18nConst:    "c",
	//	},
	//}

	ncv := reflect.ValueOf(c.Template)
	for ncv.Kind() == reflect.Ptr {
		ncv = ncv.Elem()
	}

	var deep func(dst reflect.Value) error
	deep = func(dst reflect.Value) error {
		for dst.Kind() == reflect.Ptr {
			dst = dst.Elem()
		}

		if !dst.IsValid() {
			return nil
		}

		if dst.IsZero() {
			return nil
		}

		switch dst.Kind() {
		case reflect.String:
			dst.SetString("sss")

		case reflect.Struct:
			for i := 0; i < dst.NumField(); i++ {
				if err := deep(dst.Field(i)); err != nil {
					return err
				}
			}

		case reflect.Map:
			for _, key := range dst.MapKeys() {
				if err := deep(dst.MapIndex(key)); err != nil {
					return err
				}
			}

		case reflect.Slice:
			for i := 0; i < dst.Len(); i++ {
				if err := deep(dst.Index(i)); err != nil {
					return err
				}
			}

		default:
			log.Panicf("unsupported type %s", dst.Kind())
		}

		return nil
	}

	t.Log(deep(ncv))
	t.Log(ncv.Interface())
}

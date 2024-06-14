package i18n

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/json"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Localizer interface {
	MarshalToFile(path string, v any) (err error)
	UnmarshalFromFile(path string, v any) (err error)
}

var localizer = map[string]Localizer{
	"json": jsonLocalizer,
	"yaml": yamlLocalizer,
	"yml":  yamlLocalizer,
	"toml": tomlLocalizer,
}

var (
	jsonLocalizer = NewLocalizerHandle(
		func(path string, v any) (err error) {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "\t")

			err = encoder.Encode(v)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return nil
		},
		func(path string, v any) (err error) {
			file, err := os.Open(path)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer file.Close()

			err = json.NewDecoder(file).Decode(v)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return nil
		},
	)

	yamlLocalizer = NewLocalizerHandle(
		func(path string, v any) (err error) {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer file.Close()

			err = yaml.NewEncoder(file).Encode(v)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return nil
		},
		func(path string, v any) (err error) {
			file, err := os.Open(path)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer file.Close()

			err = yaml.NewDecoder(file).Decode(v)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return nil
		},
	)

	tomlLocalizer = NewLocalizerHandle(
		func(path string, v any) (err error) {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer file.Close()

			err = toml.NewEncoder(file).Encode(v)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return nil
		},
		func(path string, v any) (err error) {
			file, err := os.Open(path)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer file.Close()

			err = toml.NewDecoder(file).Decode(v)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return nil
		},
	)
)

func RegisterLocalizer(name string, v Localizer) {
	localizer[name] = v
}

func GetLocalizer(name string) (Localizer, bool) {
	if strings.HasPrefix(name, ".") {
		name = name[1:]
	}

	l, ok := localizer[name]
	return l, ok
}

var _ Localizer = (*LocalizerHandle)(nil)

type LocalizerHandle struct {
	marshal   func(path string, v any) (err error)
	unmarshal func(path string, v any) (err error)
}

func (p *LocalizerHandle) MarshalToFile(path string, v any) (err error) {
	return p.marshal(path, v)
}

func (p *LocalizerHandle) UnmarshalFromFile(path string, v any) (err error) {
	return p.unmarshal(path, v)
}

func NewLocalizerHandle(marshal func(path string, v any) (err error), unmarshal func(path string, v any) (err error)) *LocalizerHandle {
	return &LocalizerHandle{
		marshal:   marshal,
		unmarshal: unmarshal,
	}
}

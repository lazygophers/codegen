package i18n

import (
	"bytes"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/runtime"
	"path/filepath"
	"text/template"
)

func updateAutoTrain(c *state.CfgI18nAutoTran, localizeDir string) error {
	if c.EnableRecord {
		t, err := template.New("").Parse(c.RecordPath)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		var b bytes.Buffer
		err = t.Funcs(map[string]any{
			"Join": filepath.Join,
			"Dir":  filepath.Dir,
		}).Execute(&b, map[string]any{
			"Pwd":         runtime.Pwd(),
			"LocalizeDir": localizeDir,
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		c.RecordPath = filepath.ToSlash(b.String())
	}

	return nil
}

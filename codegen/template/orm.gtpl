package {{ .PB.GoPackageName }}

import (
	"database/sql/driver"
	"github.com/lazygophers/utils"
)

{{ range $key, $value := .Models }}func (m *{{ $value }}) Scan(value interface{}) error {
	return utils.Scan(m, value)
}

func (m *{{ $value }}) Value() (driver.Value, error) {
	return utils.Value(m)
}
{{ end }}

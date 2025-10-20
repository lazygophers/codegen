package codegen

import (
	"io/fs"
	"strings"
)

const (
	// FilePermDefault 默认文件权限
	FilePermDefault = fs.FileMode(0644)
)

func clearText(content string) string {
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\r", "")

	return content
}

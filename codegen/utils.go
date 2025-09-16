package codegen

import "strings"

func clearText(content string) string {
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\r", "")

	return content
}

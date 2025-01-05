package server

import (
	"strings"
	_ "embed"
)

//go:embed script.js
var scriptContent string

func InjectScript(html []byte, addr string) string {
	injScript := strings.ReplaceAll(scriptContent, "@addr", addr)

	page := strings.Builder{}
	page.Write(html)
	page.WriteString("<script>" + injScript + "</script>")

	return page.String()
}

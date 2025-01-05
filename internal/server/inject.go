package server

import (
	"fmt"
	"strings"
)

func InjectScript(html []byte, addr string) string {
	script := fmt.Sprintf(`
		<script>
			const ws = new WebSocket("ws://%s/ws")
			ws.onopen = () => {
				console.log("[Kronos] Development server connected.")
			}
			ws.onmessage = (event) => location.reload()
			ws.onclose = () => console.error("[Kronos] Development server disconnected.")
		</script>`, addr)

	page := strings.Builder{}
	page.Write(html)
	page.WriteString(script)

	return page.String()
}

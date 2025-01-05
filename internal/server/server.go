package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	ws "trxharu.dev/kronos-build/internal/websocket"
)


func StartServer(rootpath string, addr string, ws *ws.WebSocket) (*http.Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		indexPage, _ := os.ReadFile(rootpath + "/index.html")
		ctx.Data(200, "text/html; charset=utf-8", []byte(InjectScript(indexPage, addr)))
	})

	assets := path.Join(rootpath, "assets")
	router.Static("/assets", assets)
	router.StaticFile("/favicon.ico", rootpath + "/favicon.ico")	

	// WebSocket Handler
	router.GET("/ws", gin.WrapH(ws.GetWsHandler(addr)))

	log.Printf("Local server started on http://%s.", addr)

	server := &http.Server {
		Addr: addr,
		Handler: router, 
	}

	go server.Serve(ln)
	_ = browser.OpenURL("http://" + addr)
	
	return server, err 
}


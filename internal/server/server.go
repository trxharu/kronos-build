package server

import (
	"errors"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path"

	ws "trxharu.dev/build-tool/internal/websocket"
)


func StartServer(rootpath string, addr string, ws *ws.WebSocket) (*http.Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	muxer := http.NewServeMux()
	// Static File Serving
	muxer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexPage, _ := os.ReadFile(rootpath + "/index.html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(InjectScript(indexPage, addr)))
	})

	muxer.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		favicon, err := os.ReadFile(rootpath + "/favicon.ico")	
		if errors.Is(err, fs.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte{})
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(favicon)
		}
	})

	assets := path.Join(rootpath, "assets")
	muxer.Handle("/assets", http.FileServer(http.Dir(assets)))
	// WebSocket Handler
	muxer.HandleFunc("/ws", ws.GetWsHandler(addr))

	log.Printf("Local server started on http://%s.", addr)

	server := &http.Server {
		Addr: addr,
		Handler: muxer, 
	}

	go server.Serve(ln)
	return server, err 
}


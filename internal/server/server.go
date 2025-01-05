package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"github.com/gorilla/websocket"
)


func StartServer(path string, addr string, ws func(*websocket.Conn)) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	muxer := http.NewServeMux()
	muxer.Handle("/", http.FileServer(http.Dir(path)))
	muxer.HandleFunc("/ws", getWsHandler(ws))

	log.Printf("Local server started on http://%s.", addr)

	server := &http.Server {
		Addr: addr,
		Handler: muxer, 
	}

	intSignal := make(chan os.Signal, 1)
	signal.Notify(intSignal, os.Interrupt)

	go server.Serve(ln)
	<- intSignal

	log.Println("Shutting down local server.")
	err = server.Close()
	return err 
}

func getWsHandler(wsCallback func(ws *websocket.Conn)) http.HandlerFunc {

	var upgrader = websocket.Upgrader {
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	wsHandler := func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)	
		if err != nil {
			log.Println(err)
		}
		log.Println("Client Connected")
		wsCallback(ws)
	}

	return wsHandler
}

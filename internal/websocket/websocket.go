package websocket

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type WebSocket struct {
	Conn *websocket.Conn
	WsChan chan Message
	wmu sync.Mutex
}

func (ws *WebSocket) Read() {
	_, _, err := ws.Conn.ReadMessage()
	if err != nil {
		return	
	}
}

func (ws *WebSocket) Write() {
	for {
		ws.wmu.Lock()
		if err := ws.Conn.WriteJSON(<- ws.WsChan); err != nil {
			log.Println(err.Error())	
		}
		ws.wmu.Unlock()
	}
}

func (ws *WebSocket) Close() error {
	if ws.Conn != nil {
		return ws.Conn.Close()
	} else {
		return errors.New("WebSocket connection not active.")
	}
}

func (ws *WebSocket) GetWsHandler(hostaddr string) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool { 
			if r.Host == hostaddr {
				return true 
			} else {
				return false
			}
		},
	}

	wsHandler := func(rw http.ResponseWriter, r *http.Request) {
		w, err := upgrader.Upgrade(rw, r, nil)	
		if err != nil {
			log.Println(err)
		} else {
			ws.Conn = w	
		}
		go ws.Write()
	}
	return wsHandler
}

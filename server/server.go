package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func StartServer(path string, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	indexHandler := http.FileServer(http.Dir(path))

	log.Printf("Local server started on http://%s", addr)

	server := &http.Server {
		Addr: addr,
		Handler: indexHandler, 
	}

	intSignal := make(chan os.Signal, 1)
	signal.Notify(intSignal, os.Interrupt)

	go server.Serve(ln)
	<- intSignal
	log.Println("Shutting down local server...")

	err = server.Close()
	return err 
}

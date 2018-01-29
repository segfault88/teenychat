package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "net/http/pprof"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(errors.Wrap(err, "upgrader failed"))
		return
	}

	conn.WriteJSON(map[string]interface{}{"hello": "world"})
	go func() {
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("ReadMessage failed", err)
				return
			}
			log.Printf("Message type: %d message: %s", messageType, message)
			response := append([]byte("pong: "), message...)
			if err = conn.WriteMessage(1, response); err != nil {
				log.Println("WriteMessage 'pong' failed", err)
				return
			}

		}
	}()
}

func main() {
	fmt.Printf("Teenychat\n")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndex)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	srv := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// seperate server for pprof and stuff
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			log.Printf("pprof listen: %s\n", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	spew.Dump("Server stopped") // keep spew for now
}

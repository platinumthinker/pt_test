package transport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/platinumthinker/pt_test/src/server/processor"
)

type WebSocketListener struct {
	server http.Server
	mux    *http.ServeMux
}

func NewWebSocketListener(addr string) *WebSocketListener {
	srv := http.Server{
		Addr: addr,
	}

	mux := http.NewServeMux()
	srv.Handler = mux

	ws := &WebSocketListener{
		server: srv,
		mux:    mux,
	}

	return ws
}

func (ws *WebSocketListener) Listen() error {
	if err := ws.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	return nil
}

func (ws *WebSocketListener) Handle(path string, p processor.Processor) {
	ws.mux.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			ws.handleMsg(w, r, p)
		})

}

func (ws *WebSocketListener) handleMsg(w http.ResponseWriter, r *http.Request,
	p processor.Processor) {

	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// Read messages from socket
	for {
		typeMsg, msg, err := conn.ReadMessage()
		defer conn.Close()
		if err != nil {
			return
		}

		if typeMsg == websocket.TextMessage {
			res, err := p.Process(msg)
			if err != nil {
				fmt.Printf("error: Process msg: %s, error %v", msg, err)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, res)
			if err != nil {
				return
			}
		}
	}
}

func (ws *WebSocketListener) Close(ctx context.Context) {
	if err := ws.server.Shutdown(ctx); err != nil {
		log.Printf("error: HTTP server Shutdown: %v", err)
	}
}

package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Operation string `json:"operation"`
	Arguments []int  `json:"arguments,omitempty"`
}

var (
	addr = flag.String("addr", "127.0.0.1:8080", "http service address")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			args := []int{rand.Intn(1000), rand.Intn(1000)}
			msgSum, _ := json.Marshal(Message{Operation: "sum", Arguments: args})
			msgMul, _ := json.Marshal(Message{Operation: "mul", Arguments: args})

			err := c.WriteMessage(websocket.TextMessage, msgSum)
			if err != nil {
				log.Println("write:", err)
				return
			}
			err = c.WriteMessage(websocket.TextMessage, msgMul)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

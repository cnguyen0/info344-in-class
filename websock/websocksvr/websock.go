package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

/*
TODO: Implement the code in this file, according to the comments.
If you haven't yet read the assigned reading, now would be a
good time to do so:
- Read the Overview section of the Gorilla WebSockets package
https://godoc.org/github.com/gorilla/websocket
- Read the Writing WebSocket Client Application
https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_client_applications
*/

//WebSocketsHandler is a handler for WebSocket upgrade requests
type WebSocketsHandler struct {
	notifier *Notifier
	upgrader *websocket.Upgrader
}

//NewWebSocketsHandler constructs a new WebSocketsHandler
func NewWebSocketsHandler(notifier *Notifier) *WebSocketsHandler {
	return &WebSocketsHandler{
		notifier: notifier,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

//ServeHTTP implements the http.Handler interface for the WebSocketsHandler
func (wsh *WebSocketsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("received websocket upgrade request")
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	wsh.notifier.AddClient(conn)
}

//Notifier is an object that handles WebSocket notifications
type Notifier struct {
	clients []*websocket.Conn
	eventQ  chan []byte
	sync.RWMutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	//TODO: call the .start() method on
	//a new goroutine to start the
	//event notification loop
	notifier := &Notifier{
		clients: make([]*websocket.Conn, 0),
		eventQ:  make(chan []byte, 0),
	}
	go notifier.start()
	return notifier
}

//AddClient adds a new client to the Notifier
func (n *Notifier) AddClient(client *websocket.Conn) {
	log.Println("adding new WebSockets client")
	n.RLock()
	defer n.RUnlock()
	n.clients = append(n.clients, client)
	//TODO: add the client to the `clients` slice
	//but since this can be called from multiple
	//goroutines, make sure you protect the `clients`
	//slice while you add a new connection to it!

	for {
		if _, _, err := client.NextReader(); err != nil {
			client.Close()
			break
		}
	}

	//also process incoming control messages from
	//the client, as described in this section of the docs:
	//https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
}

//Notify broadcasts the event to all WebSocket clients
func (n *Notifier) Notify(event []byte) {
	log.Printf("adding event to the queue")
	n.eventQ <- event
	//TODO: add `event` to the `n.eventQ`
}

//start starts the notification loop
func (n *Notifier) start() {
	log.Println("starting notifier loop")

	for {
		read := <-n.eventQ
		for _, client := range n.clients {
			if err := client.WriteMessage(websocket.TextMessage, read); err != nil {
				log.Println(err)
				return
			}
		}
	}

	//TODO: start a never-ending loop that reads
	//new events out of the `n.eventQ` and broadcasts
	//them to all WebSocket clients.

	//If you use additional channels instead of a mutex
	//to protext the `clients` slice, also process those
	//channels here using a non-blocking `select` statement
}

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message represents a published message to be sent to subscribers
type Message struct {
	Type int
	Body []byte
}

// pubication is an open channel of messages
var publication = make(chan Message)

// subscriptions are open connections to read published messages
var subscriptions = make(map[*websocket.Conn]bool)

// default websocket config
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	// create the websocket connection and close it when ready
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// This is where new subscriptions would be handled and given an open connection.
	// I'm not sure how to implement this without setting up an API for subscribing
	// and saving the subscriptions to a database and then opening connections for each
	// subscription, so I've added this simple single subscription.
	subscriptions[conn] = true

	for {
		// read the message to be published and pass it to the publication channel
		t, b, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		newMessage := Message{
			Type: t,
			Body: b,
		}
		publication <- newMessage
	}
}

// publish takes the next message from the publication channel and sends it to each subscription
func publish() {
	for {
		for subscription := range subscriptions {
			nextMessage := <-publication
			if err := subscription.WriteMessage(nextMessage.Type, nextMessage.Body); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handler)

	go publish()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}
}

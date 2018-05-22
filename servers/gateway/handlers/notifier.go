package handlers

import (
	"github.com/gorilla/websocket"
	"sync"
)

//Notifier is an object that handles WebSocket notifications
type Notifier struct {
	// MAY NOT NEED EVENTQ
	eventQ chan []byte
	clients map[int64]*websocket.Conn
	mx sync.RWMutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	n := &Notifier{
		eventQ: make(chan []byte, 1024), //buffered channel that can hold 1024 slices at a time
	}
	go n.start()
	return n
}

//AddClient adds a new client to the Notifier
func (n *Notifier) AddClient(client *websocket.Conn, userID int64) {
	n.mx.Lock()
	defer n.mx.Unlock()

	n.clients[userID] = client
	go n.processControlMsgs(userID)
}

func (n *Notifier) processControlMsgs(userID int64) {
	for {
		if _, _, err := n.clients[userID].NextReader(); err != nil {
			n.clients[userID].Close()
			delete(n.clients, userID)
			break
		}
	}
}

// MAY NOT NEED ANYTHING BELOW
//Notify broadcasts the event to all WebSocket clients
func (n *Notifier) Notify(event []byte) {
	n.mx.Lock()
	defer n.mx.Unlock()

	n.eventQ <- event
}

//start starts the notification loop
func (n *Notifier) start() {
	n.mx.RLock()
	defer n.mx.RUnlock()

	for evt := range n.eventQ {
		prepMsg, err := websocket.NewPreparedMessage(websocket.TextMessage, evt)
		if err != nil {
			//handle error
		}
		for user, conn := range n.clients {
			if err = conn.WritePreparedMessage(prepMsg); err != nil {
				delete(n.clients, user)
			}
		}
	}
}
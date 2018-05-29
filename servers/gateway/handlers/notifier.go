package handlers

import (
	"github.com/gorilla/websocket"
	"sync"
	"github.com/labstack/gommon/log"
)

//Notifier is an object that handles WebSocket notifications
type Notifier struct {
	clients map[int64]*websocket.Conn
	mx sync.RWMutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	n := make(map[int64]*websocket.Conn)
	return &Notifier{
		clients: n,
	}
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
			n.mx.Lock()
			n.clients[userID].Close()
			delete(n.clients, userID)
			n.mx.Unlock()
			break
		}
	}
}

func (n *Notifier) start(msg []byte, users []int64) {
	prepMsg, err := websocket.NewPreparedMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Errorf("error writing a new prepared message: %v", err)
		return
	}
	if len(users) == 0 {
		for userID, conn := range n.clients {
			_, exists := n.clients[userID]
			if exists {
				n.writePreppedMsg(prepMsg, userID, conn)
			}
		}
	} else {
		for _, userID := range users {
			_, exists := n.clients[userID]
			if exists {
				n.writePreppedMsg(prepMsg, userID, n.clients[userID])
			}
		}
	}
}

//writePreppedMsg writes a prepared message to the given connection
func (n *Notifier) writePreppedMsg(msg *websocket.PreparedMessage, userID int64, conn *websocket.Conn) {
	if err := conn.WritePreparedMessage(msg); err != nil {
		n.mx.Lock()
		delete(n.clients, userID)
		n.mx.Unlock()
	}
}
package handlers

import (
	"github.com/gorilla/websocket"
	"sync"
	"encoding/json"
	"github.com/labstack/gommon/log"
)

//Notifier is an object that handles WebSocket notifications
type Notifier struct {
	clients map[int64]*websocket.Conn
	mx sync.RWMutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	return &Notifier{}
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

func (n *Notifier) start(msg *mqMsg, users []int64) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("error marshalling mqMsg object: %v", err)
		return
	}
	prepMsg, err := websocket.NewPreparedMessage(websocket.TextMessage, jsonMsg)
	if err != nil {
		log.Errorf("error writing a new prepared message: %v", err)
		return
	}
	if len(users) == 0 {
		for userID, conn := range n.clients {
			n.writePreppedMsg(prepMsg, userID, conn)
		}
	} else {
		for _, userID := range users {
			n.writePreppedMsg(prepMsg, userID, n.clients[userID])
		}
	}
}

//writePreppedMsg writes a prepared message to the given connection
func (n *Notifier) writePreppedMsg(msg *websocket.PreparedMessage, userID int64, conn *websocket.Conn) {
	if err := conn.WritePreparedMessage(msg); err != nil {
		delete(n.clients, userID)
	}
}
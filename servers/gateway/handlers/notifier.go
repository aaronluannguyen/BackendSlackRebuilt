package handlers

import (
	"github.com/gorilla/websocket"
	"sync"
	"github.com/labstack/gommon/log"
)

//Notifier is an object that handles WebSocket notifications
type Notifier struct {
	clients map[int64][]*websocket.Conn
	mx sync.RWMutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	n := make(map[int64][]*websocket.Conn)
	return &Notifier{
		clients: n,
	}
}

//AddClient adds a new client to the Notifier
func (n *Notifier) AddClient(client *websocket.Conn, userID int64) {
	n.mx.Lock()
	defer n.mx.Unlock()

	n.clients[userID] = append(n.clients[userID], client)
	go n.processControlMsgs(userID)
}

func (n *Notifier) processControlMsgs(userID int64) {
	for {
		for i, webConn := range n.clients[userID] {
			if _, _, err := webConn.NextReader(); err != nil {
				n.mx.Lock()
				if webConn != nil {
					webConn.Close()
					if len(n.clients[userID]) == 1 {
						delete(n.clients, userID)
					} else {
						n.clients[userID] = append(n.clients[userID][:i], n.clients[userID][i + 1:]...)
					}
				}
				n.mx.Unlock()
				break
			}
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
		for userID := range n.clients {
			_, exists := n.clients[userID]
			if exists {
				for _, webConn := range n.clients[userID] {
					n.writePreppedMsg(prepMsg, userID, webConn)
				}
			}
		}
	} else {
		for _, userID := range users {
			_, exists := n.clients[userID]
			if exists {
				for _, webConn := range n.clients[userID] {
					n.writePreppedMsg(prepMsg, userID, webConn)
				}
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
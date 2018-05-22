package handlers

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/challenges-aaronluannguyen/servers/gateway/sessions"
	"fmt"
	"github.com/streadway/amqp"
	"time"
	"github.com/labstack/gommon/log"
)

const maxConnRetries = 5

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket


//WebSocketsHandler is a handler for WebSocket upgrade requests
type WebSocketsHandler struct {
	notifier *Notifier
	upgrader websocket.Upgrader
	ctx Context
}

//NewWebSocketsHandler constructs a new WebSocketsHandler
func NewWebSocketsHandler(notifier *Notifier, ctx Context) *WebSocketsHandler {
	newWSH := &WebSocketsHandler{
		notifier: notifier,
		upgrader: websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
		},
		ctx: ctx,
	}
	return newWSH
}

//ServeHTTP implements the http.Handler interface for the WebSocketsHandler
func (wsh *WebSocketsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionState := SessionState{}
	_, err := sessions.GetState(r, wsh.ctx.SigningKey, wsh.ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(w, fmt.Sprintf("error not authorized: %v", err), http.StatusUnauthorized)
	}
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error unable to upgrade to websocket: %v", err), http.StatusInternalServerError)
	}
	wsh.notifier.AddClient(conn, sessionState.User.ID)
}

func StartMQ(mqAddr string, mqName string) {
	conn, err := ConnectToMQ(mqAddr)
	if err != nil {
		log.Fatalf("error dialing MQ: err", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("error getting channel: %v", err)
	}
	q, err := channel.QueueDeclare(
		mqName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("error declaring queue: %v", err)
	}
	msgs, err := channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	go processMessages(msgs)
}

//ConnectToMQ keeps trying to connect to message queue by waiting longer each time until successful
func ConnectToMQ(addr string) (*amqp.Connection, error) {
	mqURL := "amqp://" + addr
	var conn *amqp.Connection
	var err error
	for i := 1; i <= maxConnRetries; i++ {
		conn, err = amqp.Dial(mqURL)
		if err == nil {
			return conn, nil
		}
		time.Sleep(time.Second * time.Duration(i * 2))
	}
	return nil, err
}

//
func processMessages(msgs <- chan amqp.Delivery) {
	for msg := range msgs {
		//Do stuff here
		msg.Ack(false)
	}
}
package handlers

import (
	"github.com/streadway/amqp"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"time"
)

const maxConnRetries = 5

//StartMQ starts a go routine that connects to a rabbitMQ server and reads new events
//from the message queue
func (ctx Context) StartMQ(mqAddr string, mqName string) {
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

	go ctx.processMessages(msgs)
}

//ConnectToMQ keeps trying to connect to message queue by waiting longer each time until successful
func ConnectToMQ(addr string) (*amqp.Connection, error) {
	mqURL := "amqp://" + addr
	var conn *amqp.Connection
	var err error
	for i := 1; i <= maxConnRetries; i++ {
		conn, err = amqp.Dial(mqURL)
		if err == nil {
			log.Printf("successfully connected to MQ")
			return conn, nil
		}
		time.Sleep(time.Second * time.Duration(i * 2))
	}
	return nil, err
}

//processMessages broadcasts messages to designated users if a private channel,
//if public, then broadcasts to all users that are connected via a websocket
func (ctx Context) processMessages(msgs <- chan amqp.Delivery) {
	for msg := range msgs {
		msgObj := &mqMsg{}
		byteMsg := []byte(msg.Body)
		if err := json.Unmarshal(byteMsg, msgObj); err != nil {
			log.Errorf("error unmarshalling message: %v", err)
			return
		}
		ctx.Notifier.start(msgObj, msgObj.UserIDs)
		msg.Ack(false)
	}
}

type mqMsg struct {
	MsgType 		string 			`json:"msgType"`
	Msg				interface{}		`json:"msg,omitempty"`
	UserIDs			[]int64			`json:"userIDs"`
}
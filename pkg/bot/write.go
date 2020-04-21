package bot

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tjarratt/babble"
)

type messageData struct {
	NickName       string `json:"nickname"`
	Message        string `json:"message"`
	EventSubdomain string `json:"event_subdomain"`
}

type message struct {
	Action string      `json:"action"`
	Data   messageData `json:"data"`
}

func (b bot) WriteMessages() {
	var messagesSend int

	defer func() {
		b.byeBye()

		log.WithFields(log.Fields{
			"bot":      b.NickName,
			"last msg": "bye bye",
		}).Info("disconnecting")
		b.conn.Close()
	}()

	for messagesSend = 1; messagesSend <= b.NumMessages; messagesSend++ {
		select {
		case <-b.quit:
			return
		default:
			b.writeMessage(messagesSend)
		}
	}
}

func (b bot) writeMessage(n int) bool {
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	sleepTime := time.Duration(rand.Int63n(b.Config.MaxDelay-int64(b.Config.MinDelay)) + b.Config.MinDelay)

	time.Sleep(sleepTime * time.Second)
	text := fmt.Sprintf("msg %d of %d, latency %d sec : msg %s",
		n, b.Config.NumMessages,
		sleepTime,
		babbler.Babble())

	msg := message{
		Action: "channelChatUserOnMessage",
		Data: messageData{
			NickName:       b.NickName,
			Message:        text,
			EventSubdomain: b.Config.SudDomain,
		},
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"bot":   b.NickName,
			"error": err,
		}).Error("write message marshal error")
		return false
	}

	if err := b.conn.WriteMessage(websocket.TextMessage, msgByte); err != nil {
		log.WithFields(log.Fields{
			"bot":   b.NickName,
			"error": err,
		}).Error("unable to write message to chat")
		return false
	}

	log.WithFields(log.Fields{
		"msg:": text,
		"bot":  b.NickName,
	}).Info("message written")

	return true
}

func (b bot) byeBye() {
	text := "bye bye"
	msg := message{
		Action: "channelChatUserOnMessage",
		Data: messageData{
			NickName:       b.NickName,
			Message:        text,
			EventSubdomain: b.Config.SudDomain,
		},
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"bot":   b.NickName,
			"error": err,
		}).Error("write message marshal error")
		return
	}

	if err := b.conn.WriteMessage(websocket.TextMessage, msgByte); err != nil {
		log.WithFields(log.Fields{
			"bot":   b.NickName,
			"error": err,
		}).Error("unable to write message to chat")
		return
	}
}

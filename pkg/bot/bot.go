package bot

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tjarratt/babble"
)

const (
	joinChatAction    = "channelStreamJoinUser"
	userMessageAction = "channelChatUserOnMessage"
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

type connData struct {
	EventSubdomain string `json:"event_subdomain"`
	IsOrganizer    bool   `json:"is_organizer"`
	NickName       string `json:"nickname"`
}

type connection struct {
	Action string   `json:"action"`
	Data   connData `json:"data"`
}

// Bot depics new bot
type bot struct {
	conn            *websocket.Conn
	url             url.URL
	nickName        string
	isOrganizer     bool
	eventSudbodmain string
	NumMessages     int
	SleepSecondsMin int
	SleepSecondsMax int
}

// New create new bot
func New(url url.URL, nickName string, isOrganizer bool, subdomain string) bot {
	bot := bot{
		url:             url,
		nickName:        nickName,
		isOrganizer:     isOrganizer,
		eventSudbodmain: subdomain,
	}

	return bot
}

// Connect connect with wss servie
func (b *bot) Connect() {
	var sleepTime time.Duration
	sleepTime = time.Duration(rand.Int63n(int64(b.SleepSecondsMax)-int64(b.SleepSecondsMin)) + int64(b.SleepSecondsMin))
	time.Sleep(sleepTime * time.Second)

	log.WithFields(log.Fields{
		"bot": b.nickName,
	}).Info("connected")

	c, _, err := websocket.DefaultDialer.Dial(b.url.String(), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("connection error")
		b.conn = nil
		return
	}

	b.conn = c
}

func (b bot) JoinChat() bool {
	joinChat := connection{
		Action: joinChatAction,
		Data: connData{
			EventSubdomain: b.eventSudbodmain,
			NickName:       b.nickName,
			IsOrganizer:    b.isOrganizer,
		},
	}

	msgByte, err := json.Marshal(joinChat)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.nickName,
		}).Error("joinchat marshal error")
		return false
	}

	if err := b.conn.WriteMessage(websocket.TextMessage, msgByte); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.nickName,
		}).Error("unable to join to chat")
	}
	return true
}

func (b bot) WriteMessage(text string) bool {
	msg := message{
		Action: "channelChatUserOnMessage",
		Data: messageData{
			NickName:       b.nickName,
			Message:        text,
			EventSubdomain: b.eventSudbodmain,
		},
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"bot":   b.nickName,
			"error": err,
		}).Error("write message marshal error")
		return false
	}

	if err := b.conn.WriteMessage(websocket.TextMessage, msgByte); err != nil {
		log.WithFields(log.Fields{
			"bot":   b.nickName,
			"error": err,
		}).Error("unable to write message to chat")
		return false
	}

	return true
}

func (b bot) Disconnect() {
	if err := b.conn.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.nickName,
		}).Error("unable to call close over socket")
	}
}

func (b bot) Do() {
	if b.conn == nil {
		return
	}

	if !b.JoinChat() {
		return
	}
	defer func() {
		log.WithFields(log.Fields{
			"bot": b.nickName,
		}).Info("disconnecting")

		b.Disconnect()
	}()

	babbler := babble.NewBabbler()
	babbler.Separator = " "

	for i := 1; i < b.NumMessages; i++ {
		var sleepTime time.Duration
		sleepTime = time.Duration(rand.Int63n(int64(b.SleepSecondsMax)-int64(b.SleepSecondsMin)) + int64(b.SleepSecondsMin))

		msg := fmt.Sprintf("msg %d of %d, latency %d sec : msg %s",
			i, b.NumMessages,
			sleepTime,
			babbler.Babble())

		if b.WriteMessage(msg) {
			log.WithFields(log.Fields{
				"message": msg,
				"bot":     b.nickName,
			}).Info("message written")

			time.Sleep(sleepTime * time.Second)
		} else {
			return
		}
	}

	b.WriteMessage("bye bye")
	log.WithFields(log.Fields{
		"message": "bye bye",
		"bot":     b.nickName,
	}).Infof("bye msg written, waiting 5 seconds to end bot")
	time.Sleep(5 * time.Second)
}

func (b bot) Listen() {

	for {
		msgType, msg, err := b.conn.ReadMessage()

		log.WithFields(log.Fields{
			"error": err,
			"type":  msgType,
			"msg":   string(msg),
		}).Info("message received")
	}
}

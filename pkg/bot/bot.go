package bot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	defaultMinDuraction = 10
	defaultMaxDuraction = 20
	defaultNumMessages  = 10
	defaultSchema       = "wss"
	defaultHost         = "6vfdhz6o24.execute-api.us-east-1.amazonaws.com"
	defaultPath         = "/beta"
	defaultSubDomain    = "noisescapes"

	joinChatAction    = "channelStreamJoinUser"
	userMessageAction = "channelChatUserOnMessage"
)

type channelData struct {
	EventSubdomain string `json:"event_subdomain"`
	IsOrganizer    bool   `json:"is_organizer"`
	NickName       string `json:"nickname"`
}
type joinChannel struct {
	Action string      `json:"action"`
	Data   channelData `json:"data"`
}

// Config Config for new bot
type Config struct {
	NickName    string
	SudDomain   string
	NumMessages int
	MinDelay    int64
	MaxDelay    int64
	Host        string
	Path        string
	Schema      string
}

type bot struct {
	conn  *websocket.Conn
	delay time.Duration
	Config
	quit chan bool
}

// New creates new bot
func New(conf Config, quit chan bool) bot {
	return bot{
		nil,
		0,
		conf,
		quit,
	}
}

func (b *bot) SetDelay(min, max int) {
	b.Config.MinDelay = int64(min)
	b.Config.MaxDelay = int64(max)
}

func (b *bot) SetURL(schema, host, path string) {
	if len(schema) == 0 {
		schema = defaultSchema
	}

	if len(host) == 0 {
		host = defaultHost
	}

	if len(path) == 0 {
		path = defaultPath
	}

	b.Config.Host = host
	b.Config.Path = path
	b.Config.Schema = schema
}

func (b *bot) SetSubdomain(subDomain string) {
	if len(subDomain) == 0 {
		subDomain = defaultSubDomain
	}

	b.Config.SudDomain = subDomain
}

func (b *bot) SetNumberOfMessages(n int) {
	b.Config.NumMessages = n
}

func (b *bot) SetNickName(name string) {
	b.Config.NickName = name
}

func (b *bot) Connec() bool {
	url := fmt.Sprintf("%s://%s%s", b.Schema, b.Host, b.Path)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"url":   url,
		}).Error("connection error")
		b.conn = nil
		return false
	}

	log.WithFields(log.Fields{
		"bot": b.Config.NickName,
	}).Info("connected")
	b.conn = c
	time.Sleep(1 * time.Second)

	return true
}

func (b bot) JoinChat() bool {
	joinChat := joinChannel{
		Action: joinChatAction,
		Data: channelData{
			EventSubdomain: b.Config.SudDomain,
			NickName:       b.NickName,
			IsOrganizer:    false,
		},
	}

	msgByte, err := json.Marshal(joinChat)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.NickName,
		}).Error("joinchat marshal error")
		return false
	}

	if err := b.conn.WriteMessage(websocket.TextMessage, msgByte); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.NickName,
		}).Error("unable to join to chat")
	}
	return true
}

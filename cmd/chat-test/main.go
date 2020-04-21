package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/boletia/chat-test/pkg/bot"
	log "github.com/sirupsen/logrus"
)

const subDomain = "hola"

func main() {
	var numBots, numMsgs, minLatency, maxLatency int
	var wg sync.WaitGroup

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	flag.IntVar(&numBots, "bots", 1, "bots=<num_bots>")
	flag.IntVar(&numMsgs, "messages", 2, "messages=<num_of_messaves>")
	flag.IntVar(&minLatency, "min-latency", 10, "min-latency=<minimous_of_latency>")
	flag.IntVar(&maxLatency, "max-latency", 20, "max-latency=<minimous_of_latency>")

	flag.Parse()
	log.WithFields(log.Fields{
		"numBots":     numBots,
		"messages":    numMsgs,
		"min-latency": minLatency,
		"max-latency": maxLatency,
	}).Info("params")

	var quits []interface{}

	quit := make(chan bool)
	go gossiper(&wg, quit)
	quits = append(quits, quit)

	for i := 0; i < numBots; i++ {
		time.Sleep(10 * time.Millisecond)

		log.Infof("launching bot %d", i)

		quit := make(chan bool)
		conf := bot.Config{
			NickName:    fmt.Sprintf("bot-%d", i),
			SudDomain:   "noisescapes",
			NumMessages: numMsgs,
			MinDelay:    int64(minLatency),
			MaxDelay:    int64(maxLatency),
			Host:        "6vfdhz6o24.execute-api.us-east-1.amazonaws.com",
			Path:        "/beta",
			Schema:      "wss",
		}

		wg.Add(1)
		go launchBot(conf, &wg, quit)

		quits = append(quits, quit)
	}

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		select {
		case <-interrupt:
			log.Infof("sending quick message to channels")
			for _, quit := range quits {
				if tQuite, ok := quit.(chan bool); ok {
					go func() {
						tQuite <- true
					}()
				}
			}
		}
	}()

	wg.Wait()

	log.Info("finalizing")
}

func launchBot(conf bot.Config, wg *sync.WaitGroup, quit chan bool) {
	defer wg.Done()

	bot := bot.New(conf, quit)

	if bot.Connec() {
		if bot.JoinChat() {
			bot.WriteMessages()
		}
	}
}

func gossiper(wg *sync.WaitGroup, quit chan bool) {
	defer wg.Done()

	conf := bot.Config{
		NickName:  "gossiper",
		SudDomain: "noisescapes",
		Host:      "6vfdhz6o24.execute-api.us-east-1.amazonaws.com",
		Path:      "/beta",
		Schema:    "wss",
	}

	bot := bot.New(conf, quit)

	if bot.Connec() {
		if bot.JoinChat() {
			bot.ReadMessage()
		}
	}
}

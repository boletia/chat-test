package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/boletia/chat-test/pkg/bot"
	log "github.com/sirupsen/logrus"
)

var (
	numOfMessages   = 20
	numOfBots       = 100
	sleepSecondsMin = 5
	sleepSecondsMax = 10
)

var wg sync.WaitGroup

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	numBots := flag.Int("bots", 10, "number of bots")
	messagessPerBot := flag.Int("messages", 20, "messages per bot")
	minLatency := flag.Int("min-latency", 5, "min-latency")
	maxLatency := flag.Int("max-latency", 10, "max-latency")
	flag.Parse()

	/*
		wg.Add(1)
		go gossiper(&wg)
	*/

	log.WithFields(log.Fields{
		"numBots":     *numBots,
		"msg4bot":     *messagessPerBot,
		"min latency": *minLatency,
		"max latency": *maxLatency,
	}).Info("Launching bots")

	for i := 1; i <= *numBots; i++ {
		wg.Add(1)
		go doTest(*numBots, *messagessPerBot, *minLatency, *maxLatency, i, &wg)
	}

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		select {
		case <-interrupt:
			log.Infof("signal received, waiting %d secodns to finish\n", 20)
			time.Sleep(20 * time.Second)
			os.Exit(0)
		}

	}()

	wg.Wait()
	log.Info("Saliendo")
}

func doTest(bots int, message int, min int, max int, n int, wg *sync.WaitGroup) {
	defer wg.Done()

	botName := fmt.Sprintf("bot-%d", n)

	url := url.URL{
		Scheme: "wss",
		Host:   "6vfdhz6o24.execute-api.us-east-1.amazonaws.com",
		Path:   "/beta",
	}

	bot := bot.New(url, botName, false, "noisescapes")
	bot.NumMessages = message
	bot.SleepSecondsMin = min
	bot.SleepSecondsMax = max
	bot.Connect()

	bot.Do()
}

func gossiper(wg *sync.WaitGroup) {
	defer wg.Done()

	url := url.URL{
		Scheme: "wss",
		Host:   "6vfdhz6o24.execute-api.us-east-1.amazonaws.com",
		Path:   "/beta",
	}

	botName := "gossiper"

	bot := bot.New(url, botName, false, "noisescapes")
	defer bot.Disconnect()

	if bot.JoinChat() {
		log.Info("gossiper listening...")
		bot.Listen()
	}
}

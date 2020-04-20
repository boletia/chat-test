package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/boletia/chat-test/pkg/bot"
	log "github.com/sirupsen/logrus"
)

const (
	numOfMessages   = 10
	numOfBots       = 50
	sleepSecondsMin = 5
	sleepSecondsMax = 10
)

var wg sync.WaitGroup

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	/*
		wg.Add(1)
		go gossiper(&wg)
	*/

	for i := 1; i <= numOfBots; i++ {
		wg.Add(1)
		go doTest(i, &wg)
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

func doTest(n int, wg *sync.WaitGroup) {
	defer wg.Done()

	botName := fmt.Sprintf("bot-%d", n)

	url := url.URL{
		Scheme: "wss",
		Host:   "6vfdhz6o24.execute-api.us-east-1.amazonaws.com",
		Path:   "/beta",
	}

	bot := bot.New(url, botName, false, "noisescapes")
	bot.NumMessages = numOfMessages
	bot.SleepSecondsMin = sleepSecondsMin
	bot.SleepSecondsMax = sleepSecondsMax
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

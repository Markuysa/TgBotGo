package main

import (
	tgClient "TelegramBot/clients/telegram"
	tgConsumer "TelegramBot/consumer/event-consumer"
	tgProcessor "TelegramBot/events/telegram"
	"TelegramBot/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

// 5623430249:AAHwmZrj2m0AkD1Gsn0hh1rTEWG9hBIt5U8
func main() {

	eventsProcessor := tgProcessor.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath))
	log.Print("..service started")
	consumer := tgConsumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal()
	}

}
func mustHost() string {

	token := flag.String("telegram-host", "", "host for access to tgbot")

	flag.Parse()

	if *token == "" {
		log.Fatal()
	}
	return *token
}

func mustToken() string {

	token := flag.String("telegram-token", "", "token for access to tgbot")

	flag.Parse()

	if *token == "" {
		log.Fatal()
	}
	return *token
}

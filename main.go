package main

import (
	"TelegramBot/clients/telegram"
	"flag"
	"log"
)

func main() {

	telegramCLient := telegram.New(mustHost(), mustToken())

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

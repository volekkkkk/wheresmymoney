package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mymmrac/telego"
	"github.com/volekkkkk/wheresmymoney/internal/bank"
)

func main() {
	var client *bank.MonoClient
	telegramBot, err := telego.NewBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Fail initialize Telegram Bot: %s\n", err)
	}

	botUser, err := telegramBot.GetMe()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Print Bot information
	fmt.Printf("Bot user: %+v\n", botUser)

	client, err = bank.GetClientInfo()
	if err != nil {
		log.Fatalf("Cannot get client info: %s\n", err)
	}
	log.Printf("Got client data: %v\n", *client)

	var statements []bank.Statement
	from := time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	log.Printf("from: %v to: %v", from, to)

	statements, err = bank.GetStatement(os.Getenv("MONO_ACCOUNT_ID"), strconv.FormatInt(from.Unix(), 10), strconv.FormatInt(to.Unix(), 10))
	if err != nil {
		log.Fatalf("Cannot get statement")
	}
	log.Printf("Got statement: %v\n", statements)
}

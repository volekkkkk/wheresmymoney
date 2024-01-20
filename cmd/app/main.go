package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mymmrac/telego"
	"github.com/volekkkkk/wheresmymoney/internal/bank"
	"github.com/volekkkkk/wheresmymoney/internal/environment"
)

func initEnv(envFileName string) error {
	err := environment.LoadEnv(envFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Warning: %s file not found. Using default environment variables.\n", envFileName)
			return nil
		}
		return err
	}
	return nil
}

func main() {
	err := initEnv(".env")
	if err != nil {
		log.Fatalf("Fail to init env variables: %s\n", err)
	}

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

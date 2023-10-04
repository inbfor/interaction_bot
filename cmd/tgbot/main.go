package main

import (
	"encoding/json"
	"flag"
	"log"
	"transaction_bot/internal/db"
	"transaction_bot/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nats-io/nats.go"
)

func main() {
	apiKey := flag.String("apiKey", "", "Api Key for bot From BotFather")
	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(*apiKey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	dbConn, _ := db.Connect("users.db")
	// stateOf := make(map[string]config.State)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	nc, _ := nats.Connect(nats.DefaultURL)

	updates := bot.GetUpdatesChan(u)

	nc.Subscribe("event", func(msg *nats.Msg) {

		var event model.Event
		json.Unmarshal(msg.Data, &event)

		addresses, err := db.SelectUsers(event.Activity[0].FromAddress, event.Activity[0].ToAddress, dbConn)

		log.Println(addresses)

		if err != nil {
			log.Println(err)
		}

		for _, address := range addresses {
			msgSend := tgbotapi.NewMessage(address.Chat_id, event.Network)
			log.Println(event.Network)
			bot.Send(msgSend)
		}

	})

	for update := range updates {
		// TODO
		db.InsertIntoTable(update.Message.From.ID, update.Message.From.UserName, "0x74a5edd315951ccbcc91fc85f7d7ee67b60d71d8", dbConn)
		msg := tgbotapi.NewMessage(update.Message.From.ID, "TODO")
		bot.Send(msg)

	}

}

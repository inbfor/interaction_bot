package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"transaction_bot/internal/db"
	"transaction_bot/internal/model"
	"transaction_bot/internal/utils"

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

		log.Println(event)
		log.Println(addresses)

		if err != nil {
			log.Println(err)
		}

		for _, address := range addresses {
			if db.CheckAddr(event.Activity[0].FromAddress, dbConn) && db.CheckAddr(event.Activity[0].ToAddress, dbConn) {
				val := utils.HexToEth(event.Activity[0].RawContract.RawValue)

				msgSend := tgbotapi.NewMessage(address.Chat_id, fmt.Sprintf("Вы отправили %s eth с вами отслеживаемого кошелька %s на вами отслеживаемый кошелек %s.", val.Text('f', -1), event.Activity[0].FromAddress, event.Activity[0].ToAddress))
				bot.Send(msgSend)

			} else {
				val := utils.HexToEth(event.Activity[0].RawContract.RawValue)

				if event.Activity[0].FromAddress == address.Eth_address {
					msgSend := tgbotapi.NewMessage(address.Chat_id, fmt.Sprintf("Вы отправили %s eth на адрес %s.", val.Text('f', -1), event.Activity[0].ToAddress))
					bot.Send(msgSend)
				} else {
					msgSend := tgbotapi.NewMessage(address.Chat_id, fmt.Sprintf("Вам отправили %s eth с адреса %s.", val.Text('f', -1), event.Activity[0].FromAddress))
					bot.Send(msgSend)
				}
			}
		}

	})

	for update := range updates {

		message := update.Message.Text

		if db.CheckNumberAddr(update.Message.Chat.UserName, 10, dbConn) {
			if utils.CheckIfValidAddress(message) {
				if !db.CheckAddr(message, dbConn) {
					if err != nil {
						log.Printf("%s insert error", err)
					}

					channelJson := model.ChannelJson{
						Chat_id: update.Message.From.ID,
						TgNick:  update.Message.From.UserName,
						Address: message,
					}

					messageChannel, err := json.Marshal(&channelJson)

					if err != nil {
						log.Println(err)
					}

					msg := tgbotapi.NewMessage(update.Message.From.ID, "Вы теперь отслеживаете этот адрес!")
					nc.Publish("addresses", messageChannel)
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.From.ID, "Вы уже отслеживаете этот адрес!")
					bot.Send(msg)
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.From.ID, "Это не адрес Ethereum кошелька :(")
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.From.ID, "Вы больше не можете добавлять аддреса :(")
			bot.Send(msg)
		}

	}

}

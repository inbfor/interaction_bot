package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"transaction_bot/internal/db"
	"transaction_bot/internal/model"
	notify "transaction_bot/internal/webhooks"

	"github.com/nats-io/nats.go"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func main() {

	apiKeyAlchemy := os.Getenv("ALCHEMY_API_KEY")

	mux := http.NewServeMux()

	tun, err := runNgrok(context.Background())
	dbConn, _ := db.Connect("users.db")

	if err != nil {
		log.Println(err)
	}
	nc, _ := nats.Connect(nats.DefaultURL)

	nc.Subscribe("addresses", func(msg *nats.Msg) {

		var msgChannel model.ChannelJson

		json.Unmarshal(msg.Data, &msgChannel)

		jsn, _ := createWebHook(tun.URL(), msgChannel.Address, apiKeyAlchemy)

		log.Println(msgChannel)

		err := db.InsertIntoTable(msgChannel.Chat_id, msgChannel.TgNick, msgChannel.Address, jsn.Data.SigningKey, dbConn)

		if err != nil {
			log.Println(err)
		}
	})
	// Register handler for Alchemy Notify webhook events
	mux.Handle(
		// TODO: update to your own webhook path
		"/",
		// Middleware needed to validate the alchemy signature
		notify.NewAlchemyRequestHandlerMiddleware(handleWebhook, dbConn, nc),
	)

	// Listen to Alchemy Notify webhook events
	log.Printf("Example Alchemy Notify app listening at %s\n", "something")
	err = http.Serve(tun, mux)
	log.Fatal(err)
}

func handleWebhook(w http.ResponseWriter, req *http.Request, event *model.EventJson, nc *nats.Conn) {
	// Do stuff with with webhook event here!
	// var webHookEvent model.EventJson

	log.Println(event)

	nc.Publish("event", event.Event)
}

func runNgrok(ctx context.Context) (ngrok.Tunnel, error) {
	tun, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return nil, err
	}

	log.Println("tunnel created:", tun.URL())

	return tun, nil
}

func createWebHook(urlOfNgrok string, addr string, apiKey string) (CreatedWebHook, error) {

	var jsonCreatedWebhook CreatedWebHook

	url := "https://dashboard.alchemy.com/api/create-webhook"

	payload := strings.NewReader(
		fmt.Sprintf("{\"network\":\"ETH_SEPOLIA\",\"webhook_type\":\"ADDRESS_ACTIVITY\",\"addresses\":[\"%s\"],\"webhook_url\":\"%s\"}", addr, urlOfNgrok))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Alchemy-Token", apiKey)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	json.Unmarshal(body, &jsonCreatedWebhook)

	return jsonCreatedWebhook, nil
}

type CreatedWebHook struct {
	Data struct {
		ID          string `json:"id"`
		Network     string `json:"network"`
		WebhookType string `json:"webhook_type"`
		WebhookURL  string `json:"webhook_url"`
		IsActive    bool   `json:"is_active"`
		TimeCreated int64  `json:"time_created"`
		SigningKey  string `json:"signing_key"`
		Version     string `json:"version"`
	} `json:"data"`
}

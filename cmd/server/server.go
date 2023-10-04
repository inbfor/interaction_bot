package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"transaction_bot/internal/model"
	notify "transaction_bot/internal/webhooks"

	"github.com/nats-io/nats.go"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func main() {

	mux := http.NewServeMux()

	tun, err := runNgrok(context.Background())

	if err != nil {
		log.Println(err)
	}

	json, _ := createWebHook(fmt.Sprintf("%s/handle", tun.URL()), "0x74A5EDD315951cCBcC91FC85F7D7ee67b60d71D8")
	nc, _ := nats.Connect(nats.DefaultURL)
	// Register handler for Alchemy Notify webhook events
	mux.Handle(
		// TODO: update to your own webhook path
		"/handle",
		// Middleware needed to validate the alchemy signature
		notify.NewAlchemyRequestHandlerMiddleware(handleWebhook, json.Data.SigningKey, nc),
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
	// Be sure to respond with 200 when you successfully process the event
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

func createWebHook(urlOfNgrok string, addr string) (CreatedWebHook, error) {

	var jsonCreatedWebhook CreatedWebHook

	url := "https://dashboard.alchemy.com/api/create-webhook"

	payload := strings.NewReader(
		fmt.Sprintf("{\"network\":\"ETH_SEPOLIA\",\"webhook_type\":\"ADDRESS_ACTIVITY\",\"addresses\":[\"%s\"],\"webhook_url\":\"%s\"}", addr, urlOfNgrok))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Alchemy-Token", "Vw5p6rjCZccltsOHEpVnuimej-pbWdb1")
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

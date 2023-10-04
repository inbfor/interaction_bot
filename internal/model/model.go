package model

import (
	"encoding/json"
	"time"
)

type EventJson struct {
	WebhookID string          `json:"webhookId"`
	ID        string          `json:"id"`
	CreatedAt time.Time       `json:"createdAt"`
	Type      string          `json:"type"`
	Event     json.RawMessage `json:"event"`
}

type Event struct {
	Network  string `json:"network"`
	Activity []struct {
		Category      string `json:"category"`
		FromAddress   string `json:"fromAddress"`
		ToAddress     string `json:"toAddress"`
		Erc721TokenID string `json:"erc721TokenId"`
		RawContract   struct {
			RawValue string `json:"rawValue"`
			Address  string `json:"address"`
		} `json:"rawContract"`
		Log struct {
			Removed bool     `json:"removed"`
			Address string   `json:"address"`
			Data    string   `json:"data"`
			Topics  []string `json:"topics"`
		} `json:"log"`
	} `json:"activity"`
}

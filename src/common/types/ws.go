package types

import "encoding/json"

type WsMessage struct {
	Id      Uuid            `json:"id"`
	Action  string          `json:"action"`
	Token   string          `json:"token"`
	Payload json.RawMessage `json:"payload"`
}

type WsResponse struct {
	Id      Uuid   `json:"id"`
	Action  string `json:"action"`
	Payload any    `json:"payload"`
}

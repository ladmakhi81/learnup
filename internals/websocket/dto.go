package websocket

import "encoding/json"

type WsMessagePayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

func NewWsMessagePayload(event string, data json.RawMessage) *WsMessagePayload {
	return &WsMessagePayload{
		Event: event,
		Data:  data,
	}
}

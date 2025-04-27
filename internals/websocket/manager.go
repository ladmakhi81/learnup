package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

type WsManager struct {
	connections map[uint]*websocket.Conn
	mu          sync.Mutex
}

func NewWsManager() *WsManager {
	return &WsManager{
		connections: make(map[uint]*websocket.Conn),
	}
}

func (manager *WsManager) AddConnection(userID uint, conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.connections[userID] = conn
}

func (manager *WsManager) RemoveConnection(userID uint) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.connections, userID)
}

func (manager *WsManager) PublishMessage(userID uint, eventType string, data json.RawMessage) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	conn := manager.connections[userID]
	payload := NewWsMessagePayload(eventType, data)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, payloadBytes)
}

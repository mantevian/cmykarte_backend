package ws

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"mantevian.xyz/cmykarte_backend/common/types"
)

type Hub struct {
	mu          sync.RWMutex
	clients     map[string]map[types.Uuid]*Client // group -> client id -> client
	clientGroup map[types.Uuid]string             // client id -> group
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]map[string]*Client),
		clientGroup: make(map[string]string),
	}
}

func (h *Hub) GetClientsInGroup(group string) []*Client {
	h.mu.RLock()
	m := h.clients[group]
	h.mu.RUnlock()

	var clients []*Client
	for _, c := range m {
		clients = append(clients, c)
	}
	return clients
}

func (h *Hub) Broadcast(group string, action string, payload []byte) {
	h.mu.RLock()
	m := h.clients[group]
	h.mu.RUnlock()

	msg := types.WsMessage{
		Action:  action,
		Payload: payload,
	}
	msgBytes, _ := json.Marshal(msg)
	fmt.Println("broadcast >>", group, "/", action, ">>", string(msgBytes))
	for _, c := range m {
		_ = c.Conn.WriteMessage(websocket.TextMessage, msgBytes)
	}
}

func (h *Hub) AddClient(group string, c *Client) {
	h.mu.Lock()
	m, ok := h.clients[group]
	if !ok {
		m = make(map[string]*Client)
		h.clients[group] = m
	}
	m[c.Id] = c
	h.clientGroup[c.Id] = group
	h.mu.Unlock()
}

func (h *Hub) RemoveClient(id string) {
	h.mu.Lock()
	group, ok := h.clientGroup[id]
	if ok {
		if m, ok2 := h.clients[group]; ok2 {
			delete(m, id)
			if len(m) == 0 {
				delete(h.clients, group)
			}
		}
		delete(h.clientGroup, id)
	}
	h.mu.Unlock()
}

func (h *Hub) MoveClient(id string, group string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	curGroup := h.clientGroup[id]
	if curGroup == group {
		return
	}

	if m, ok := h.clients[curGroup]; ok {
		client := m[id]
		delete(m, id)
		if len(m) == 0 {
			delete(h.clients, curGroup)
		}

		dst, ok := h.clients[group]
		if !ok {
			dst = make(map[string]*Client)
			h.clients[group] = dst
		}
		dst[id] = client
		h.clientGroup[id] = group
	}
}

func (h *Hub) SetClientUserId(clientId types.Uuid, userId types.Uuid) {
	h.mu.Lock()
	g := h.clientGroup[clientId]
	c := h.clients[g][clientId]
	c.UserId = userId
	h.mu.Unlock()
}

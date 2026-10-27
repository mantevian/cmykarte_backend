package ws

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"mantevian.xyz/cmykarte_backend/common/types"
)

var Hub struct {
	mu          sync.RWMutex
	clients     map[string]map[types.Uuid]*Client // group -> client id -> client
	clientGroup map[types.Uuid]string             // client id -> group
	userClient  map[types.Uuid]*Client            // user id -> client
}

func GetClientsInGroup(group string) []*Client {
	Hub.mu.RLock()
	m := Hub.clients[group]
	Hub.mu.RUnlock()

	var clients []*Client
	for _, c := range m {
		clients = append(clients, c)
	}
	return clients
}

func SendToUser(userId types.Uuid, action string, payload any) {
	Hub.mu.RLock()
	client, ok := Hub.userClient[userId]

	if !ok {
		return
	}

	Hub.mu.RUnlock()

	msg := types.WsResponse{
		Action:  action,
		Payload: payload,
	}
	msgBytes, _ := json.Marshal(msg)
	fmt.Println("send to user >>", userId, "/", action, ">>", string(msgBytes))

	_ = client.Conn.WriteMessage(websocket.TextMessage, msgBytes)
}

func Broadcast(group string, action string, payload any) {
	Hub.mu.RLock()
	m := Hub.clients[group]
	Hub.mu.RUnlock()

	msg := types.WsResponse{
		Action:  action,
		Payload: payload,
	}
	msgBytes, _ := json.Marshal(msg)
	fmt.Println("broadcast >>", group, "/", action, ">>", string(msgBytes))
	for _, c := range m {
		_ = c.Conn.WriteMessage(websocket.TextMessage, msgBytes)
	}
}

func AddClient(group string, c *Client) {
	Hub.mu.Lock()
	m, ok := Hub.clients[group]
	if !ok {
		m = make(map[string]*Client)
		Hub.clients[group] = m
	}
	m[c.Id] = c
	Hub.clientGroup[c.Id] = group
	Hub.mu.Unlock()
}

func RemoveClient(id string) {
	Hub.mu.Lock()
	group, ok := Hub.clientGroup[id]
	if ok {
		if m, ok2 := Hub.clients[group]; ok2 {
			c := m[id]
			delete(Hub.userClient, c.UserId)
			delete(m, id)

			c.Conn.Close()

			if len(m) == 0 {
				delete(Hub.clients, group)
			}
		}
		delete(Hub.clientGroup, id)
	}

	Hub.mu.Unlock()
}

func MoveUser(userId types.Uuid, group string) {
	Hub.mu.Lock()
	defer Hub.mu.Unlock()

	client := Hub.userClient[userId]

	curGroup := Hub.clientGroup[client.Id]
	if curGroup == group {
		return
	}

	if m, ok := Hub.clients[curGroup]; ok {
		delete(m, client.Id)
		if len(m) == 0 {
			delete(Hub.clients, curGroup)
		}

		dst, ok := Hub.clients[group]
		if !ok {
			dst = make(map[string]*Client)
			Hub.clients[group] = dst
		}
		dst[client.Id] = client
		Hub.clientGroup[client.Id] = group
	}
}

func SetClientUserId(clientId types.Uuid, userId types.Uuid) {
	Hub.mu.Lock()
	g := Hub.clientGroup[clientId]
	c := Hub.clients[g][clientId]
	c.UserId = userId

	if _, ok := Hub.userClient[userId]; ok {

	}

	Hub.userClient[userId] = c
	Hub.mu.Unlock()
}

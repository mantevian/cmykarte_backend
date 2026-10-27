package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"uuid"

	"github.com/gorilla/websocket"
	"mantevian.xyz/cmykarte_backend/common/types"
	"mantevian.xyz/cmykarte_backend/functions"
	"mantevian.xyz/cmykarte_backend/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func genericResponse(w http.ResponseWriter, code int, resp types.GenericResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

func HandleWs(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			genericResponse(w, http.StatusInternalServerError, types.GenericResponseError("Server error"))
			return
		}

		id := uuid.New().String()

		client := ws.Client{
			Id:   id,
			Conn: conn,
		}

		ws.AddClient("", &client)

		defer ws.RemoveClient(id)

		var ctx = types.FunctionCallCtx{
			DB: db,
		}

		for {
			_, bytes, err := conn.ReadMessage()

			if err != nil {
				break
			}

			log.Println("incoming >>", string(bytes))

			var message types.WsMessage
			json.Unmarshal(bytes, &message)

			validateTokenRes := functions.ValidateToken(types.ValidateTokenRequest{Token: message.Token}, ctx)
			authSuccess := validateTokenRes.Success

			ws.SetClientUserId(client.Id, validateTokenRes.Claims.UserId)

			var res any

			switch message.Action {
			case "register":
				var req types.RegisterRequest
				json.Unmarshal(message.Payload, &req)
				res = functions.Register(req, ctx)
			case "login":
				var req types.LoginRequest
				json.Unmarshal(message.Payload, &req)
				loginRes := functions.Login(req, ctx)
				res = loginRes
				ws.SetClientUserId(client.Id, loginRes.UserId)
			case "validate_token":
				res = validateTokenRes
			}

			if authSuccess {
				ctx.UserId = validateTokenRes.Claims.UserId

				switch message.Action {
				case "queue":
					var req types.QueueRequest
					json.Unmarshal(message.Payload, &req)
					queueRes := functions.Queue(req, ctx)
					res = queueRes
				case "unqueue":
					var req types.QueueRequest
					json.Unmarshal(message.Payload, &req)
					queueRes := functions.Unqueue(req, ctx)
					res = queueRes
				}
			}

			response := types.WsResponse{
				Action:  message.Action,
				Id:      message.Id,
				Payload: res,
			}

			responseJson, _ := json.Marshal(response)
			log.Println("response >>", string(responseJson))

			if err := conn.WriteMessage(websocket.TextMessage, responseJson); err != nil {
				break
			}
		}

		// game.RemoveQueueEntry(ctx.UserId)
	}
}

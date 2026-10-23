package ws

import (
	"github.com/gorilla/websocket"
	"mantevian.xyz/cmykarte_backend/common/types"
)

type Client struct {
	Id     types.Uuid
	Conn   *websocket.Conn
	UserId types.Uuid
}

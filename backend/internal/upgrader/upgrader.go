package upgrader

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func CreateUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

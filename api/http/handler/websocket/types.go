package websocket

import (
	"github.com/hazik1024/portainer/api"
)

type webSocketRequestParams struct {
	ID       string
	nodeName string
	endpoint *portainer.Endpoint
}

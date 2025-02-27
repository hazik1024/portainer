package websocket

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	httperror "github.com/portainer/libhttp/error"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/security"
)

// Handler is the HTTP handler used to handle websocket operations.
type Handler struct {
	*mux.Router
	EndpointService    portainer.EndpointService
	SignatureService   portainer.DigitalSignatureService
	requestBouncer     *security.RequestBouncer
	connectionUpgrader websocket.Upgrader
}

// NewHandler creates a handler to manage websocket operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:             mux.NewRouter(),
		connectionUpgrader: websocket.Upgrader{},
		requestBouncer:     bouncer,
	}
	h.PathPrefix("/websocket/exec").Handler(
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.websocketExec)))
	h.PathPrefix("/websocket/attach").Handler(
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.websocketAttach)))
	return h
}

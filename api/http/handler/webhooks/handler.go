package webhooks

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	portainer "github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/docker"
	"github.com/hazik1024/portainer/api/http/security"
)

// Handler is the HTTP handler used to handle webhook operations.
type Handler struct {
	*mux.Router
	WebhookService      portainer.WebhookService
	EndpointService     portainer.EndpointService
	DockerClientFactory *docker.ClientFactory
}

// NewHandler creates a handler to manage settings operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/webhooks",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.webhookCreate))).Methods(http.MethodPost)
	h.Handle("/webhooks",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.webhookList))).Methods(http.MethodGet)
	h.Handle("/webhooks/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.webhookDelete))).Methods(http.MethodDelete)
	h.Handle("/webhooks/{token}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.webhookExecute))).Methods(http.MethodPost)
	return h
}

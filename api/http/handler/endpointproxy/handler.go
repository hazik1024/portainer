package endpointproxy

import (
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/proxy"
	"github.com/hazik1024/portainer/api/http/security"
)

// Handler is the HTTP handler used to proxy requests to external APIs.
type Handler struct {
	*mux.Router
	requestBouncer  *security.RequestBouncer
	EndpointService portainer.EndpointService
	ProxyManager    *proxy.Manager
}

// NewHandler creates a handler to proxy requests to external APIs.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
	}
	h.PathPrefix("/{id}/azure").Handler(
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.proxyRequestsToAzureAPI)))
	h.PathPrefix("/{id}/docker").Handler(
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.proxyRequestsToDockerAPI)))
	h.PathPrefix("/{id}/storidge").Handler(
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.proxyRequestsToStoridgeAPI)))
	return h
}

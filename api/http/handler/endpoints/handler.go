package endpoints

import (
	httperror "github.com/portainer/libhttp/error"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/proxy"
	"github.com/hazik1024/portainer/api/http/security"

	"net/http"

	"github.com/gorilla/mux"
)

const (
	// ErrEndpointManagementDisabled is an error raised when trying to access the endpoints management endpoints
	// when the server has been started with the --external-endpoints flag
	ErrEndpointManagementDisabled = portainer.Error("Endpoint management is disabled")
)

func hideFields(endpoint *portainer.Endpoint) {
	endpoint.AzureCredentials = portainer.AzureCredentials{}
}

// Handler is the HTTP handler used to handle endpoint operations.
type Handler struct {
	*mux.Router
	authorizeEndpointManagement bool
	requestBouncer              *security.RequestBouncer
	EndpointService             portainer.EndpointService
	EndpointGroupService        portainer.EndpointGroupService
	FileService                 portainer.FileService
	ProxyManager                *proxy.Manager
	Snapshotter                 portainer.Snapshotter
	JobService                  portainer.JobService
}

// NewHandler creates a handler to manage endpoint operations.
func NewHandler(bouncer *security.RequestBouncer, authorizeEndpointManagement bool) *Handler {
	h := &Handler{
		Router:                      mux.NewRouter(),
		authorizeEndpointManagement: authorizeEndpointManagement,
		requestBouncer:              bouncer,
	}

	h.Handle("/endpoints",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointCreate))).Methods(http.MethodPost)
	h.Handle("/endpoints/snapshot",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointSnapshots))).Methods(http.MethodPost)
	h.Handle("/endpoints",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointList))).Methods(http.MethodGet)
	h.Handle("/endpoints/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointInspect))).Methods(http.MethodGet)
	h.Handle("/endpoints/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointUpdate))).Methods(http.MethodPut)
	h.Handle("/endpoints/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointDelete))).Methods(http.MethodDelete)
	h.Handle("/endpoints/{id}/extensions",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointExtensionAdd))).Methods(http.MethodPost)
	h.Handle("/endpoints/{id}/extensions/{extensionType}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointExtensionRemove))).Methods(http.MethodDelete)
	h.Handle("/endpoints/{id}/job",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointJob))).Methods(http.MethodPost)
	h.Handle("/endpoints/{id}/snapshot",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointSnapshot))).Methods(http.MethodPost)
	return h
}

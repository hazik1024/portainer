package endpointgroups

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/security"
)

// Handler is the HTTP handler used to handle endpoint group operations.
type Handler struct {
	*mux.Router
	EndpointService      portainer.EndpointService
	EndpointGroupService portainer.EndpointGroupService
}

// NewHandler creates a handler to manage endpoint group operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/endpoint_groups",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointGroupCreate))).Methods(http.MethodPost)
	h.Handle("/endpoint_groups",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointGroupList))).Methods(http.MethodGet)
	h.Handle("/endpoint_groups/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointGroupInspect))).Methods(http.MethodGet)
	h.Handle("/endpoint_groups/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointGroupUpdate))).Methods(http.MethodPut)
	h.Handle("/endpoint_groups/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.endpointGroupDelete))).Methods(http.MethodDelete)

	return h
}

func (handler *Handler) checkForGroupUnassignment(endpoint portainer.Endpoint, associatedEndpoints []portainer.EndpointID) error {
	for _, id := range associatedEndpoints {
		if id == endpoint.ID {
			return nil
		}
	}

	endpoint.GroupID = portainer.EndpointGroupID(1)
	return handler.EndpointService.UpdateEndpoint(endpoint.ID, &endpoint)
}

func (handler *Handler) checkForGroupAssignment(endpoint portainer.Endpoint, groupID portainer.EndpointGroupID, associatedEndpoints []portainer.EndpointID) error {
	for _, id := range associatedEndpoints {

		if id == endpoint.ID {
			endpoint.GroupID = groupID
			return handler.EndpointService.UpdateEndpoint(endpoint.ID, &endpoint)
		}
	}
	return nil
}

func (handler *Handler) updateEndpointGroup(endpoint portainer.Endpoint, groupID portainer.EndpointGroupID, associatedEndpoints []portainer.EndpointID) error {
	if endpoint.GroupID == groupID {
		return handler.checkForGroupUnassignment(endpoint, associatedEndpoints)
	} else if endpoint.GroupID == portainer.EndpointGroupID(1) {
		return handler.checkForGroupAssignment(endpoint, groupID, associatedEndpoints)
	}
	return nil
}

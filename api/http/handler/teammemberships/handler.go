package teammemberships

import (
	httperror "github.com/portainer/libhttp/error"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/security"

	"net/http"

	"github.com/gorilla/mux"
)

// Handler is the HTTP handler used to handle team membership operations.
type Handler struct {
	*mux.Router
	TeamMembershipService  portainer.TeamMembershipService
	ResourceControlService portainer.ResourceControlService
}

// NewHandler creates a handler to manage team membership operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/team_memberships",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.teamMembershipCreate))).Methods(http.MethodPost)
	h.Handle("/team_memberships",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.teamMembershipList))).Methods(http.MethodGet)
	h.Handle("/team_memberships/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.teamMembershipUpdate))).Methods(http.MethodPut)
	h.Handle("/team_memberships/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.teamMembershipDelete))).Methods(http.MethodDelete)

	return h
}

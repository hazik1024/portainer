package roles

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/security"
)

// Handler is the HTTP handler used to handle role operations.
type Handler struct {
	*mux.Router
	RoleService portainer.RoleService
}

// NewHandler creates a handler to manage role operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/roles",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.roleList))).Methods(http.MethodGet)

	return h
}

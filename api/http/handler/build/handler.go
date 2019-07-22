package build

import (
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/http/proxy"
	"github.com/portainer/portainer/api/http/security"

	"net/http"

	"github.com/gorilla/mux"
)

// Handler is the HTTP handler used to handle build operations.
type Handler struct {
	*mux.Router
}

// NewHandler creates a handler to manage build operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter()
	}
	h.Handle("/build/{certificate:(?:ca|cert|key)}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.build))).Methods(http.MethodPost)
	return h
}

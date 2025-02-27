package motd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hazik1024/portainer/api/http/security"
)

// Handler is the HTTP handler used to handle MOTD operations.
type Handler struct {
	*mux.Router
}

// NewHandler returns a new Handler
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/motd",
		bouncer.AuthorizedAccess(http.HandlerFunc(h.motd))).Methods(http.MethodGet)

	return h
}

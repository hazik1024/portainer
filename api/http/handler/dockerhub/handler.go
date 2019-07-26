package dockerhub

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/security"
)

func hideFields(dockerHub *portainer.DockerHub) {
	dockerHub.Password = ""
}

// Handler is the HTTP handler used to handle DockerHub operations.
type Handler struct {
	*mux.Router
	DockerHubService portainer.DockerHubService
}

// NewHandler creates a handler to manage Dockerhub operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/dockerhub",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.dockerhubInspect))).Methods(http.MethodGet)
	h.Handle("/dockerhub",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.dockerhubUpdate))).Methods(http.MethodPut)

	return h
}

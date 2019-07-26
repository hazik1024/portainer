package build

import (
	"log"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"

	// "github.com/portainer/libhttp/response"
	"github.com/hazik1024/portainer/api/http/security"
)

// Handler 编译镜像
type Handler struct {
	*mux.Router
	requestBouncer *security.RequestBouncer
	service        *Service
}

// NewHandler 返回新的Handler
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
		service:        &Service{},
	}
	h.PathPrefix("/build").Handler(bouncer.RestrictedAccess(httperror.LoggerHandler(h.proxyBuild)))
	h.PathPrefix("/build/history").Handler(bouncer.RestrictedAccess(httperror.LoggerHandler(h.proxyBuildHistory)))
	return h
}

func (handler *Handler) proxyBuild() *httperror.HandlerError {
	log.Fatal("test_proxyBuild")
	// return response.JSON(w, filteredRegistries)
	return "{'data':[]}"
}

func (handler *Handler) proxyBuildHistory() *httperror.HandlerError {
	log.Fatal("test_proxyBuildHistory")
	// return response.JSON(w, filteredRegistries)
	return "{'data':[]}"
}

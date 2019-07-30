package build

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hazik1024/portainer/api/http/security"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

// Handler 编译镜像
type Handler struct {
	*mux.Router
	requestBouncer *security.RequestBouncer
	service        *Service
}

// Resp 响应格式
type Resp struct {
	ID   int    `json:"id"`
	Type int    `json:"type"`
	Data string `json:"data"`
}

// NewHandler 返回新的Handler
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
		service:        &Service{},
	}
	// h.PathPrefix("/build").Handler(httperror.LoggerHandler(h.proxyBuild))
	// h.PathPrefix("/build/history").Handler(httperror.LoggerHandler(h.proxyBuildHistory))
	h.PathPrefix("/build").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBuild)))
	h.PathPrefix("/build/history").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBuildHistory)))
	return h
}

func (handler *Handler) proxyBuild(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Fatal("test_proxyBuild")
	resp := &Resp{
		ID:   1,
		Type: 2,
		Data: "proxyBuild",
	}
	return response.JSON(w, resp)
}

func (handler *Handler) proxyBuildHistory(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Fatal("test_proxyBuildHistory")
	resp := &Resp{
		ID:   1,
		Type: 2,
		Data: "proxyBuildHistory",
	}
	return response.JSON(w, resp)
}

package build

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hazik1024/portainer/api/http/security"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

type (
	// Resp 响应格式
	Resp struct {
		ID   int    `json:"id"`
		Type int    `json:"type"`
		Data string `json:"data"`
	}

	// Handler 编译镜像
	Handler struct {
		*mux.Router
		requestBouncer *security.RequestBouncer
		service        *Service
	}
)

// NewHandler 返回新的Handler
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
		service:        &Service{},
	}
	// h.Handle("/build", bouncer.AuthorizedAccess(httperror.LoggerHandler(h.proxyBuild))).Methods(http.MethodPost)
	// h.Handle("/build/history", bouncer.AuthorizedAccess(httperror.LoggerHandler(h.proxyBuildHistory))).Methods(http.MethodPost)
	h.PathPrefix("/build").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBuild))).Methods(http.MethodPost)
	h.PathPrefix("/build/history").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBuildHistory))).Methods(http.MethodPost)
	return h
}

func (handler *Handler) proxyBuild(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Println("test_proxyBuild 1111")
	resp := &Resp{
		ID:   1,
		Type: 2,
		Data: "proxyBuild",
	}
	log.Println("test_proxyBuild 2222")
	return response.JSON(w, resp)
}

func (handler *Handler) proxyBuildHistory(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Println("proxyBuildHistory 1111")
	resp := &Resp{
		ID:   1,
		Type: 2,
		Data: "proxyBuildHistory",
	}
	log.Println("proxyBuildHistory 2222")
	return response.JSON(w, resp)
}

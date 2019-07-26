package stack

import (
	"log"
	"net/http"

	"github.com/hazik1024/portainer/api/http/security"
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
	h.PathPrefix("/stack").Handler(bouncer.RestrictedAccess(httperror.LoggerHandler(h.proxyBackup)))
	h.PathPrefix("/stack/history").Handler(bouncer.RestrictedAccess(httperror.LoggerHandler(h.proxyBackupHistory)))
	return h
}

func (handler *Handler) proxyBackup(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Fatal("test_proxyBuild")
	resp := Resp{
		ID:   1,
		Type: 2,
		Data: "proxyBackup",
	}
	return response.JSON(w, resp)
	// return "{'data':[]}"
}

func (handler *Handler) proxyBackupHistory(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Fatal("test_proxyBuildHistory")
	resp := Resp{
		ID:   1,
		Type: 2,
		Data: "proxyBackupHistory",
	}
	return response.JSON(w, resp)
	// return "{'data':[]}"
}

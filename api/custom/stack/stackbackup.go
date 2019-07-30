package stack

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
	service        *BackupService
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
		service:        &BackupService{},
	}
	// h.PathPrefix("/stack").Handler(httperror.LoggerHandler(h.proxyBackup))
	// h.PathPrefix("/stack/history").Handler(httperror.LoggerHandler(h.proxyBackupHistory))
	h.Handle("/stack", bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBackup))).Methods(http.MethodPost)
	h.Handle("/stack/history", bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBackupHistory))).Methods(http.MethodPost)
	// h.PathPrefix("/stack").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBackup)))
	// h.PathPrefix("/stack/history").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBackupHistory)))
	return h
}

func (handler *Handler) proxyBackup(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Fatal("test_proxyBackup")
	resp := &Resp{
		ID:   1,
		Type: 2,
		Data: "proxyBackup",
	}
	return response.JSON(w, resp)
	// return "{'data':[]}"
}

func (handler *Handler) proxyBackupHistory(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Fatal("test_proxyBackupHistory")
	resp := &Resp{
		ID:   1,
		Type: 2,
		Data: "proxyBackupHistory",
	}
	return response.JSON(w, resp)
	// return "{'data':[]}"
}

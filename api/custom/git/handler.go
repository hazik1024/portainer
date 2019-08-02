package git

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hazik1024/portainer/api/http/security"
	httperror "github.com/portainer/libhttp/error"
)

type (
	// Payload 请求结构 (Type 类型   1 - GitLab  2 - GitHub)
	Payload struct {
		GitID    int    `json:"gitid"`
		GitName  string `json:"gitname"`
		UserName string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Project  string `json:"project"`
		Branch   string `json:"branch"`
		Type     int    `json:"type"`
	}

	// RespGit 响应结构
	RespGit struct {
		GitID    int    `json:"gitid"`
		GitName  string `json:"gitname"`
		UserName string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Project  string `json:"project"`
		Branch   string `json:"branch"`
		Type     int    `json:"type"`
		AddTime  string `json:"addtime"`
		LastTime string `json:"lasttime"`
	}

	// Handler 编译镜像
	Handler struct {
		*mux.Router
		requestBouncer *security.RequestBouncer
		service        *Service
	}
)

// NewHandler 返回新的Handler
func NewHandler(bouncer *security.RequestBouncer, service *Service) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
		service:        service,
	}
	h.PathPrefix("/git").Handler(bouncer.AuthorizedAccess(httperror.LoggerHandler(h.gitCreate))).Methods(http.MethodPost)
	h.PathPrefix("/git").Handler(bouncer.AuthorizedAccess(httperror.LoggerHandler(h.getGits))).Methods(http.MethodGet)
	h.PathPrefix("/git/{id}").Handler(bouncer.AuthorizedAccess(httperror.LoggerHandler(h.getGit))).Methods(http.MethodGet)
	h.PathPrefix("/git/{id}").Handler(bouncer.AuthorizedAccess(httperror.LoggerHandler(h.updateGit))).Methods(http.MethodPut)
	h.PathPrefix("/git/{id}").Handler(bouncer.AuthorizedAccess(httperror.LoggerHandler(h.gitDelete))).Methods(http.MethodDelete)
	return h
}

func (handler *Handler) getGits(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	return nil
}

func (handler *Handler) getGit(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	return nil
}

func (handler *Handler) updateGit(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	return nil
}

func (handler *Handler) deleteGit(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	return nil
}

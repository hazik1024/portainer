package build

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hazik1024/portainer/api/http/security"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

type (
	// CustomBuildResponseID  CustomResp ID
	CustomBuildResponseID int
	// CustomBuildResponseType CustomResp Type
	CustomBuildResponseType int

	// CustomBuildResponse 响应格式
	CustomBuildResponse struct {
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
	// h.PathPrefix("/build").Handler(httperror.LoggerHandler(h.proxyBuild))
	// h.PathPrefix("/build/history").Handler(httperror.LoggerHandler(h.proxyBuildHistory))
	// h.Handle("/build", bouncer.AuthorizedAccess(httperror.LoggerHandler(h.proxyBuild))).Methods(http.MethodPost)
	// h.Handle("/build/history", bouncer.AuthorizedAccess(httperror.LoggerHandler(h.proxyBuildHistory))).Methods(http.MethodPost)
	h.PathPrefix("/build").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBuild))).Methods(http.MethodPost)
	h.PathPrefix("/build/history").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBuildHistory))).Methods(http.MethodPost)
	return h
}

func (handler *Handler) proxyBuild(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Fatal("test_proxyBuild aaaa")
	return response.Empty(w)
	// jsonStr := `{"id": 1,"type":2,"data":"proxyBuild"}`
	// log.Fatal("test_proxyBuild222")
	// var resp *CustomBuildResponse
	// log.Fatal("test_proxyBuild3333")
	// err := json.Unmarshal([]byte(jsonStr), resp)
	// log.Fatal("test_proxyBuild4444")
	// if err != nil {
	// 	log.Fatal("test_proxyBuild5555")
	// 	return &httperror.HandlerError{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "parse error",
	// 		Err:        err,
	// 	}
	// }
	// log.Fatal("test_proxyBuild6666")
	// return response.JSON(w, resp)
}

func (handler *Handler) proxyBuildHistory(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	log.Fatal("test_proxyBuild111")
	var resp *CustomBuildResponse
	log.Fatal("test_proxyBuild222")
	jsonStr := `{"id": 1,"type":2,"data":"proxyBuild"}`
	err := json.Unmarshal([]byte(jsonStr), &resp)
	log.Fatal("test_proxyBuild3333")
	if err != nil {
		log.Fatal("test_proxyBuild4444")
		return &httperror.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Message:    "parse error",
			Err:        err,
		}
	}
	log.Fatal("test_proxyBuild5555")
	return response.JSON(w, resp)
}

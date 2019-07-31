package build

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hazik1024/portainer/api/custom/mysqldb"
	"github.com/hazik1024/portainer/api/http/security"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

type (
	// Resp 响应结构
	Resp struct {
		ID   int    `json:"id"`
		Type int    `json:"type"`
		Data string `json:"data"`
	}

	// ReqPayload 请求结构
	reqPayload struct {
		//
		GitPath      string `json:"gitpath"`
		GitBranch    string `json:"gitbranch"`
		ImageName    string `json:"imagename"`
		RegistryPath string `json:"registrypath"`
		RegistryName string `json:"registryname"`
	}

	// Handler 编译镜像
	Handler struct {
		*mux.Router
		requestBouncer *security.RequestBouncer
		service        *Service
		mysqlDb        *mysqldb.MySQLDb
	}
)

// Validate 请求参数校验
func (payload *reqPayload) Validate(r *http.Request) error {
	// if govalidator.IsNull(payload.Username) || govalidator.Contains(payload.Username, " ") {
	// 	return portainer.Error("Invalid username. Must not contain any whitespace")
	// }

	// if payload.Role != 1 && payload.Role != 2 {
	// 	return portainer.Error("Invalid role value. Value must be one of: 1 (administrator) or 2 (regular user)")
	// }
	return nil
}

// NewHandler 返回新的Handler
func NewHandler(bouncer *security.RequestBouncer, service *Service) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
		service:        service,
	}
	// h.Handle("/build", bouncer.AuthorizedAccess(httperror.LoggerHandler(h.proxyBuild))).Methods(http.MethodPost)
	// h.Handle("/build/history", bouncer.AuthorizedAccess(httperror.LoggerHandler(h.proxyBuildHistory))).Methods(http.MethodPost)
	h.PathPrefix("/build").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBuild))).Methods(http.MethodPost)
	h.PathPrefix("/build/history").Handler(bouncer.PublicAccess(httperror.LoggerHandler(h.proxyBuildHistory))).Methods(http.MethodPost)
	return h
}

func (handler *Handler) proxyBuild(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var req reqPayload
	err := request.DecodeAndValidateJSONPayload(r, &req)
	if err != nil {
		return &httperror.HandlerError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
			Err:        err,
		}
	}
	handler.service.BuildAndPushImage(req)
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

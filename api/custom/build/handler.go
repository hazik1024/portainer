package build

import (
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	portainer "github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/security"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

type (
	// Resp 响应结构
	Resp struct {
		OutInfos [6]string `json:"outinfos"`
	}

	// ReqPayload 请求结构
	reqPayload struct {
		//
		RegistryID int    `json:"registryid"`
		GitPath    string `json:"gitpath"`
		GitBranch  string `json:"gitbranch"`
		GitUser    string `json:"gituser"`
		GitPwd     string `json:"gitpwd"`
		ImageName  string `json:"imagename"`
		ImageTag   string `json:"imagetag"`
	}

	// Handler 编译镜像
	Handler struct {
		*mux.Router
		requestBouncer *security.RequestBouncer
		service        *Service
	}
)

// Validate 请求参数校验
func (payload *reqPayload) Validate(r *http.Request) error {
	if payload.RegistryID < 1 {
		return portainer.Error("Invalid registryid. Must greater than 0")
	}
	if govalidator.IsNull(payload.GitPath) || govalidator.Contains(payload.GitPath, " ") {
		return portainer.Error("Invalid gitpath. Must not contain any whitespace")
	}

	if govalidator.IsNull(payload.GitBranch) || govalidator.Contains(payload.GitBranch, " ") {
		return portainer.Error("Invalid gitbranch. Must not contain any whitespace")
	}

	if govalidator.IsNull(payload.GitUser) || govalidator.Contains(payload.GitUser, " ") {
		return portainer.Error("Invalid gituser. Must not contain any whitespace")
	}

	if govalidator.IsNull(payload.GitPwd) || govalidator.Contains(payload.GitPwd, " ") {
		return portainer.Error("Invalid gitpwd. Must not contain any whitespace")
	}

	if govalidator.IsNull(payload.ImageName) || govalidator.Contains(payload.ImageName, " ") {
		return portainer.Error("Invalid imagename. Must not contain any whitespace")
	}

	if govalidator.IsNull(payload.ImageTag) || govalidator.Contains(payload.ImageTag, " ") {
		return portainer.Error("Invalid imagetag. Must not contain any whitespace")
	}

	if payload.GitBranch == "" {
		payload.GitBranch = "latest"
	}
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
	outs, handlerErr := handler.service.buildAndPushImage(req)
	if handlerErr != nil {
		return handlerErr
	}
	log.Println(outs)
	resp := &Resp{
		OutInfos: outs,
	}
	return response.JSON(w, resp)
}

func (handler *Handler) proxyBuildHistory(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	return response.Empty(w)
}

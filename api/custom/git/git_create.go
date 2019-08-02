package git

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	portainer "github.com/hazik1024/portainer/api"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

// Validate 请求参数校验
func (req *Payload) Validate(r *http.Request) error {
	if govalidator.IsNull(req.GitName) || govalidator.Contains(req.GitName, " ") {
		return portainer.Error("Invalid gitname. Must not contain any whitespace")
	}

	if govalidator.IsNull(req.UserName) || govalidator.Contains(req.UserName, " ") {
		return portainer.Error("Invalid password. Must not contain any whitespace")
	}

	if govalidator.IsNull(req.Password) || govalidator.Contains(req.Password, " ") {
		return portainer.Error("Invalid password. Must not contain any whitespace")
	}

	if govalidator.IsNull(req.Address) || govalidator.Contains(req.Address, " ") {
		return portainer.Error("Invalid address. Must not contain any whitespace")
	}

	if govalidator.IsNull(req.Project) || govalidator.Contains(req.Project, " ") {
		return portainer.Error("Invalid project. Must not contain any whitespace")
	}

	if govalidator.IsNull(req.Branch) || govalidator.Contains(req.Branch, " ") {
		return portainer.Error("Invalid branch. Must not contain any whitespace")
	}

	if req.Type < 1 {
		req.Type = 1
	}
	return nil
}

func (handler *Handler) gitCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload Payload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{
			StatusCode: http.StatusBadRequest,
			Message:    "参数错误",
			Err:        err,
		}
	}

	gitInfo, err := handler.service.getGitByAddress(payload.Address)
	if err != nil {
		return &httperror.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Message:    "地址验证失败",
			Err:        err,
		}
	}
	if gitInfo.GitID >= 0 {
		return &httperror.HandlerError{
			StatusCode: http.StatusBadRequest,
			Message:    "Git地址已存在",
			Err:        err,
		}
	}

	gitID, err := handler.service.createGit(payload)
	if err != nil {
		return &httperror.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Message:    "保存Git信息失败",
			Err:        err,
		}
	}
	resp := &RespGit{
		GitID: gitID,
	}
	return response.JSON(w, resp)
}

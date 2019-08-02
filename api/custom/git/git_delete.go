package git

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

func (handler *Handler) gitDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	gitID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	_, err = handler.service.deleteGit(gitID)
	if err != nil {
		return &httperror.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Message:    "删除Git信息失败",
			Err:        err,
		}
	}
	return response.Empty(w)
}

package build

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	"github.com/portainer/portainer/api"
)

func (handler *Handler) build(r *http.Request) *httperror.HandlerError {
	clone();
	build_image();
	push();
	return response.JSON(w, [])
}

func clone() {
	
}

func build_image() {

}

func push() {

}
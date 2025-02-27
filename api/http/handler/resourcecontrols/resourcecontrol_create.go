package resourcecontrols

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/security"
)

type resourceControlCreatePayload struct {
	ResourceID     string
	Type           string
	Public         bool
	Users          []int
	Teams          []int
	SubResourceIDs []string
}

func (payload *resourceControlCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.ResourceID) {
		return portainer.Error("Invalid resource identifier")
	}

	if govalidator.IsNull(payload.Type) {
		return portainer.Error("Invalid type")
	}

	if len(payload.Users) == 0 && len(payload.Teams) == 0 && !payload.Public {
		return portainer.Error("Invalid resource control declaration. Must specify Users, Teams or Public")
	}
	return nil
}

// POST request on /api/resource_controls
func (handler *Handler) resourceControlCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload resourceControlCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	var resourceControlType portainer.ResourceControlType
	switch payload.Type {
	case "container":
		resourceControlType = portainer.ContainerResourceControl
	case "service":
		resourceControlType = portainer.ServiceResourceControl
	case "volume":
		resourceControlType = portainer.VolumeResourceControl
	case "network":
		resourceControlType = portainer.NetworkResourceControl
	case "secret":
		resourceControlType = portainer.SecretResourceControl
	case "stack":
		resourceControlType = portainer.StackResourceControl
	case "config":
		resourceControlType = portainer.ConfigResourceControl
	default:
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid type value. Value must be one of: container, service, volume, network, secret, stack or config", portainer.ErrInvalidResourceControlType}
	}

	rc, err := handler.ResourceControlService.ResourceControlByResourceID(payload.ResourceID)
	if err != nil && err != portainer.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve resource controls from the database", err}
	}
	if rc != nil {
		return &httperror.HandlerError{http.StatusConflict, "A resource control is already associated to this resource", portainer.ErrResourceControlAlreadyExists}
	}

	var userAccesses = make([]portainer.UserResourceAccess, 0)
	for _, v := range payload.Users {
		userAccess := portainer.UserResourceAccess{
			UserID:      portainer.UserID(v),
			AccessLevel: portainer.ReadWriteAccessLevel,
		}
		userAccesses = append(userAccesses, userAccess)
	}

	var teamAccesses = make([]portainer.TeamResourceAccess, 0)
	for _, v := range payload.Teams {
		teamAccess := portainer.TeamResourceAccess{
			TeamID:      portainer.TeamID(v),
			AccessLevel: portainer.ReadWriteAccessLevel,
		}
		teamAccesses = append(teamAccesses, teamAccess)
	}

	resourceControl := portainer.ResourceControl{
		ResourceID:     payload.ResourceID,
		SubResourceIDs: payload.SubResourceIDs,
		Type:           resourceControlType,
		Public:         payload.Public,
		UserAccesses:   userAccesses,
		TeamAccesses:   teamAccesses,
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	if !security.AuthorizedResourceControlCreation(&resourceControl, securityContext) {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to create a resource control for the specified resource", portainer.ErrResourceAccessDenied}
	}

	err = handler.ResourceControlService.CreateResourceControl(&resourceControl)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the resource control inside the database", err}
	}

	return response.JSON(w, resourceControl)
}

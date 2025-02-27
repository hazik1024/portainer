package registries

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	"github.com/hazik1024/portainer/api"
)

type registryUpdatePayload struct {
	Name               string
	URL                string
	Authentication     bool
	Username           string
	Password           string
	UserAccessPolicies portainer.UserAccessPolicies
	TeamAccessPolicies portainer.TeamAccessPolicies
}

func (payload *registryUpdatePayload) Validate(r *http.Request) error {
	if payload.Authentication && (govalidator.IsNull(payload.Username) || govalidator.IsNull(payload.Password)) {
		return portainer.Error("Invalid credentials. Username and password must be specified when authentication is enabled")
	}
	return nil
}

// PUT request on /api/registries/:id
func (handler *Handler) registryUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	registryID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid registry identifier route variable", err}
	}

	var payload registryUpdatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	registry, err := handler.RegistryService.Registry(portainer.RegistryID(registryID))
	if err == portainer.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a registry with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a registry with the specified identifier inside the database", err}
	}

	registries, err := handler.RegistryService.Registries()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve registries from the database", err}
	}
	for _, r := range registries {
		if r.URL == payload.URL && r.ID != registry.ID {
			return &httperror.HandlerError{http.StatusConflict, "Another registry with the same URL already exists", portainer.ErrRegistryAlreadyExists}
		}
	}

	if payload.Name != "" {
		registry.Name = payload.Name
	}

	if payload.URL != "" {
		registry.URL = payload.URL
	}

	if payload.Authentication {
		registry.Authentication = true
		registry.Username = payload.Username
		registry.Password = payload.Password
	} else {
		registry.Authentication = false
		registry.Username = ""
		registry.Password = ""
	}

	if payload.UserAccessPolicies != nil {
		registry.UserAccessPolicies = payload.UserAccessPolicies
	}

	if payload.TeamAccessPolicies != nil {
		registry.TeamAccessPolicies = payload.TeamAccessPolicies
	}

	err = handler.RegistryService.UpdateRegistry(registry.ID, registry)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist registry changes inside the database", err}
	}

	return response.JSON(w, registry)
}

package webhooks

import (
	"context"
	"net/http"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	"github.com/hazik1024/portainer/api"
)

// Acts on a passed in token UUID to restart the docker service
func (handler *Handler) webhookExecute(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {

	webhookToken, err := request.RetrieveRouteVariableValue(r, "token")

	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Invalid service id parameter", err}
	}

	webhook, err := handler.WebhookService.WebhookByToken(webhookToken)

	if err == portainer.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a webhook with this token", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve webhook from the database", err}
	}

	resourceID := webhook.ResourceID
	endpointID := webhook.EndpointID
	webhookType := webhook.WebhookType

	endpoint, err := handler.EndpointService.Endpoint(portainer.EndpointID(endpointID))
	if err == portainer.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an endpoint with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an endpoint with the specified identifier inside the database", err}
	}

	imageTag, _ := request.RetrieveQueryParameter(r, "tag", true)

	switch webhookType {
	case portainer.ServiceWebhook:
		return handler.executeServiceWebhook(w, endpoint, resourceID, imageTag)
	default:
		return &httperror.HandlerError{http.StatusInternalServerError, "Unsupported webhook type", portainer.ErrUnsupportedWebhookType}
	}
}

func (handler *Handler) executeServiceWebhook(w http.ResponseWriter, endpoint *portainer.Endpoint, resourceID string, imageTag string) *httperror.HandlerError {
	dockerClient, err := handler.DockerClientFactory.CreateClient(endpoint, "")
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Error creating docker client", err}
	}
	defer dockerClient.Close()

	service, _, err := dockerClient.ServiceInspectWithRaw(context.Background(), resourceID, dockertypes.ServiceInspectOptions{InsertDefaults: true})
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Error looking up service", err}
	}

	service.Spec.TaskTemplate.ForceUpdate++

	if imageTag != "" {
		service.Spec.TaskTemplate.ContainerSpec.Image = strings.Split(service.Spec.TaskTemplate.ContainerSpec.Image, ":")[0] + ":" + imageTag
	} else {
		service.Spec.TaskTemplate.ContainerSpec.Image = strings.Split(service.Spec.TaskTemplate.ContainerSpec.Image, "@sha")[0]
	}

	_, err = dockerClient.ServiceUpdate(context.Background(), resourceID, service.Version, service.Spec, dockertypes.ServiceUpdateOptions{QueryRegistry: true})

	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Error updating service", err}
	}
	return response.Empty(w)
}

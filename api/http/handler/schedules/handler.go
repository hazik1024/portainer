package schedules

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/http/security"
)

// Handler is the HTTP handler used to handle schedule operations.
type Handler struct {
	*mux.Router
	ScheduleService portainer.ScheduleService
	EndpointService portainer.EndpointService
	SettingsService portainer.SettingsService
	FileService     portainer.FileService
	JobService      portainer.JobService
	JobScheduler    portainer.JobScheduler
}

// NewHandler creates a handler to manage schedule operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	h.Handle("/schedules",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.scheduleList))).Methods(http.MethodGet)
	h.Handle("/schedules",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.scheduleCreate))).Methods(http.MethodPost)
	h.Handle("/schedules/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.scheduleInspect))).Methods(http.MethodGet)
	h.Handle("/schedules/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.scheduleUpdate))).Methods(http.MethodPut)
	h.Handle("/schedules/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.scheduleDelete))).Methods(http.MethodDelete)
	h.Handle("/schedules/{id}/file",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.scheduleFile))).Methods(http.MethodGet)
	h.Handle("/schedules/{id}/tasks",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.scheduleTasks))).Methods(http.MethodGet)
	return h
}

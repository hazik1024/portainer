package schedules

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/cron"
)

type scheduleUpdatePayload struct {
	Name           *string
	Image          *string
	CronExpression *string
	Recurring      *bool
	Endpoints      []portainer.EndpointID
	FileContent    *string
	RetryCount     *int
	RetryInterval  *int
}

func (payload *scheduleUpdatePayload) Validate(r *http.Request) error {
	if payload.Name != nil && !govalidator.Matches(*payload.Name, `^[a-zA-Z0-9][a-zA-Z0-9_.-]+$`) {
		return errors.New("Invalid schedule name format. Allowed characters are: [a-zA-Z0-9_.-]")
	}
	return nil
}

func (handler *Handler) scheduleUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	settings, err := handler.SettingsService.Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Unable to retrieve settings", err}
	}
	if !settings.EnableHostManagementFeatures {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Host management features are disabled", portainer.ErrHostManagementFeaturesDisabled}
	}

	scheduleID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid schedule identifier route variable", err}
	}

	var payload scheduleUpdatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	schedule, err := handler.ScheduleService.Schedule(portainer.ScheduleID(scheduleID))
	if err == portainer.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a schedule with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a schedule with the specified identifier inside the database", err}
	}

	updateJobSchedule := updateSchedule(schedule, &payload)

	if payload.FileContent != nil {
		_, err := handler.FileService.StoreScheduledJobFileFromBytes(strconv.Itoa(scheduleID), []byte(*payload.FileContent))
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist script file changes on the filesystem", err}
		}
		updateJobSchedule = true
	}

	if updateJobSchedule {
		jobContext := cron.NewScriptExecutionJobContext(handler.JobService, handler.EndpointService, handler.FileService)
		jobRunner := cron.NewScriptExecutionJobRunner(schedule, jobContext)
		err := handler.JobScheduler.UpdateJobSchedule(jobRunner)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update job scheduler", err}
		}
	}

	err = handler.ScheduleService.UpdateSchedule(portainer.ScheduleID(scheduleID), schedule)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist schedule changes inside the database", err}
	}

	return response.JSON(w, schedule)
}

func updateSchedule(schedule *portainer.Schedule, payload *scheduleUpdatePayload) bool {
	updateJobSchedule := false

	if payload.Name != nil {
		schedule.Name = *payload.Name
	}

	if payload.Endpoints != nil {
		schedule.ScriptExecutionJob.Endpoints = payload.Endpoints
		updateJobSchedule = true
	}

	if payload.CronExpression != nil {
		schedule.CronExpression = *payload.CronExpression
		updateJobSchedule = true
	}

	if payload.Recurring != nil {
		schedule.Recurring = *payload.Recurring
		updateJobSchedule = true
	}

	if payload.Image != nil {
		schedule.ScriptExecutionJob.Image = *payload.Image
		updateJobSchedule = true
	}

	if payload.RetryCount != nil {
		schedule.ScriptExecutionJob.RetryCount = *payload.RetryCount
		updateJobSchedule = true
	}

	if payload.RetryInterval != nil {
		schedule.ScriptExecutionJob.RetryInterval = *payload.RetryInterval
		updateJobSchedule = true
	}

	return updateJobSchedule
}

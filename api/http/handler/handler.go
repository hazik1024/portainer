package handler

import (
	"net/http"
	"strings"

	"github.com/hazik1024/portainer/api/http/handler/schedules"

	"github.com/hazik1024/portainer/api/http/handler/roles"

	"github.com/hazik1024/portainer/api/http/handler/auth"
	"github.com/hazik1024/portainer/api/http/handler/dockerhub"
	"github.com/hazik1024/portainer/api/http/handler/endpointgroups"
	"github.com/hazik1024/portainer/api/http/handler/endpointproxy"
	"github.com/hazik1024/portainer/api/http/handler/endpoints"
	"github.com/hazik1024/portainer/api/http/handler/extensions"
	"github.com/hazik1024/portainer/api/http/handler/file"
	"github.com/hazik1024/portainer/api/http/handler/motd"
	"github.com/hazik1024/portainer/api/http/handler/registries"
	"github.com/hazik1024/portainer/api/http/handler/resourcecontrols"
	"github.com/hazik1024/portainer/api/http/handler/settings"
	"github.com/hazik1024/portainer/api/http/handler/stacks"
	"github.com/hazik1024/portainer/api/http/handler/status"
	"github.com/hazik1024/portainer/api/http/handler/tags"
	"github.com/hazik1024/portainer/api/http/handler/teammemberships"
	"github.com/hazik1024/portainer/api/http/handler/teams"
	"github.com/hazik1024/portainer/api/http/handler/templates"
	"github.com/hazik1024/portainer/api/http/handler/upload"
	"github.com/hazik1024/portainer/api/http/handler/users"
	"github.com/hazik1024/portainer/api/http/handler/webhooks"
	"github.com/hazik1024/portainer/api/http/handler/websocket"

	"github.com/hazik1024/portainer/api/custom/build"
	"github.com/hazik1024/portainer/api/custom/stackbackup"
)

// Handler is a collection of all the service handlers.
type Handler struct {
	AuthHandler            *auth.Handler
	DockerHubHandler       *dockerhub.Handler
	EndpointGroupHandler   *endpointgroups.Handler
	EndpointHandler        *endpoints.Handler
	EndpointProxyHandler   *endpointproxy.Handler
	FileHandler            *file.Handler
	MOTDHandler            *motd.Handler
	ExtensionHandler       *extensions.Handler
	RegistryHandler        *registries.Handler
	ResourceControlHandler *resourcecontrols.Handler
	RoleHandler            *roles.Handler
	SchedulesHanlder       *schedules.Handler
	SettingsHandler        *settings.Handler
	StackHandler           *stacks.Handler
	StatusHandler          *status.Handler
	TagHandler             *tags.Handler
	TeamMembershipHandler  *teammemberships.Handler
	TeamHandler            *teams.Handler
	TemplatesHandler       *templates.Handler
	UploadHandler          *upload.Handler
	UserHandler            *users.Handler
	WebSocketHandler       *websocket.Handler
	WebhookHandler         *webhooks.Handler
	BuildHandler           *build.Handler
	StackBackupHandler     *stackbackup.Handler
}

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/api/auth"):
		http.StripPrefix("/api", h.AuthHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/dockerhub"):
		http.StripPrefix("/api", h.DockerHubHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/endpoint_groups"):
		http.StripPrefix("/api", h.EndpointGroupHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/endpoints"):
		switch {
		case strings.Contains(r.URL.Path, "/docker/"):
			http.StripPrefix("/api/endpoints", h.EndpointProxyHandler).ServeHTTP(w, r)
		case strings.Contains(r.URL.Path, "/storidge/"):
			http.StripPrefix("/api/endpoints", h.EndpointProxyHandler).ServeHTTP(w, r)
		case strings.Contains(r.URL.Path, "/azure/"):
			http.StripPrefix("/api/endpoints", h.EndpointProxyHandler).ServeHTTP(w, r)
		default:
			http.StripPrefix("/api", h.EndpointHandler).ServeHTTP(w, r)
		}
	case strings.HasPrefix(r.URL.Path, "/api/extensions"):
		http.StripPrefix("/api", h.ExtensionHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/motd"):
		http.StripPrefix("/api", h.MOTDHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/registries"):
		http.StripPrefix("/api", h.RegistryHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/resource_controls"):
		http.StripPrefix("/api", h.ResourceControlHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/roles"):
		http.StripPrefix("/api", h.RoleHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/schedules"):
		http.StripPrefix("/api", h.SchedulesHanlder).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/settings"):
		http.StripPrefix("/api", h.SettingsHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/stacks"):
		http.StripPrefix("/api", h.StackHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/status"):
		http.StripPrefix("/api", h.StatusHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/tags"):
		http.StripPrefix("/api", h.TagHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/templates"):
		http.StripPrefix("/api", h.TemplatesHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/upload"):
		http.StripPrefix("/api", h.UploadHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/users"):
		http.StripPrefix("/api", h.UserHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/teams"):
		http.StripPrefix("/api", h.TeamHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/team_memberships"):
		http.StripPrefix("/api", h.TeamMembershipHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/websocket"):
		http.StripPrefix("/api", h.WebSocketHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/webhooks"):
		http.StripPrefix("/api", h.WebhookHandler).ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/custom"):
		switch {
		case strings.Contains(r.URL.Path, "/build"):
			http.StripPrefix("/api/custom", h.BuildHandler).ServeHTTP(w, r)
		case strings.Contains(r.URL.Path, "/stack"):
			http.StripPrefix("/api/custom", h.StackBackupHandler).ServeHTTP(w, r)
		}
	case strings.HasPrefix(r.URL.Path, "/"):
		h.FileHandler.ServeHTTP(w, r)
	}
}

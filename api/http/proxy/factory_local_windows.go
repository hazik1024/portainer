// +build windows

package proxy

import (
	"net"
	"net/http"
	"github.com/Microsoft/go-winio"
	
	portainer "github.com/hazik1024/portainer/api"
)

func (factory *proxyFactory) newLocalProxy(path string, endpointID portainer.EndpointID) http.Handler {
	proxy := &localProxy{}
	transport := &proxyTransport{
		enableSignature:        false,
		ResourceControlService: factory.ResourceControlService,
		TeamMembershipService:  factory.TeamMembershipService,
		SettingsService:        factory.SettingsService,
		RegistryService:        factory.RegistryService,
		DockerHubService:       factory.DockerHubService,
		dockerTransport:        newNamedPipeTransport(path),
		endpointIdentifier:     endpointID,
	}
	proxy.Transport = transport
	return proxy
}

func newNamedPipeTransport(namedPipePath string) *http.Transport {
	return &http.Transport{
		Dial: func(proto, addr string) (conn net.Conn, err error) {
			return winio.DialPipe(namedPipePath, nil)
		},
	}
}

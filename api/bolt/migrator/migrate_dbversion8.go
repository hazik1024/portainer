package migrator

import "github.com/hazik1024/portainer/api"

func (m *Migrator) updateEndpointsToVersion9() error {
	legacyEndpoints, err := m.endpointService.Endpoints()
	if err != nil {
		return err
	}

	for _, endpoint := range legacyEndpoints {
		endpoint.GroupID = portainer.EndpointGroupID(1)
		err = m.endpointService.UpdateEndpoint(endpoint.ID, &endpoint)
		if err != nil {
			return err
		}
	}

	return nil
}

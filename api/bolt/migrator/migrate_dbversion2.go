package migrator

import "github.com/hazik1024/portainer/api"

func (m *Migrator) updateSettingsToDBVersion3() error {
	legacySettings, err := m.settingsService.Settings()
	if err != nil {
		return err
	}

	legacySettings.AuthenticationMethod = portainer.AuthenticationInternal
	legacySettings.LDAPSettings = portainer.LDAPSettings{
		TLSConfig: portainer.TLSConfiguration{},
		SearchSettings: []portainer.LDAPSearchSettings{
			portainer.LDAPSearchSettings{},
		},
	}

	return m.settingsService.UpdateSettings(legacySettings)
}

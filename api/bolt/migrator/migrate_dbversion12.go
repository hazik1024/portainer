package migrator

import "github.com/hazik1024/portainer/api"

func (m *Migrator) updateSettingsToVersion13() error {
	legacySettings, err := m.settingsService.Settings()
	if err != nil {
		return err
	}

	legacySettings.LDAPSettings.AutoCreateUsers = false
	legacySettings.LDAPSettings.GroupSearchSettings = []portainer.LDAPGroupSearchSettings{
		portainer.LDAPGroupSearchSettings{},
	}

	return m.settingsService.UpdateSettings(legacySettings)
}

package stack

import "github.com/hazik1024/portainer/api/http/security"

// BackupHandler Stack配置文件备份
type BackupHandler struct {
	backupService *BackupService
}

// NewBackupHandler 返回新的BackupHandler
func NewBackupHandler(bouncer *security.RequestBouncer) *BackupHandler {
	return &BackupHandler{
		backupService: &BackupService{},
	}
}

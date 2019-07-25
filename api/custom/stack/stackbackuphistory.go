package stack

import "github.com/portainer/portainer/api/http/security"

// BackupHistoryHandler 编译镜像
type BackupHistoryHandler struct {
	backupService *BackupService
}

// NewBackupHistoryHandler 返回新的Handler
func NewBackupHistoryHandler(bouncer *security.RequestBouncer) *BackupHistoryHandler {
	return &BackupHistoryHandler{
		backupService: &BackupService{},
	}
}

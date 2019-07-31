package stackbackup

// BackupService 定义BuildService
type BackupService struct {
	// conn mysql.conn
}

// NewBackupService 返回BackupService
func NewBackupService() *BackupService {
	return &BackupService{}
}

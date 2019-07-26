package build

import (
	"github.com/hazik1024/portainer/api/custom/mysqldb"
)

// Service 定义BuildService
type Service struct {
	db mysqldb.MySQLDb
}

// NewService 返回BuildService指针
func NewService() *Service {
	newDb := mysqldb.NewMySQLDb()
	return &Service{
		db: newDb,
	}
}

func (s *Service) queryWarehouse() {

}

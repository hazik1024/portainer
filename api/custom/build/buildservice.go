package build

import (
	"github.com/hazik1024/portainer/api/custom/mysqldb"
	"github.com/hazik1024/portainer/api"
)

// Service 定义BuildService
type Service struct {
	db mysqldb.MySQLDb
}

// NewService 返回BuildService指针
func NewService() *Service {
	return &Service{
		db : mysqldb.NewMySQLDb{}
	}
}

func (s *Service)queryWarehouse() {
	
}
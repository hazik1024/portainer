package build

import (
	"github.com/portainer/portainer/api/custom/mysqldb"
	"github.com/portainer/portainer/api"
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
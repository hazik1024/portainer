package build

import (
	"log"

	"github.com/hazik1024/portainer/api/bolt/registry"
	"github.com/hazik1024/portainer/api/custom/mysqldb"
)

// Service 定义BuildService
type Service struct {
	db              *mysqldb.MySQLDb
	registryService *registry.Service
}

// NewService 返回BuildService指针
func NewService(db *mysqldb.MySQLDb, registryService *registry.Service) *Service {
	return &Service{
		db:              db,
		registryService: registryService,
	}
}

// BuildAndPushImage 编译并推送
func (s *Service) buildAndPushImage(req reqPayload) {
	log.Println(req.GitPath, req.GitBranch, req.ImageName, req.RegistryName, req.RegistryPath)
	s.copySourceCode()
	s.buildSourceCodeToImage()
	s.pushImage()
}

// cloneSourceCode 拷贝源码到本地
func (s *Service) copySourceCode() {

}

// buildImage 编译镜像
func (s *Service) buildSourceCodeToImage() {

}

// pushImage 推送镜像
func (s *Service) pushImage() {

}

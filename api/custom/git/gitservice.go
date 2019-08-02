package git

import (
	"time"

	"github.com/hazik1024/portainer/api/custom/mysqldb"
)

// Service 定义BuildService
type Service struct {
	db *mysqldb.MySQLDb
}

// NewService 返回BuildService指针
func NewService(db *mysqldb.MySQLDb) *Service {
	return &Service{
		db: db,
	}
}

func (service *Service) createGit(payload Payload) (int, error) {
	sql := "insert into t_git (gitname, username, password, address, project, branch, type, addtime, lasttime) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	gitID, err := service.db.Insert(sql, payload.GitName, payload.UserName, payload.Password, payload.Address, payload.Project, payload.Branch, payload.Type, time.Now(), time.Now())
	return gitID, err
}

func (service *Service) deleteGit(id int) (int, error) {
	sql := "delete from t_git where id = ?"
	affects, err := service.db.Delete(sql, id)
	return affects, err
}

func (service *Service) updateGit(id int, payload Payload) (int, error) {
	sql := "update t_git set where id = ?"
	affects, err := service.db.Update(sql, id)
	return affects, err
}

func (service *Service) getGitByID(id int) (*mysqldb.TableGit, error) {
	sql := "select gitid, gitname, username, password, address, project, branch, type, addtime, lasttime from t_git where gitid = ?"
	row := service.db.QueryOne(sql, id)
	var tableGit *mysqldb.TableGit
	err := row.Scan(&tableGit.GitID, &tableGit.GitName, &tableGit.UserName, &tableGit.Password, &tableGit.Address, &tableGit.Project, &tableGit.Branch, &tableGit.Type, &tableGit.AddTime, &tableGit.LastTime)
	if err != nil {
		return nil, err
	}
	return tableGit, nil
}

func (service *Service) getGitByAddress(address string) (*mysqldb.TableGit, error) {
	sql := "select gitid, gitname, username, password, address, project, branch, type, addtime, lasttime from t_git where address = ?"
	row := service.db.QueryOne(sql, address)
	var tableGit *mysqldb.TableGit
	err := row.Scan(&tableGit.GitID, &tableGit.GitName, &tableGit.UserName, &tableGit.Password, &tableGit.Address, &tableGit.Project, &tableGit.Branch, &tableGit.Type, &tableGit.AddTime, &tableGit.LastTime)
	if err != nil {
		return nil, err
	}
	return tableGit, nil
}

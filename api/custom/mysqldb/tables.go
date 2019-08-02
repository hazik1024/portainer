package mysqldb

import "time"

// TableRegistry Git地址表
type TableRegistry struct {
	RegistryID int       `db:"registryid"`
	Name       string    `db:"name"`
	Address    string    `db:"address"`
	Project    string    `db:"project"`
	UserName   string    `db:"username"`
	Password   string    `db:"password"`
	Status     int       `db:"status"`
	Type       int       `db:"type"`
	AddTime    time.Time `db:"addtime"`
	LastTime   time.Time `db:"lasttime"`
}

// TableName 表名
func (info *TableRegistry) TableName() string {
	return "t_registry"
}

// TableGit Git地址表
type TableGit struct {
	GitID    int       `db:"gitid"`
	GitName  string    `db:"gitname"`
	UserName string    `db:"username"`
	Password string    `db:"password"`
	Address  string    `db:"address"`
	Project  string    `db:"project"`
	Branch   string    `db:"branch"`
	Type     int       `db:"type"`
	AddTime  time.Time `db:"addtime"`
	LastTime time.Time `db:"lasttime"`
}

// TableName 表名
func (info *TableGit) TableName() string {
	return "t_git"
}

// TableBuild Git地址表
type TableBuild struct {
	BuildID    int       `db:"buildid"`
	GitID      int       `db:"gitid"`
	RegistryID int       `db:"registryid"`
	ImageName  string    `db:"imagename"`
	Branch     string    `db:"branch"`
	Tag        string    `db:"tag"`
	Status     int       `db:"status"`
	AdminName  string    `db:"adminname"`
	AddTime    time.Time `db:"addtime"`
	LastTime   time.Time `db:"lasttime"`
}

// TableName 表名
func (info *TableBuild) TableName() string {
	return "t_build"
}

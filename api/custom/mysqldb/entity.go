package mysqldb

// WareHouse 镜像地址表
type WareHouse struct {
	id       uint   `db:"id"`
	name     string `db:"name"`
	username string `db:"username"`
	passwd   string `db:"passwd"`
	address  string `db:"address"`
	project  string `db:"project"`
	typ      string `db:"type"`
}

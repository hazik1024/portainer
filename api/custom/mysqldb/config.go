package mysqldb

import "time"

const (
	dbtype      = "mysql"
	maxopenconn = 20
	maxidleconn = 5
	maxlifetime = 60 * time.Second
)

const (
	username  = "devrelease"
	password  = "devrelease"
	network   = "tcp"
	ip        = "10.0.0.143"
	port      = 3306
	addr      = "10.0.0.143:3306"
	dbname    = "devrelease"
	charset   = "utf8"
	parsetime = "true"
)

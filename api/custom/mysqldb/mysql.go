package mysqldb

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

var (
	db   *sql.DB
	lock sync.Mutex
)

// MySQLDb 操作类
type MySQLDb struct {
}

// NewMySQLDb MySQLDb实例化函数
func NewMySQLDb() *MySQLDb {
	initDb()
	return &MySQLDb{}
}

func initDb() {
	if db == nil {
		lock.Lock()
		if db == nil {
			dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=%s&&parseTime=%s", username, password, network, ip, port, dbname, charset, parsetime)
			var err error
			db, err = sql.Open(dbtype, dsn)
			if err != nil {
				log.Println(err)
				lock.Unlock()
				return
			}
			db.SetConnMaxLifetime(maxlifetime)
			db.SetMaxOpenConns(maxopenconn)
			db.SetMaxIdleConns(maxidleconn)
		}
		lock.Unlock()
	}
}

// Close 关闭db
func (mydb *MySQLDb) Close() {
	if db != nil {
		db.Close()
	}
}

// Insert 插入
func (mydb *MySQLDb) Insert(query string) int64 {
	result, err := mydb.Execute(query)
	if err != nil {
		log.Println("insert fail", err)
		return -1
	}
	lastInsertID, err1 := result.LastInsertId()
	if err1 != nil {
		log.Println("get lastInsertID fail", err)
		return -1
	}
	return lastInsertID
}

// Delete 删除
func (mydb *MySQLDb) Delete(query string) int64 {
	_, err := mydb.Execute(query)
	if err != nil {
		log.Println("delete fail", err)
		return -1
	}
	return 0
}

// Update 更新
func (mydb *MySQLDb) Update(query string) int64 {
	_, err := mydb.Execute(query)
	if err != nil {
		log.Println("update fail", err)
		return -1
	}
	return 0
}

// QueryOne 单挑查询
func (*MySQLDb) QueryOne(query string) *sql.Row {
	row := db.QueryRow(query)
	return row
}

// Query 查询
func (*MySQLDb) Query(query string) (*sql.Rows, error) {
	return db.Query(query)
}

// Execute 执行
func (*MySQLDb) Execute(query string) (sql.Result, error) {
	return db.Exec(query)
}

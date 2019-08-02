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
func (mydb *MySQLDb) Insert(query string, args ...interface{}) (int, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(args)
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(lastInsertID), nil
}

// Delete 删除
func (mydb *MySQLDb) Delete(query string, args ...interface{}) (int, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(args)
	if err != nil {
		return -1, err
	}
	affects, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return int(affects), nil
}

// Update 更新
func (mydb *MySQLDb) Update(query string, args ...interface{}) (int, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(args)
	if err != nil {
		return -1, err
	}
	affects, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return int(affects), nil
}

// QueryOne 单挑查询
func (*MySQLDb) QueryOne(query string, args ...interface{}) *sql.Row {
	row := db.QueryRow(query, args)
	return row
}

// Query 查询
func (*MySQLDb) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args)
}

package utils

import (
	"fmt"
	mySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type Mysql struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	Timeout  string

	lock sync.Mutex // 单例模式锁
}

func DefaultMySQLDB() *Mysql {
	return &Mysql{
		Username: "douyin",
		Password: "douyin",
		Host:     "43.143.14.234",
		Port:     "3306",
		DBName:   "mini_douyin",
		Timeout:  "10s",
	}
}

func (m *Mysql) GetDB(db *gorm.DB) (*gorm.DB, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if db == nil {
		conn, err := m.connectDefaultDB()
		if err != nil {
			return nil, fmt.Errorf("GetDB失败, %s", err)
		}
		db = conn
	}
	return db, nil
}

func (m *Mysql) connectDefaultDB() (*gorm.DB, error) {
	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&loc=Local&timeout=%s", m.Username, m.Password, m.Host, m.Port, m.DBName, m.Timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mySQL.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败, %s", err)
	}
	// 连接成功
	fmt.Println("数据库连接成功")
	return db, nil
}

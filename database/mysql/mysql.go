package mysql

import (
	"errors"
	"fmt"
	"github.com/fundata-varena/fundata-resource-server/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var shared *xorm.Engine

func Init() error {
	config, err := conf.GetConf()
	if err != nil {
		return errors.New("conf is nil")
	}

	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		config.Mysql.User,
		config.Mysql.Password,
		config.Mysql.Host,
		config.Mysql.Port,
		config.Mysql.Db)
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		return errors.New("db connect failed " + err.Error())
	}

	err = engine.Ping()
	if err != nil {
		return errors.New("ping db err " + err.Error())
	}

	shared = engine

	return nil
}

func GetInstance() (*xorm.Engine, error) {
	if shared == nil {
		return nil, errors.New("init mysql first please")
	}
	return shared, nil
}

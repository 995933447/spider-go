package saverclient

import (
	"github.com/go-xorm/xorm"
	"spider-go/logger/mysql"
	mysqlSaver "spider-go/persist/mysql"
)

func GetMysqlClient() *xorm.Engine {
	client, err := mysqlSaver.NewClient()
	if err != nil {
		panic(err)
	}

	if err = mysql.ListenSql(client); err != nil {
		panic(err)
	}

	return client
}
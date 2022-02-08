package db_op

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlDb *gorm.DB

func SqlInit() {
	db, err := gorm.Open("mysql", "root:qingjiuling1491!@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		defer db.Close()
		return
	}
	MysqlDb = db
}

func SqlClose() {
	MysqlDb.Close()
}

package main

import (
	"awesomeProject/db_op"
	"awesomeProject/types"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	//test123
	db_op.SqlInit()
	defer db_op.SqlClose()
	db_op.RedisInit()
	defer db_op.RedisClose()
	defer fmt.Println(1123)
	conn := db_op.MysqlDb
	// 不要复数表名
	conn.SingularTable(true)

	// 借助 gorm 创建数据库表.
	conn.AutoMigrate(&types.MemberSql{})

	if !conn.HasTable(types.MemberSql{}) {
		conn.AutoMigrate(types.MemberSql{})
		if conn.HasTable(types.MemberSql{}) {
			fmt.Println("balance表创建成功")
		} else {
			fmt.Println("balance表创建失败")
		}
	} else {
		fmt.Println("表已存在")
	}
	var member = types.MemberSql{Nickname: "JudgeAdmin", Username: "JudgeAdmin", UserType: 1, PassWord: "JudgePassword2022"}
	db_op.MysqlDb.Create(&member)

	// 借助 gorm 创建数据库表.
	conn.AutoMigrate(&types.CourseSql{})

	if !conn.HasTable(types.CourseSql{}) {
		conn.AutoMigrate(types.CourseSql{})
		if conn.HasTable(types.CourseSql{}) {
			fmt.Println("balance表创建成功")
		} else {
			fmt.Println("balance表创建失败")
		}
	} else {
		fmt.Println("表已存在")
	}
	var course = types.CourseSql{CourseName: "math", Cap: 10}
	db_op.MysqlDb.Create(&course)

	// 1.创建路由

	r := gin.Default()
	types.RegisterRouter(r)
	r.Run(":8080")

}

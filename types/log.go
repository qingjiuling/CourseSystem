package types

import (
	"awesomeProject/db_op"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

//登录

func Login(c *gin.Context) {

	var request LoginRequest
	var response LoginResponse

	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		c.JSON(http.StatusOK, response)
		return
	}
	username := request.Username
	password := request.Password
	conn := db_op.MysqlDb
	var memberSql MemberSql

	// 根据username查表，若无匹配则return
	if err := conn.Where("username = ?", username).First(&memberSql).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		response.Code = UserNotExisted
		c.JSON(http.StatusOK, response)
		return
	}

	// 查看是否已删除
	if memberSql.Deleted == false {
		response.Code = UserHasDeleted
		c.JSON(http.StatusOK, response)
		return
	}

	// 查看密码是否正确
	if memberSql.PassWord != password {
		response.Code = WrongPassword
		c.JSON(http.StatusOK, response)
		return
	}

	// 设置cookie
	c.SetCookie("camp-session", username, 10000, "/", "127.0.0.1", false, true)
	response.Code = OK
	response.Data.UserID = strconv.FormatInt(memberSql.UserID, 10)
	c.JSON(http.StatusOK, response)

}

func Logout(c *gin.Context) {
	var response LogoutResponse
	cookie, err := c.Cookie("camp-session")
	//根据cookie判断是否已登录
	if err != nil {
		response.Code = LoginRequired
		c.JSON(http.StatusOK, response)
		return
	}
	//清除cookie
	c.SetCookie("camp-session", cookie, -1, "/", "127.0.0.1", false, true)

	response.Code = OK
	c.JSON(http.StatusOK, response)
}

func Whoami(c *gin.Context) {
	var response WhoAmIResponse
	cookie, err := c.Cookie("camp-session")
	//根据cookie判断是否已登录
	if err != nil {
		response.Code = LoginRequired
		c.JSON(http.StatusOK, response)
		return
	}
	//获取个人信息
	var memberSql MemberSql
	username := cookie
	conn := db_op.MysqlDb
	conn.Where("username = ?", username).First(&memberSql)
	response.Code = OK
	response.Data.UserID = strconv.FormatInt(memberSql.UserID, 10)
	response.Data.UserType = memberSql.UserType
	response.Data.Username = memberSql.Username
	response.Data.Nickname = memberSql.Nickname
	c.JSON(http.StatusOK, response)
}

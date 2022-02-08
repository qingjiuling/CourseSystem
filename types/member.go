package types

import (
	"awesomeProject/db_op"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateMember(c *gin.Context) {
	var request CreateMemberRequest
	var response CreateMemberResponse
	cookie, err := c.Cookie("camp-session")
	//根据cookie判断是否已登录
	if err != nil {
		response.Code = LoginRequired
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//获取操作者个人信息
	var memberSql MemberSql
	usernameRe := cookie
	conn := db_op.MysqlDb
	conn.Where("username = ?", usernameRe).First(&memberSql)
	if memberSql.UserType != Admin {
		response.Code = PermDenied
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//获取昵称、用户名、密码、用户类型等信息
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	nickname := request.Nickname
	username := request.Username
	password := request.Password
	usertype := request.UserType
	if len(nickname) < 4 || len(nickname) > 20 {
		response.Code = ParamInvalid
		//fmt.Println(2)
		c.JSON(http.StatusBadRequest, response)
	}
	if len(username) < 8 || len(username) > 20 {
		response.Code = ParamInvalid
		//fmt.Println(3)
		c.JSON(http.StatusBadRequest, response)
	} else {
		for _, ch := range username {
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
				continue
			} else {
				//fmt.Println(4, ch)
				response.Code = ParamInvalid
				c.JSON(http.StatusBadRequest, response)
				return
			}
		}
	}
	var members []MemberSql
	conn.Where("username = ?", username).Find(&members)
	if len(members) == 1 {
		response.Code = UserHasExisted
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if len(password) < 8 || len(password) > 20 {
		//fmt.Println(5)
		response.Code = ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	} else {
		flag1 := false
		flag2 := false
		flag3 := false
		for _, ch := range password {
			if ch >= 'a' && ch <= 'z' {
				flag1 = true
			}
			if ch >= 'A' && ch <= 'Z' {
				flag2 = true
			}
			if ch >= '0' && ch <= '9' {
				flag3 = true
			}
		}
		if !(flag1 && flag2 && flag3) {
			response.Code = ParamInvalid
			c.JSON(http.StatusBadRequest, response)
			//fmt.Println(6, flag1, flag2, flag3)
			return
		}

	}
	if usertype != Admin && usertype != Student && usertype != Teacher {
		//fmt.Println(7)
		response.Code = ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var member = MemberSql{Nickname: nickname, Username: username, UserType: usertype, PassWord: password}
	db_op.MysqlDb.Create(&member)
	response.Code = OK
	conn.Where("username = ?", username).First(&member)
	response.Data.UserID = strconv.FormatInt(member.UserID, 10)
	c.JSON(http.StatusOK, response)
}

func GetMember(c *gin.Context) {
	var request GetMemberRequest
	var response GetMemberResponse
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	userid, _ := strconv.ParseInt(request.UserID, 10, 64)
	conn := db_op.MysqlDb
	var members []MemberSql
	conn.Where("user_id = ?", userid).Find(&members)
	if len(members) == 0 {
		response.Code = UserNotExisted
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if members[0].Deleted == false {
		response.Code = UserHasDeleted
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response.Code = OK
	response.Data.UserID = strconv.FormatInt(userid, 10)
	response.Data.Nickname = members[0].Nickname
	response.Data.Username = members[0].Username
	response.Data.UserType = members[0].UserType
	c.JSON(http.StatusOK, response)
}

func GetMemberList(c *gin.Context) {
	var request GetMemberListRequest
	var response GetMemberListResponse
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	limit := request.Limit
	offset := request.Offset
	conn := db_op.MysqlDb
	var members []MemberSql
	conn.Offset(offset).Limit(limit).Find(&members)
	tmembers := make([]TMember, len(members))
	for i, member := range members {
		tmembers[i].UserType = member.UserType
		tmembers[i].Nickname = member.Nickname
		tmembers[i].UserID = strconv.FormatInt(member.UserID, 10)
		tmembers[i].Username = member.Username
	}
	response.Code = OK
	response.Data.MemberList = tmembers
	c.JSON(http.StatusOK, response)
}

func UpdateMember(c *gin.Context) {
	var request UpdateMemberRequest
	var response UpdateMemberResponse
	//json解析
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//userid查询
	nickname := request.Nickname
	userid, _ := strconv.ParseInt(request.UserID, 10, 64)
	conn := db_op.MysqlDb
	var member MemberSql
	if conn.Where("user_id = ?", userid).First(&member).RecordNotFound() {
		response.Code = UserNotExisted
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//判断new nickname是否正确
	if len(nickname) < 4 || len(nickname) > 20 {
		response.Code = ParamInvalid
		//fmt.Println(2)
		c.JSON(http.StatusBadRequest, response)
	}
	//更新nickname
	conn.Model(&member).UpdateColumn("nickname", nickname)
	response.Code = OK
	c.JSON(http.StatusOK, response)
}

func DeleteMember(c *gin.Context) {
	var request DeleteMemberRequest
	var response DeleteMemberResponse
	//json解析
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//userid查询
	userid, _ := strconv.ParseInt(request.UserID, 10, 64)
	conn := db_op.MysqlDb
	var member MemberSql
	if conn.Where("user_id = ?", userid).First(&member).RecordNotFound() {
		response.Code = UserNotExisted
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if member.Deleted == false {
		response.Code = UserHasDeleted
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//更新nickname
	conn.Model(&member).UpdateColumn("deleted", false)
	response.Code = OK
	c.JSON(http.StatusOK, response)
}

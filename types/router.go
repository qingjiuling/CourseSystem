package types

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create", CreateMember)
	g.GET("/member", GetMember)
	g.GET("/member/list", GetMemberList)
	g.POST("/member/update", UpdateMember)
	g.POST("/member/delete", DeleteMember)

	// 登录
	g.POST("/auth/login", Login)
	g.POST("/auth/logout", Logout)
	g.GET("/auth/whoami", Whoami)

	// 排课
	g.POST("/course/create", CreateCourse)
	g.GET("/course/get", GetCourse)

	g.POST("/teacher/bind_course", BindCourse)
	g.POST("/teacher/unbind_course", UnbindCourse)
	g.GET("/teacher/get_course", GetTeacherCourse)
	g.POST("/course/schedule", ScheduleCourse)

	// 抢课
	g.POST("/student/book_course", BookCourse)
	g.GET("/student/course", GetStudentCourse)

}

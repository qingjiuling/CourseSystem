package types

import (
	"awesomeProject/db_op"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateCourse(c *gin.Context) {
	var request CreateCourseRequest
	var response CreateCourseResponse
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusOK, response)
		return
	}
	courseCap := request.Cap
	courseName := request.Name
	conn := db_op.MysqlDb
	var course = CourseSql{CourseName: courseName, Cap: courseCap}
	conn.Create(&course)
	conn.Where("courseName'?", courseName).First(&course)
	response.Code = OK
	response.Data.CourseID = strconv.FormatInt(course.CourseID, 10)
	c.JSON(http.StatusOK, response)

}

func GetCourse(c *gin.Context) {
	var request GetCourseRequest
	var response GetCourseResponse
	if err := c.ShouldBindUri(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusOK, response)
		return
	}
	courseID, err := strconv.ParseInt(request.CourseID, 10, 64)
	if err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusOK, response)
		return
	}
	var course CourseSql
	conn := db_op.MysqlDb
	if conn.Where("course_id=?", courseID).First(&course).RecordNotFound() {
		response.Code = CourseNotExisted
		c.JSON(http.StatusOK, response)
		return
	}
	response.Code = OK
	response.Data.Name = course.CourseName
	response.Data.CourseID = strconv.FormatInt(course.CourseID, 10)
	response.Data.TeacherID = course.TeacherID
	c.JSON(http.StatusOK, response)
}

func BindCourse(c *gin.Context) {
	var request BindCourseRequest
	var response BindCourseResponse
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusOK, response)
		return
	}
	teacherID := request.TeacherID
	courseID, _ := strconv.ParseInt(request.CourseID, 10, 64)
	conn := db_op.MysqlDb
	var course CourseSql
	m := conn.Where("course_id=?", courseID).First(&course)
	if m.RecordNotFound() {
		response.Code = CourseNotExisted
		c.JSON(http.StatusOK, response)
		return
	}
	if course.TeacherID != "" {
		response.Code = CourseHasBound
		c.JSON(http.StatusOK, response)
		return
	}
	m.Update("teacher_id", teacherID)
	response.Code = OK
	c.JSON(http.StatusOK, response)
}

func UnbindCourse(c *gin.Context) {
	var request UnbindCourseRequest
	var response UnbindCourseResponse
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusOK, response)
		return
	}
	teacherID := request.TeacherID
	courseID, _ := strconv.ParseInt(request.CourseID, 10, 64)
	conn := db_op.MysqlDb
	var course CourseSql
	m := conn.Where("course_id=?", courseID).First(&course)
	if course.TeacherID != teacherID || course.TeacherID == "" {
		response.Code = CourseNotBind
		c.JSON(http.StatusOK, response)
		return
	}
	m.Update("teacher_id", "")
	response.Code = OK
	c.JSON(http.StatusOK, response)
}

func GetTeacherCourse(c *gin.Context) {
	var request GetTeacherCourseRequest
	var response GetTeacherCourseResponse
	if err := c.ShouldBindUri(&request); err != nil {
		response.Code = ParamInvalid
		//fmt.Println(1)
		c.JSON(http.StatusOK, response)
		return
	}
	teacherID := request.TeacherID
	conn := db_op.MysqlDb
	var courses []CourseSql
	conn.Where("teacher_id=?", teacherID).Find(&courses)
	tCourses := make([]TCourse, len(courses))
	pCourses := make([]*TCourse, len(courses))
	for i, course := range courses {
		tCourses[i].CourseID = strconv.FormatInt(course.CourseID, 10)
		tCourses[i].Name = course.CourseName
		tCourses[i].TeacherID = course.TeacherID
		pCourses[i] = &tCourses[i]
	}
	response.Data.CourseList = pCourses
	c.JSON(http.StatusOK, response)
}

func ScheduleCourse(c *gin.Context) {
	/*
		var request ScheduleCourseRequest
		var response ScheduleCourseResponse
		if err := c.ShouldBindJSON(&request); err != nil {
			response.Code = ParamInvalid
			//fmt.Println(1)
			c.JSON(http.StatusOK, response)
			return
		}
		teCoRe := request.TeacherCourseRelationShip //获取老师到课程的关系表
		coTeRe := make(map[string][]string)         //建立课程到老师的关系表
		for teacher := range teCoRe {
			courses := teCoRe[teacher]
			for _, course := range courses {
				key := course
				value, ok := coTeRe[key]
				if !ok {
					value = make([]string, 0, 2)
				}
				value = append(value, teacher)
				coTeRe[key] = value
			}
		}

	*/

}

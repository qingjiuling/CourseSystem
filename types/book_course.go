package types

import (
	"awesomeProject/db_op"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func BookCourse(c *gin.Context) {
	var request BookCourseRequest
	var response BookCourseResponse
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	rdb := db_op.RedisDb.Get()
	defer rdb.Close()
	conn := db_op.MysqlDb
	courseID, err := strconv.ParseInt(request.CourseID, 10, 64)
	if err != nil {
		response.Code = ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	studentID, err := strconv.ParseInt(request.StudentID, 10, 64)
	if err != nil {
		response.Code = ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var course CourseSql
	if conn.Where("course_id=?", courseID).First(&course).RecordNotFound() {
		response.Code = CourseNotExisted
		c.JSON(http.StatusBadRequest, response)
		return
	}
	cap, err := redis.Int(rdb.Do("Get", courseID))
	// 课程容量不够时
	if cap == 0 {
		response.Code = CourseNotAvailable
		c.JSON(http.StatusBadRequest, response)
		return
	}
	cap -= 1
	rdb.Do("Set", courseID, cap)
	var bookcourse = BookCourseSql{StudentID: studentID, CourseID: courseID}
	conn.Create(&bookcourse)
	response.Code = OK
	c.JSON(http.StatusOK, response)
}

func GetStudentCourse(c *gin.Context) {
	var request GetStudentCourseRequest
	var response GetStudentCourseResponse
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Code = ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	conn := db_op.MysqlDb
	studentID := request.StudentID
	var courses []CourseSql
	conn.Where("student_id=?", studentID).Find(&courses)
	tCourses := make([]TCourse, len(courses))
	for i, course := range courses {
		tCourses[i].CourseID = strconv.FormatInt(course.CourseID, 10)
		tCourses[i].Name = course.CourseName
		tCourses[i].TeacherID = course.TeacherID
	}
	response.Data.CourseList = tCourses
	response.Code = OK
	c.JSON(http.StatusOK, response)
}

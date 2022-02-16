package types

import (
	"awesomeProject/db_op"
	"fmt"
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
		c.JSON(http.StatusOK, response)
		return
	}
	// 从redis线程池取一个线程
	rdb := db_op.RedisDb.Get()
	defer rdb.Close()
	// 与mysql建立连接
	conn := db_op.MysqlDb
	//将couurseID转换为整形
	courseID, err := strconv.ParseInt(request.CourseID, 10, 64)
	if err != nil {
		response.Code = ParamInvalid
		c.JSON(http.StatusOK, response)
		return
	}
	//将studentID转换为整形
	studentID, err := strconv.ParseInt(request.StudentID, 10, 64)
	if err != nil {
		response.Code = ParamInvalid
		c.JSON(http.StatusOK, response)
		return
	}
	//检查学生是否存在
	var student MemberSql
	if conn.Where("user_id=? AND user_type=?", studentID, 2).First(&student).RecordNotFound() {
		response.Code = StudentNotExisted
		c.JSON(http.StatusOK, response)
		return
	}
	//检查课程是否存在
	var course CourseSql
	if conn.Where("course_id=?", courseID).First(&course).RecordNotFound() {
		response.Code = CourseNotExisted
		c.JSON(http.StatusOK, response)
		return
	}
	//检查学生是否已经选上该课程
	tmp1, _ := redis.Int(rdb.Do("SISMEMBER", "student_"+request.StudentID, courseID))
	if tmp1 == 1 {
		response.Code = StudentHasCourse
		c.JSON(http.StatusOK, response)
		return
	}
	//利用DECR来减少库存，防止读写分离
	tmp2, _ := redis.Int(rdb.Do("DECR", "course_"+request.CourseID))
	if tmp2 < 0 {
		response.Code = CourseNotAvailable
		c.JSON(http.StatusOK, response)
		return
	}
	// 将选的课插进学生的集合里
	_, err = rdb.Do("SADD", "student_"+request.StudentID, courseID)
	if err != nil {
		fmt.Print(err)
	}
	response.Code = OK
	c.JSON(http.StatusOK, response)
}

func GetStudentCourse(c *gin.Context) {
	var request GetStudentCourseRequest
	var response GetStudentCourseResponse
	if err := c.ShouldBindQuery(&request); err != nil {
		response.Code = ParamInvalid
		c.JSON(http.StatusOK, response)
		return
	}
	// 从redis线程池取一个线程
	rdb := db_op.RedisDb.Get()
	defer rdb.Close()
	conn := db_op.MysqlDb
	//将student转为整型
	studentID, err := strconv.ParseInt(request.StudentID, 10, 64)
	if err != nil {
		response.Code = ParamInvalid
		c.JSON(http.StatusOK, response)
		return
	}
	//再来判断学生是否存在
	var student MemberSql
	if conn.Where("user_id=? AND user_type=?", studentID, 2).First(&student).RecordNotFound() {
		response.Code = StudentNotExisted
		c.JSON(http.StatusOK, response)
		return
	}
	//先来判断学生是否有课程
	courses, err := redis.Ints(rdb.Do("SMEMBERS", "student_"+request.StudentID))
	courses_len := len(courses)
	fmt.Println(courses_len)
	if courses_len == 0 {
		response.Code = StudentHasNoCourse
		c.JSON(http.StatusOK, response)
		return
	}
	tCourses := make([]TCourse, courses_len)
	for i := 0; i < courses_len; i++ {
		tCourses[i].CourseID = strconv.Itoa(courses[i])
		var course CourseSql
		conn.Where("select course_id=?", tCourses[i].CourseID).First(&course)
		tCourses[i].Name = course.CourseName
		tCourses[i].TeacherID = course.TeacherID
	}
	response.Data.CourseList = tCourses
	response.Code = OK
	c.JSON(http.StatusOK, response)
}

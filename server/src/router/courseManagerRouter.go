package router

import (
	"net/http"
	"orkidslearning/src/controller"
	"orkidslearning/src/models"

	"github.com/gin-gonic/gin"
)

func GetAllCourses(c *gin.Context) {
	var courses []models.Course
	courses = controller.GetAllCourses()
	if courses == nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
	c.JSON(http.StatusOK, courses)
}

func GetCourseById(c *gin.Context) {
	var course *models.Course
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	course = controller.GetCourseById(id)
	if course == nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, course)
}

func AddCourse(c *gin.Context) {
	var course models.AddCourse
	c.BindJSON(&course)
	controller.AddCourse(course)
	c.JSON(http.StatusOK, gin.H{"message": "Course added successfully"})
}

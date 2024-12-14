package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"orkidslearning/src/controller"
	"orkidslearning/src/models"
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
	// Make call to database
	// get course for specific Id
	// display the course
}

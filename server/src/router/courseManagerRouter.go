package router

import (
	"net/http"
	"orkidslearning/src/controller"
	models "orkidslearning/src/models/database"

	"github.com/gin-gonic/gin"
)

func GetAllCourses(c *gin.Context) {
	courses := controller.GetAllCourses()
	if courses == nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
	c.JSON(http.StatusOK, courses)
}

func GetCourseById(c *gin.Context) {
	var enrollInCourse models.EnrollInCourse
	if err := c.ShouldBindJSON(&enrollInCourse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	course := controller.GetCourseById(id)
	if course == nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	isEnrolled, err := controller.IsUserEnrolledInCourse(enrollInCourse.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll in course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enrolled": isEnrolled, "course": course})
}

func AddCourse(c *gin.Context) {
	var course models.AddCourse
	err := c.BindJSON(&course)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	addedCourse, err := controller.AddCourse(course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add course"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Course added successfully", "course": addedCourse})
}

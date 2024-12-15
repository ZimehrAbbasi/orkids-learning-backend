package router

import (
	"log"
	"net/http"
	"orkidslearning/src/controller"
	"orkidslearning/src/models"

	"github.com/gin-gonic/gin"
)

func EnrollInCourse(c *gin.Context) {
	log.Println("Enrolling in course", c.Request)
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
	err := controller.EnrollInCourse(enrollInCourse.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll in course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enrolled in course"})
}

func UnenrollFromCourse(c *gin.Context) {
	var unenrollFromCourse models.UnenrollFromCourse
	if err := c.ShouldBindJSON(&unenrollFromCourse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	err := controller.UnenrollFromCourse(unenrollFromCourse.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unenroll from course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unenrolled from course"})
}

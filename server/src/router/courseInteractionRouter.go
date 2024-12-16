package router

import (
	"log"
	"net/http"
	"orkidslearning/src/controller"
	models "orkidslearning/src/models/database"
	"orkidslearning/src/models/response"

	"github.com/gin-gonic/gin"
)

func EnrollInCourse(c *gin.Context) {
	log.Println("Enrolling in course", c.Request)
	var enrollInCourse models.EnrollInCourse
	if err := c.ShouldBindJSON(&enrollInCourse); err != nil {
		c.JSON(http.StatusBadRequest, response.EnrollInCourseResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.EnrollInCourseResponse{
			Message: "id is required",
			Error:   "id is required",
		})
		return
	}
	err := controller.EnrollInCourse(enrollInCourse.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.EnrollInCourseResponse{
			Message: "Failed to enroll in course",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.EnrollInCourseResponse{
		Message:  "Enrolled in course",
		Enrolled: true,
		Username: enrollInCourse.Username,
		CourseId: id,
	})
}

func UnenrollFromCourse(c *gin.Context) {
	var unenrollFromCourse models.UnenrollFromCourse
	if err := c.ShouldBindJSON(&unenrollFromCourse); err != nil {
		c.JSON(http.StatusBadRequest, response.UnenrollFromCourseResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.UnenrollFromCourseResponse{
			Message: "id is required",
			Error:   "id is required",
		})
		return
	}
	err := controller.UnenrollFromCourse(unenrollFromCourse.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.UnenrollFromCourseResponse{
			Message: "Failed to unenroll from course",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.UnenrollFromCourseResponse{
		Message:    "Unenrolled from course",
		Unenrolled: true,
		Username:   unenrollFromCourse.Username,
		CourseId:   id,
	})
}

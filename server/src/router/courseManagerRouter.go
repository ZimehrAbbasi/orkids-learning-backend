package router

import (
	"context"
	"net/http"
	"orkidslearning/src/controller"
	models "orkidslearning/src/models/database"
	"orkidslearning/src/models/response"
	"orkidslearning/src/services"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllCourses(c *gin.Context) {

	contextService, exists := c.MustGet("contextService").(*services.ContextService)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contextService"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	courses, err := controller.GetAllCourses(ctx, contextService.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.GetCoursesResponse{
			Message: "Failed to get courses",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.GetCoursesResponse{
		Message: "Courses retrieved successfully",
		Courses: courses,
	})
}

func GetCourseById(c *gin.Context) {

	contextService, exists := c.MustGet("contextService").(*services.ContextService)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contextService"})
		return
	}

	var enrollInCourse models.EnrollInCourse
	if err := c.ShouldBindJSON(&enrollInCourse); err != nil {
		c.JSON(http.StatusBadRequest, response.GetCourseResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.GetCourseResponse{
			Message: "id is required",
			Error:   "id is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	course, err := controller.GetCourseById(ctx, contextService.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.GetCourseResponse{
			Message: "Failed to get course",
			Error:   err.Error(),
		})
		return
	}

	var isEnrolled bool = false
	if enrollInCourse.CheckEnrollment {
		isEnrolled, err = controller.IsUserEnrolledInCourse(ctx, contextService.GetDB(), enrollInCourse.Username, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GetCourseResponse{
				Message: "Failed to enroll in course",
				Error:   err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, response.GetCourseResponse{
		Message:  "Course retrieved successfully",
		Course:   *course,
		Enrolled: isEnrolled,
	})
}

func AddCourse(c *gin.Context) {

	contextService, exists := c.MustGet("contextService").(*services.ContextService)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contextService"})
		return
	}

	var course models.AddCourse
	err := c.BindJSON(&course)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.AddCourseResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	addedCourse, err := controller.AddCourse(ctx, contextService.GetDB(), course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.AddCourseResponse{
			Message: "Failed to add course",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.AddCourseResponse{
		Message: "Course added successfully",
		Course:  *addedCourse,
		Added:   true,
	})
}

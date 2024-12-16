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
	"go.opentelemetry.io/otel"
)

func EnrollInCourse(c *gin.Context) {
	tracer := otel.Tracer("router")
	ctx, span := tracer.Start(c.Request.Context(), "EnrollInCourse")
	defer span.End()

	contextService, exists := c.MustGet("contextService").(*services.ContextService)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contextService"})
		return
	}

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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := controller.EnrollInCourse(ctx, contextService.GetDB(), enrollInCourse.Username, id)
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

	tracer := otel.Tracer("router")
	ctx, span := tracer.Start(c.Request.Context(), "UnenrollFromCourse")
	defer span.End()

	contextService, exists := c.MustGet("contextService").(*services.ContextService)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contextService"})
		return
	}

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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := controller.UnenrollFromCourse(ctx, contextService.GetDB(), unenrollFromCourse.Username, id)
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

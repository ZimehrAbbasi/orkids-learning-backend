package router

import (
	"net/http"
	"orkidslearning/src/controller"
	models "orkidslearning/src/models/database"
	"orkidslearning/src/models/response"

	"github.com/gin-gonic/gin"
)

func GetAllCourses(c *gin.Context) {
	courses, err := controller.GetAllCourses()
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
	course, err := controller.GetCourseById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.GetCourseResponse{
			Message: "Failed to get course",
			Error:   err.Error(),
		})
		return
	}

	isEnrolled, err := controller.IsUserEnrolledInCourse(enrollInCourse.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.GetCourseResponse{
			Message: "Failed to enroll in course",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.GetCourseResponse{
		Message:  "Course retrieved successfully",
		Course:   *course,
		Enrolled: isEnrolled,
	})
}

func AddCourse(c *gin.Context) {
	var course models.AddCourse
	err := c.BindJSON(&course)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.AddCourseResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}
	addedCourse, err := controller.AddCourse(course)
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

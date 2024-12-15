package controller

import (
	"log"
	"orkidslearning/src/database"
	"orkidslearning/src/models"
)

func GetAllCourses() []models.Course {
	var courses []models.Course
	courses, err := database.GetAllCoursesHandler()
	if err != nil {
		log.Println("Error getting all courses ", err)
		return nil
	}
	return courses
}

func GetCourseById(id string) *models.Course {
	var course *models.Course
	course, err := database.GetCourseByIdHandler(id)
	if err != nil {
		log.Println("Error getting course by id ", err)
		return nil
	}
	return course
}

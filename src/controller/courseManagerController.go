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

func GetCourseById(id string) {
	// search database for id
	// return course
}

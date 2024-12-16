package controller

import (
	"log"
	database "orkidslearning/src/database"
	models "orkidslearning/src/models/database"
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

func AddCourse(course models.AddCourse) (*models.Course, error) {
	addedCourse, err := database.AddCourseHandler(course)
	if err != nil {
		log.Println("Error adding course ", err)
		return nil, err
	}
	return addedCourse, nil
}

package controller

import (
	"log"

	"orkidslearning/src/database"
)

func EnrollInCourse(username string, courseId string) error {
	err := database.CheckIfUserExistsByUsername(username)
	if err != nil {
		log.Println("User does not exist", err)
		return err
	}

	err = database.CheckIfCourseExists(courseId)
	if err != nil {
		log.Println("Course does not exist", err)
		return err
	}

	isEnrolled, err := database.CheckIfUserIsEnrolledInCourse(username, courseId)
	if err != nil {
		log.Println("User is already enrolled in course", err)
		return err
	}

	if isEnrolled {
		return nil
	}

	err = database.AddUserToCourse(username, courseId)
	if err != nil {
		log.Println("Failed to add user to course", err)
		return err
	}
	return nil
}

func IsUserEnrolledInCourse(username string, courseId string) (bool, error) {
	isEnrolled, err := database.CheckIfUserIsEnrolledInCourse(username, courseId)
	if err != nil {
		log.Println("User is already enrolled in course", err)
		return false, err
	}

	return isEnrolled, nil
}

func UnenrollFromCourse(username string, courseId string) error {
	err := database.CheckIfUserExistsByUsername(username)
	if err != nil {
		log.Println("User does not exist", err)
		return err
	}

	err = database.CheckIfCourseExists(courseId)
	if err != nil {
		log.Println("Course does not exist", err)
		return err
	}

	isEnrolled, err := database.CheckIfUserIsEnrolledInCourse(username, courseId)
	if err != nil {
		log.Println("User is already enrolled in course", err)
		return err
	}

	if !isEnrolled {
		return nil
	}

	err = database.RemoveUserFromCourse(username, courseId)
	if err != nil {
		log.Println("Failed to remove user from course", err)
		return err
	}
	return nil
}

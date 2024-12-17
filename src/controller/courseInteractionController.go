package controller

import (
	"context"
	"log"
	services "orkidslearning/src/services"

	"go.opentelemetry.io/otel"
)

func EnrollInCourse(ctx context.Context, contextService *services.ContextService, username string, courseId string) error {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "EnrollInCourse")
	defer span.End()

	_, userSpan := tracer.Start(ctx, "CheckIfUserExistsByUsername")
	err := contextService.GetPostgres().CheckIfUserExistsByUsername(username)
	userSpan.End()
	if err != nil {
		log.Println("User does not exist", err)
		return err
	}

	_, courseSpan := tracer.Start(ctx, "CheckIfCourseExists")
	err = contextService.GetPostgres().CheckIfCourseExists(courseId)
	courseSpan.End()
	if err != nil {
		log.Println("Course does not exist", err)
		return err
	}

	_, isEnrolledSpan := tracer.Start(ctx, "CheckIfUserIsEnrolledInCourse")
	isEnrolled, err := contextService.GetPostgres().CheckIfUserIsEnrolledInCourse(username, courseId)
	isEnrolledSpan.End()
	if err != nil {
		log.Println("User is already enrolled in course", err)
		return err
	}

	if isEnrolled {
		return nil
	}

	_, addUserToCourseSpan := tracer.Start(ctx, "AddUserToCourse")
	err = contextService.GetPostgres().AddUserToCourse(username, courseId)
	addUserToCourseSpan.End()
	if err != nil {
		log.Println("Failed to add user to course", err)
		return err
	}
	return nil
}

func IsUserEnrolledInCourse(ctx context.Context, contextService *services.ContextService, username string, courseId string) (bool, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "IsUserEnrolledInCourse")
	defer span.End()

	_, isEnrolledSpan := tracer.Start(ctx, "CheckIfUserIsEnrolledInCourse")
	isEnrolled, err := contextService.GetPostgres().CheckIfUserIsEnrolledInCourse(username, courseId)
	isEnrolledSpan.End()
	if err != nil {
		log.Println("User is already enrolled in course", err)
		return false, err
	}

	return isEnrolled, nil
}

func UnenrollFromCourse(ctx context.Context, contextService *services.ContextService, username string, courseId string) error {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "UnenrollFromCourse")
	defer span.End()

	_, userSpan := tracer.Start(ctx, "CheckIfUserExistsByUsername")
	err := contextService.GetPostgres().CheckIfUserExistsByUsername(username)
	userSpan.End()
	if err != nil {
		log.Println("User does not exist", err)
		return err
	}

	_, courseSpan := tracer.Start(ctx, "CheckIfCourseExists")
	err = contextService.GetPostgres().CheckIfCourseExists(courseId)
	courseSpan.End()
	if err != nil {
		log.Println("Course does not exist", err)
		return err
	}

	_, isEnrolledSpan := tracer.Start(ctx, "CheckIfUserIsEnrolledInCourse")
	isEnrolled, err := contextService.GetPostgres().CheckIfUserIsEnrolledInCourse(username, courseId)
	isEnrolledSpan.End()
	if err != nil {
		log.Println("User is already enrolled in course", err)
		return err
	}

	if !isEnrolled {
		return nil
	}

	_, removeUserFromCourseSpan := tracer.Start(ctx, "RemoveUserFromCourse")
	err = contextService.GetPostgres().RemoveUserFromCourse(username, courseId)
	removeUserFromCourseSpan.End()
	if err != nil {
		log.Println("Failed to remove user from course", err)
		return err
	}
	return nil
}

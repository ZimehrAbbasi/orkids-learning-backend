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

	err := contextService.GetPostgres().CheckIfUserExistsByUsername(ctx, username)
	if err != nil {
		log.Println("User does not exist", err)
		return err
	}

	err = contextService.GetPostgres().CheckIfCourseExists(ctx, courseId)
	if err != nil {
		log.Println("Course does not exist", err)
		return err
	}

	isEnrolled, err := contextService.GetPostgres().CheckIfUserIsEnrolledInCourse(ctx, username, courseId)
	if err != nil {
		log.Println("User is already enrolled in course", err)
		return err
	}

	if isEnrolled {
		return nil
	}

	err = contextService.GetPostgres().AddUserToCourse(ctx, username, courseId)
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

	isEnrolled, err := contextService.GetPostgres().CheckIfUserIsEnrolledInCourse(ctx, username, courseId)
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

	err := contextService.GetPostgres().CheckIfUserExistsByUsername(ctx, username)
	if err != nil {
		log.Println("User does not exist", err)
		return err
	}

	err = contextService.GetPostgres().CheckIfCourseExists(ctx, courseId)
	if err != nil {
		log.Println("Course does not exist", err)
		return err
	}

	isEnrolled, err := contextService.GetPostgres().CheckIfUserIsEnrolledInCourse(ctx, username, courseId)
	if err != nil {
		log.Println("User is already enrolled in course", err)
		return err
	}

	if !isEnrolled {
		return nil
	}

	err = contextService.GetPostgres().RemoveUserFromCourse(ctx, username, courseId)
	if err != nil {
		log.Println("Failed to remove user from course", err)
		return err
	}
	return nil
}

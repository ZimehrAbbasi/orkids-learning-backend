package controller

import (
	"context"
	"fmt"
	"log"
	services "orkidslearning/src/services"

	"go.opentelemetry.io/otel"
)

func EnrollInCourse(ctx context.Context, contextService *services.ContextService, username string, courseId string) error {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "EnrollInCourse")
	defer span.End()

	exists, err := contextService.GetPostgres().DoesUserExistsByUsername(ctx, username)
	if err != nil {
		log.Println("User does not exist", err)
		return err
	}
	if !exists {
		return fmt.Errorf("user does not exist")
	}

	exists, err = contextService.GetPostgres().DoesCourseExist(ctx, courseId)
	if err != nil {
		log.Println("Course does not exist", err)
		return err
	}
	if !exists {
		return fmt.Errorf("course does not exist")
	}

	isEnrolled, err := contextService.GetPostgres().IsUserEnrolledInCourse(ctx, username, courseId)
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

	isEnrolled, err := contextService.GetPostgres().IsUserEnrolledInCourse(ctx, username, courseId)
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

	exists, err := contextService.GetPostgres().DoesUserExistsByUsername(ctx, username)
	if err != nil {
		log.Println("User does not exist", err)
		return err
	}
	if !exists {
		return fmt.Errorf("user does not exist")
	}

	exists, err = contextService.GetPostgres().DoesCourseExist(ctx, courseId)
	if err != nil {
		log.Println("Course does not exist", err)
		return err
	}
	if !exists {
		return fmt.Errorf("course does not exist")
	}

	isEnrolled, err := contextService.GetPostgres().IsUserEnrolledInCourse(ctx, username, courseId)
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

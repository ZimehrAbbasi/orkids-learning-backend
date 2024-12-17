package controller

import (
	"context"
	"log"
	models "orkidslearning/src/models/database"
	services "orkidslearning/src/services"

	"go.opentelemetry.io/otel"
)

func GetAllCourses(ctx context.Context, contextService *services.ContextService) ([]models.CoursePostgres, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "GetAllCourses")
	defer span.End()

	var courses []models.CoursePostgres
	_, coursesSpan := tracer.Start(ctx, "GetAllCoursesFromDatabase")
	courses, err := contextService.GetPostgres().GetAllCoursesFromDatabase()
	coursesSpan.End()
	if err != nil {
		log.Println("Error getting all courses ", err)
		return nil, err
	}
	return courses, nil
}

func GetCourseById(ctx context.Context, contextService *services.ContextService, id string) (*models.CoursePostgres, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "GetCourseById")
	defer span.End()

	var course *models.CoursePostgres
	_, courseSpan := tracer.Start(ctx, "GetCourseByIdFromDatabase")
	course, err := contextService.GetPostgres().GetCourseByIdFromDatabase(id)
	courseSpan.End()
	if err != nil {
		log.Println("Error getting course by id ", err)
		return nil, err
	}
	return course, nil
}

func AddCourse(ctx context.Context, contextService *services.ContextService, course models.AddCourse) (*models.CoursePostgres, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "AddCourse")
	defer span.End()

	_, addCourseSpan := tracer.Start(ctx, "AddCourseToDatabase")
	addedCourse, err := contextService.GetPostgres().AddCourseToDatabase(course)
	addCourseSpan.End()
	if err != nil {
		log.Println("Error adding course ", err)
		return nil, err
	}
	return addedCourse, nil
}

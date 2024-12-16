package controller

import (
	"context"
	"log"
	database "orkidslearning/src/database"
	models "orkidslearning/src/models/database"

	"go.opentelemetry.io/otel"
)

func GetAllCourses(ctx context.Context, db *database.Database) ([]models.Course, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "GetAllCourses")
	defer span.End()

	var courses []models.Course
	courses, err := db.GetAllCourses(ctx)
	if err != nil {
		log.Println("Error getting all courses ", err)
		return nil, err
	}
	return courses, nil
}

func GetCourseById(ctx context.Context, db *database.Database, id string) (*models.Course, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "GetCourseById")
	defer span.End()

	var course *models.Course
	course, err := db.GetCourseByID(ctx, id)
	if err != nil {
		log.Println("Error getting course by id ", err)
		return nil, err
	}
	return course, nil
}

func AddCourse(ctx context.Context, db *database.Database, course models.AddCourse) (*models.Course, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "AddCourse")
	defer span.End()

	addedCourse, err := db.AddCourse(ctx, course)
	if err != nil {
		log.Println("Error adding course ", err)
		return nil, err
	}
	return addedCourse, nil
}

package database

import (
	"context"
	"fmt"
	"log"
	"orkidslearning/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect(ctx context.Context, uri string) error {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
		return err
	}
	fmt.Println("Connected to MongoDB!")
	return nil
}

func GetAllCoursesHandler() ([]models.Course, error) {
	collection := client.Database("orkidslearning").Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch all documents
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Find error:", err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("Failed to close search cursor: ", err)
		}
	}(cursor, ctx)

	// Parse results
	var courses []models.Course
	if err = cursor.All(ctx, &courses); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return courses, nil
}

func GetCourseByIdHandler(id string) (*models.Course, error) {
	collection := client.Database("orkidslearning").Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectId: %v", err)
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	var course models.Course
	err = collection.FindOne(ctx, filter).Decode(&course)
	if err != nil {
		log.Println("FindOne error:", err)
		return nil, err
	}
	return &course, nil
}

func AddCourseHandler(course models.AddCourse) (*models.Course, error) {
	collection := client.Database("orkidslearning").Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var addedCourse models.Course
	result, err := collection.InsertOne(ctx, course)
	if err != nil {
		log.Println("InsertOne error:", err)
		return nil, err
	}

	addedCourse.Id = result.InsertedID.(primitive.ObjectID)
	addedCourse.Title = course.Title
	addedCourse.Description = course.Description

	return &addedCourse, nil
}

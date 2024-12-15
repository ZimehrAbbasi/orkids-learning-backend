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

func GetUserByEmail(email string) (*models.User, error) {
	collection := client.Database("orkidslearning").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CheckIfUserExists(username, email string) error {
	collection := client.Database("orkidslearning").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser models.User

	// Check email
	if err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser); err == nil {
		return fmt.Errorf("email already in use")
	}

	// Check username
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&existingUser); err == nil {
		return fmt.Errorf("username already in use")
	}

	return nil
}

func CheckIfUserExistsByUsername(username string) error {
	collection := client.Database("orkidslearning").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser models.User

	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
	if err != nil {
		return fmt.Errorf("user does not exist")
	}

	return nil
}

func AddUser(user models.AddUser) (*models.User, error) {
	collection := client.Database("orkidslearning").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert into MongoDB
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	// Return the newly created user
	return &models.User{
		ID:       result.InsertedID.(primitive.ObjectID),
		Username: user.Username,
		Email:    user.Email,
		Password: "", // Do not return the password
	}, nil
}

func CheckIfCourseExists(courseId string) error {
	collection := client.Database("orkidslearning").Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Checking if course exists", courseId)
	objectId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		log.Printf("Invalid ObjectId: %v", err)
		return err
	}

	var course models.Course
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&course)
	if err != nil {
		return fmt.Errorf("course does not exist")
	}
	return nil
}

func CheckIfUserIsEnrolledInCourse(username string, courseId string) (bool, error) {
	collection := client.Database("orkidslearning").Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		log.Printf("Invalid ObjectId: %v", err)
		return false, err
	}

	var course models.Course
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&course)
	if err != nil {
		return false, err
	}

	for _, user := range course.EnrolledUsers {
		if user == username {
			return true, nil
		}
	}

	return false, nil
}

func AddUserToCourse(username string, courseId string) error {
	collection := client.Database("orkidslearning").Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		log.Printf("Invalid ObjectId: %v", err)
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$push": bson.M{"enrolledUsers": username}})
	if err != nil {
		log.Println("UpdateOne error:", err)
		return err
	}

	return nil
}

func RemoveUserFromCourse(username string, courseId string) error {
	collection := client.Database("orkidslearning").Collection("courses")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		log.Printf("Invalid ObjectId: %v", err)
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$pull": bson.M{"enrolledUsers": username}})
	if err != nil {
		log.Println("UpdateOne error:", err)
		return err
	}

	return nil
}

package database

import (
	"context"
	"fmt"
	"log"
	"slices"

	models "orkidslearning/src/models/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database encapsulates the MongoDB client and provides methods to interact with the database
type Database struct {
	client     *mongo.Client
	dbName     string
	courseColl string
	userColl   string
}

// NewDatabase creates a new Database instance
func NewDatabase(ctx context.Context, uri, dbName string) (*Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")

	return &Database{
		client:     client,
		dbName:     dbName,
		courseColl: "courses",
		userColl:   "users",
	}, nil
}

// Disconnect closes the database connection
func (db *Database) Disconnect(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}

// GetAllCourses retrieves all courses
func (db *Database) GetAllCourses(ctx context.Context) ([]models.Course, error) {
	collection := db.client.Database(db.dbName).Collection(db.courseColl)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Find error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var courses []models.Course
	if err = cursor.All(ctx, &courses); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return courses, nil
}

// GetCourseByID retrieves a course by its ID
func (db *Database) GetCourseByID(ctx context.Context, id string) (*models.Course, error) {
	collection := db.client.Database(db.dbName).Collection(db.courseColl)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectId: %v", err)
		return nil, err
	}

	var course models.Course
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&course)
	if err != nil {
		log.Println("FindOne error:", err)
		return nil, err
	}
	return &course, nil
}

// AddCourse adds a new course
func (db *Database) AddCourse(ctx context.Context, course models.AddCourse) (*models.Course, error) {
	collection := db.client.Database(db.dbName).Collection(db.courseColl)

	result, err := collection.InsertOne(ctx, course)
	if err != nil {
		log.Println("InsertOne error:", err)
		return nil, err
	}

	return &models.Course{
		Id:          result.InsertedID.(primitive.ObjectID),
		Title:       course.Title,
		Description: course.Description,
	}, nil
}

// GetUserByEmail retrieves a user by email
func (db *Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := db.client.Database(db.dbName).Collection(db.userColl)

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckIfUserExists checks if a user exists by username or email
func (db *Database) CheckIfUserExists(ctx context.Context, username, email string) error {
	collection := db.client.Database(db.dbName).Collection(db.userColl)

	if err := collection.FindOne(ctx, bson.M{"email": email}).Err(); err == nil {
		return fmt.Errorf("email already in use")
	}

	if err := collection.FindOne(ctx, bson.M{"username": username}).Err(); err == nil {
		return fmt.Errorf("username already in use")
	}

	return nil
}

// AddUser adds a new user
func (db *Database) AddUser(ctx context.Context, user models.AddUser) (*models.User, error) {
	collection := db.client.Database(db.dbName).Collection(db.userColl)

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	return &models.User{
		ID:       result.InsertedID.(primitive.ObjectID),
		Username: user.Username,
		Email:    user.Email,
		Password: "", // Do not return the password
	}, nil
}

// AddUserToCourse enrolls a user in a course
func (db *Database) AddUserToCourse(ctx context.Context, username, courseId string) error {
	collection := db.client.Database(db.dbName).Collection(db.courseColl)

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

func (db *Database) CheckIfUserExistsByUsername(ctx context.Context, username string) error {
	collection := db.client.Database(db.dbName).Collection(db.userColl)

	err := collection.FindOne(ctx, bson.M{"username": username}).Err()
	if err != nil {
		return fmt.Errorf("user does not exist")
	}
	return nil
}

func (db *Database) CheckIfCourseExists(ctx context.Context, courseId string) error {
	collection := db.client.Database(db.dbName).Collection(db.courseColl)

	objectId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		log.Printf("Invalid ObjectId: %v", err)
		return err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Err()
	if err != nil {
		return fmt.Errorf("course does not exist")
	}
	return nil
}

func (db *Database) CheckIfUserIsEnrolledInCourse(ctx context.Context, username, courseId string) (bool, error) {
	collection := db.client.Database(db.dbName).Collection(db.courseColl)

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
	return slices.Contains(course.EnrolledUsers, username), nil
}

func (db *Database) RemoveUserFromCourse(ctx context.Context, username, courseId string) error {
	collection := db.client.Database(db.dbName).Collection(db.courseColl)

	objectId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		log.Printf("Invalid ObjectId: %v", err)
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$pull": bson.M{"enrolledUsers": username}})
	if err != nil {
		return err
	}
	return nil
}

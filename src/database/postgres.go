package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	models "orkidslearning/src/models/database"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

type PostgresDatabase struct {
	conn *pgx.Conn
}

// NewDatabase creates a new Database instance
func NewPostgresDatabase(ctx context.Context, connConfig pgx.ConnConfig) (*PostgresDatabase, error) {
	conn, err := pgx.Connect(connConfig)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return nil, err
	}
	fmt.Println("Connected to PostgreSQL!")
	return &PostgresDatabase{conn: conn}, nil
}

// Disconnect closes the database connection
func (db *PostgresDatabase) Disconnect() error {
	return db.conn.Close()
}

// GetAllCoursesFromDatabase retrieves all courses
func (db *PostgresDatabase) GetAllCoursesFromDatabase() ([]models.CoursePostgres, error) {
	query := "SELECT id, title, description FROM courses"
	rows, err := db.conn.Query(query)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var courses []models.CoursePostgres
	for rows.Next() {
		var course models.CoursePostgres
		var id pgtype.UUID
		if err := rows.Scan(&id, &course.Title, &course.Description); err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		course.Id = fmt.Sprintf("%x", id.Bytes)
		courses = append(courses, course)
	}
	return courses, nil
}

// GetCourseByIdFromDatabase retrieves a course by its ID
func (db *PostgresDatabase) GetCourseByIdFromDatabase(courseId string) (*models.CoursePostgres, error) {
	query := "SELECT id, title, description FROM courses WHERE id = $1"
	var course models.CoursePostgres
	var id pgtype.UUID
	err := db.conn.QueryRow(query, courseId).Scan(&id, &course.Title, &course.Description)
	if err != nil {
		log.Println("QueryRow error:", err)
		return nil, err
	}
	course.Id = fmt.Sprintf("%x", id.Bytes)
	return &course, nil
}

// AddCourseToDatabase adds a new course
func (db *PostgresDatabase) AddCourseToDatabase(course models.AddCourse) (*models.CoursePostgres, error) {
	query := "INSERT INTO courses (title, description) VALUES ($1, $2) RETURNING id"
	var id pgtype.UUID
	err := db.conn.QueryRow(query, course.Title, course.Description).Scan(&id)
	if err != nil {
		log.Println("Insert error:", err)
		return nil, err
	}
	return &models.CoursePostgres{
		Id:          fmt.Sprintf("%x", id.Bytes),
		Title:       course.Title,
		Description: course.Description,
	}, nil
}

// GetUserByEmail retrieves a user by email
func (db *PostgresDatabase) GetUserByEmail(email string) (*models.UserPostgres, error) {
	query := "SELECT id, username, email, password FROM users WHERE email = $1"
	var user models.UserPostgres

	// Use a temporary variable if needed for type conversion
	var id int
	err := db.conn.QueryRow(query, email).Scan(&id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}

	// Convert id to string if necessary
	user.Id = strconv.Itoa(id)

	return &user, nil
}

// CheckIfUserExists checks if a user exists by username or email
func (db *PostgresDatabase) CheckIfUserExists(username, email string) error {
	query := "SELECT 1 FROM users WHERE username = $1 OR email = $2"
	var exists int
	err := db.conn.QueryRow(query, username, email).Scan(&exists)
	if err == pgx.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.New("user already exists")
}

func (db *PostgresDatabase) CheckIfUserExistsByUsername(username string) error {
	query := "SELECT 1 FROM users WHERE username = $1"
	var exists int
	err := db.conn.QueryRow(query, username).Scan(&exists)
	if err == pgx.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.New("user already exists")
}

func (db *PostgresDatabase) CheckIfCourseExists(courseId string) error {
	query := "SELECT 1 FROM courses WHERE id = $1"
	var exists int
	err := db.conn.QueryRow(query, courseId).Scan(&exists)
	if err == pgx.ErrNoRows {
		// Return an error if the course does not exist
		return fmt.Errorf("course with ID '%s' does not exist", courseId)
	}
	if err != nil {
		// Handle unexpected errors
		return fmt.Errorf("error checking course existence: %w", err)
	}
	return nil
}

func (db *PostgresDatabase) CheckIfUserIsEnrolledInCourse(username, courseId string) (bool, error) {
	query := "SELECT 1 FROM course_enrollments WHERE username = $1 AND id = $2"
	var exists int
	err := db.conn.QueryRow(query, username, courseId).Scan(&exists)
	if err == pgx.ErrNoRows {
		// User is not enrolled in the course
		return false, nil
	}
	if err != nil {
		// Handle unexpected errors
		return false, fmt.Errorf("error checking user enrollment: %w", err)
	}
	// User is enrolled in the course
	return exists == 1, nil
}

// AddUser adds a new user
func (db *PostgresDatabase) AddUser(user models.AddUser) (*models.UserPostgres, error) {
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	var id int
	err := db.conn.QueryRow(query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}
	return &models.UserPostgres{
		Id:       strconv.Itoa(id),
		Username: user.Username,
		Email:    user.Email,
		Password: "", // Do not return the password
	}, nil
}

// AddUserToCourse enrolls a user in a course
func (db *PostgresDatabase) AddUserToCourse(username, courseId string) error {
	query := "INSERT INTO course_enrollments (username, id) VALUES ($1, $2)"
	_, err := db.conn.Exec(query, username, courseId)
	if err != nil {
		log.Println("Insert error:", err)
		return err
	}
	return nil
}

// RemoveUserFromCourse removes a user from a course
func (db *PostgresDatabase) RemoveUserFromCourse(username, courseId string) error {
	query := "DELETE FROM course_enrollments WHERE username = $1 AND id = $2"
	_, err := db.conn.Exec(query, username, courseId)
	if err != nil {
		return err
	}
	return nil
}

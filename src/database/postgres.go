package database

import (
	"context"
	"fmt"
	"log"
	"strconv"

	models "orkidslearning/src/models/database"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"go.opentelemetry.io/otel"
)

type PostgresDatabase struct {
	conn *pgx.Conn
}

// --- Initialization and Connection Management ---

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

// --- Course Management ---

// GetAllCourses retrieves all courses
func (db *PostgresDatabase) GetAllCourses(ctx context.Context) ([]models.CoursePostgres, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "GetAllCourses")
	defer span.End()

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

// GetCourseByID retrieves a course by its ID
func (db *PostgresDatabase) GetCourseByID(ctx context.Context, courseId string) (*models.CoursePostgres, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "GetCourseByID")
	defer span.End()

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

// AddCourse adds a new course
func (db *PostgresDatabase) AddCourse(ctx context.Context, course models.AddCourse) (*models.CoursePostgres, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "AddCourse")
	defer span.End()

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

// --- User Management ---

// GetUserByEmail retrieves a user by email
func (db *PostgresDatabase) GetUserByEmail(ctx context.Context, email string) (*models.UserPostgres, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "GetUserByEmail")
	defer span.End()

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

// AddUser adds a new user
func (db *PostgresDatabase) AddUser(ctx context.Context, user models.AddUser) (*models.UserPostgres, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "AddUser")
	defer span.End()

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
func (db *PostgresDatabase) AddUserToCourse(ctx context.Context, username, courseId string) error {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "AddUserToCourse")
	defer span.End()

	query := "INSERT INTO course_enrollments (username, id) VALUES ($1, $2)"
	_, err := db.conn.Exec(query, username, courseId)
	if err != nil {
		log.Println("Insert error:", err)
		return err
	}
	return nil
}

// RemoveUserFromCourse removes a user from a course
func (db *PostgresDatabase) RemoveUserFromCourse(ctx context.Context, username, courseId string) error {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "RemoveUserFromCourse")
	defer span.End()

	query := "DELETE FROM course_enrollments WHERE username = $1 AND id = $2"
	_, err := db.conn.Exec(query, username, courseId)
	if err != nil {
		return err
	}
	return nil
}

// --- Existence Checks ---

// DoesUserExist checks if a user exists by username or email
func (db *PostgresDatabase) DoesUserExist(ctx context.Context, username, email string) (bool, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "DoesUserExist")
	defer span.End()

	query := "SELECT 1 FROM users WHERE username = $1 OR email = $2"
	var exists int
	err := db.conn.QueryRow(query, username, email).Scan(&exists)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *PostgresDatabase) DoesUserExistsByUsername(ctx context.Context, username string) (bool, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "DoesUserExistsByUsername")
	defer span.End()

	query := "SELECT 1 FROM users WHERE username = $1"
	var exists int
	err := db.conn.QueryRow(query, username).Scan(&exists)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *PostgresDatabase) DoesCourseExist(ctx context.Context, courseId string) (bool, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "DoesCourseExist")
	defer span.End()

	query := "SELECT 1 FROM courses WHERE id = $1"
	var exists int
	err := db.conn.QueryRow(query, courseId).Scan(&exists)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking course existence: %w", err)
	}
	return true, nil
}

func (db *PostgresDatabase) IsUserEnrolledInCourse(ctx context.Context, username, courseId string) (bool, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "IsUserEnrolledInCourse")
	defer span.End()

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

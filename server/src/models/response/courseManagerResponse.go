package response

import (
	models "orkidslearning/src/models/database"
)

type GetCoursesResponse struct {
	Message string          `json:"message"`
	Error   string          `json:"error"`
	Courses []models.Course `json:"courses"`
}

type GetCourseResponse struct {
	Message  string        `json:"message"`
	Error    string        `json:"error"`
	Course   models.Course `json:"course"`
	Enrolled bool          `json:"enrolled" default:"false"`
}

type AddCourseResponse struct {
	Message string        `json:"message"`
	Error   string        `json:"error"`
	Course  models.Course `json:"course"`
	Added   bool          `json:"added" default:"false"`
}

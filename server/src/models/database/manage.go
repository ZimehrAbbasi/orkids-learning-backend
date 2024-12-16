package models

type EnrollInCourse struct {
	Username        string `json:"username"`
	CheckEnrollment bool   `json:"checkEnrollment" default:"false"`
}

type UnenrollFromCourse struct {
	Username string `json:"username"`
}

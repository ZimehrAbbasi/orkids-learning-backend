package models

type EnrollInCourse struct {
	Username string `json:"username"`
}

type UnenrollFromCourse struct {
	Username string `json:"username"`
}

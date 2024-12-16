package response

type EnrollInCourseResponse struct {
	Message  string `json:"message"`
	Error    string `json:"error"`
	Enrolled bool   `json:"enrolled" default:"false"`
	Username string `json:"username"`
	CourseId string `json:"courseId"`
}

type UnenrollFromCourseResponse struct {
	Message    string `json:"message"`
	Error      string `json:"error"`
	Unenrolled bool   `json:"unenrolled" default:"false"`
	Username   string `json:"username"`
	CourseId   string `json:"courseId"`
}

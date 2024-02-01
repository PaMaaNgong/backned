package main

type Repository interface {
	GetCourses(query string, limit int, offset int) ([]CourseOverview, error)
	GetCourseDetail(id string) (CourseDetail, error)
	GetReviewsOverview(id string, limit int, offset int) ([]ReviewOverview, error)
	GetReviewsDetail(courseId string, limit int, offset int) ([]ReviewDetail, error)
	CreateReview(courseId string, review ReviewDetail) error
	AddCourse(course CourseDetail)
}

package main

type Repository interface {
	GetCourses(query string, limit int, offset int) ([]CourseOverview, error)
	GetCourseDetail(id string) (CourseDetail, error)
	AddCourse(course CourseDetail)
	GetCourseGrades(id string) ([]Grade, error)
	GetReviewsOverview(id string, limit int, offset int) ([]ReviewOverview, error)
	GetReviewsDetail(courseId string, limit int, offset int) ([]ReviewDetail, error)
	CreateReview(userId uint64, courseId string, review ReviewDetail) error
	EditReview(userId uint64, courseId string, reviewId uint64, review ReviewDetail) error
	DeleteReview(userId uint64, courseId string, reviewId uint64) error
}

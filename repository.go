package main

type Repository interface {
	GetCourses(query string) ([]CourseOverview, error)
	GetCourseDetail(id string) (CourseDetail, error)
	AddCourse(course CourseDetail)
	GetCourseGrades(id string) ([]Grade, error)
	GetReviewsOverview(id string) ([]ReviewOverview, error)
	GetReviewsDetail(courseId string) ([]ReviewDetail, error)
	CreateReview(userId uint64, courseId string, review ReviewDetail) error
	EditReview(userId uint64, courseId string, reviewId uint64, review ReviewDetail) error
	DeleteReview(userId uint64, courseId string, reviewId uint64) error
	GetCoursesByUser(userId uint64) ([]CourseOverview, error)
	GetReviewByUser(userId uint64, courseId string) (ReviewDetail, error)
}

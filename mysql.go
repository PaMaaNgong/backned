package main

import (
	"errors"
	"fmt"
	"math"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLRepository struct {
	db *gorm.DB
}

func NewMySQLRepository(host string, username string, password string, port string, dbName string) MySQLRepository {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName)
	driver := mysql.Open(dsn)
	db, _ := gorm.Open(driver, &gorm.Config{})
	err := db.AutoMigrate(&CourseDetail{}, &ReviewDetail{}, &User{})
	if err != nil {
		panic(err)
	}
	return MySQLRepository{
		db: db,
	}
}

func (s MySQLRepository) GetCourses(query string) ([]CourseOverview, error) {
	coursesDetail := make([]CourseDetail, 0)
	s.db.Where("id <> ?", query).
		Find(&coursesDetail)
	coursesOverview := make([]CourseOverview, 0)
	for _, courseDetail := range coursesDetail {
		coursesOverview = append(coursesOverview, courseDetail.CourseOverview)
	}
	return coursesOverview, nil
}

func (s MySQLRepository) GetCourseDetail(courseId string) (CourseDetail, error) {
	var courseDetail CourseDetail
	result := s.db.Where("id = ?", courseId).Find(&courseDetail)
	if result.RowsAffected == 0 {
		return CourseDetail{}, ErrCourseNotFound{}
	}
	return courseDetail, nil
}

func (s MySQLRepository) AddCourse(course CourseDetail) {
	s.db.Create(&course)
}

func (s MySQLRepository) GetCourseGrades(id string) ([]Grade, error) {
	if s.noCourse(id) {
		return []Grade{}, ErrCourseNotFound{}
	}

	reviewsDetail := make([]ReviewDetail, 0)
	s.db.Model(&ReviewDetail{}).
		Where("course_id = ?", id).
		Find(&reviewsDetail)
	grades := make([]Grade, 0)
	for _, reviewDetail := range reviewsDetail {
		grades = append(grades, reviewDetail.Grade)
	}
	return grades, nil
}

func (s MySQLRepository) GetReviewsOverview(courseId string) ([]ReviewOverview, error) {
	if s.noCourse(courseId) {
		return []ReviewOverview{}, ErrCourseNotFound{}
	}

	reviewsDetail := make([]ReviewDetail, 0)
	s.db.Where("course_id = ?", courseId).
		Find(&reviewsDetail)
	reviewsOverview := make([]ReviewOverview, 0)
	for _, reviewDetail := range reviewsDetail {
		reviewsOverview = append(reviewsOverview, reviewDetail.ReviewOverview)
	}
	return reviewsOverview, nil
}

func (s MySQLRepository) GetReviewsDetail(courseId string) ([]ReviewDetail, error) {
	if s.noCourse(courseId) {
		return []ReviewDetail{}, ErrCourseNotFound{}
	}
	reviewsDetail := make([]ReviewDetail, 0)
	result := s.db.Where("course_id = ?", courseId).Find(&reviewsDetail)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []ReviewDetail{}, ErrCourseNotFound{}
	}
	return reviewsDetail, nil
}

func (s MySQLRepository) CreateReview(userId uint64, courseId string, review ReviewDetail) error {
	if s.noCourse(courseId) {
		return ErrCourseNotFound{}
	}
	review.CourseID = courseId
	review.OwnerID = userId
	result := s.db.Create(&review)
	s.updateCourse(courseId)
	return result.Error
}

func (s MySQLRepository) EditReview(userId uint64, courseId string, reviewId uint64, review ReviewDetail) error {
	if s.noReview(courseId, reviewId) {
		return ErrCourseNotFound{}
	}
	review.ID = reviewId
	review.CourseID = courseId
	review.OwnerID = userId
	s.db.Where("id = ? AND course_id = ? AND owner_id = ?", reviewId, courseId, userId).Updates(review)
	s.updateCourse(courseId)
	return nil
}

func (s MySQLRepository) DeleteReview(userId uint64, courseId string, reviewId uint64) error {
	if s.noReview(courseId, reviewId) {
		return ErrCourseNotFound{}
	}
	s.db.Where("id = ? AND course_id = ? AND owner_id = ?", reviewId, courseId, userId).Delete(&ReviewDetail{})
	s.updateCourse(courseId)
	return nil
}

func (s MySQLRepository) GetCoursesByUser(userId uint64) ([]CourseOverview, error) {
	reviews := make([]ReviewDetail, 0)
	courses := make([]CourseOverview, 0)
	s.db.Where("owner_id = ?", userId).Find(&reviews)
	for _, review := range reviews {
		course := CourseDetail{}
		s.db.Where("id = ?", review.CourseID).First(&course)
		courses = append(courses, course.CourseOverview)
	}
	return courses, nil
}

func (s MySQLRepository) GetReviewByUser(userId uint64, courseId string) (ReviewDetail, error) {
	var review ReviewDetail
	if s.db.Where("course_id = ? AND owner_id = ?", courseId, userId).First(&review).Error == nil {
		return review, nil
	}
	return ReviewDetail{}, ErrCourseNotFound{}
}

func (s MySQLRepository) noCourse(courseId string) bool {
	var course CourseDetail
	return errors.Is(s.db.Where("id = ?", courseId).First(&course).Error, gorm.ErrRecordNotFound)
}

func (s MySQLRepository) noReview(courseId string, reviewId uint64) bool {
	var review ReviewDetail
	return errors.Is(s.db.Where("id = ? AND course_id = ?", reviewId, courseId).First(&review).Error, gorm.ErrRecordNotFound)
}

func (s MySQLRepository) updateCourse(courseId string) {
	var new_rating float64
	var total_reviews uint64
	reviewsDetail := make([]ReviewDetail, 0)
	s.db.Where("course_id = ?", courseId).Find(&reviewsDetail)
	for _, review := range reviewsDetail {
		new_rating += float64(review.Rating)
		total_reviews += 1
	}
	new_rating /= float64(total_reviews)
	if math.IsNaN(new_rating) {
		new_rating = 0
	}
	new_rating = (math.Round(new_rating * 100)) / 100
	s.db.Model(&CourseDetail{}).Where("id = ?", courseId).Update("rating", new_rating).Update("total_reviews", total_reviews)
}

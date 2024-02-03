package main

import (
	"errors"
	"fmt"

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
	err := db.AutoMigrate(&CourseDetail{}, &ReviewDetail{})
	if err != nil {
		panic(err)
	}
	return MySQLRepository{
		db: db,
	}
}

func (s MySQLRepository) GetCourses(query string, limit int, offset int) ([]CourseOverview, error) {
	coursesDetail := make([]CourseDetail, 0)
	s.db.Limit(limit).
		Offset(offset).
		Where("id <> ?", query).
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

func (s MySQLRepository) GetReviewsOverview(courseId string, limit int, offset int) ([]ReviewOverview, error) {
	if s.noCourse(courseId) {
		return []ReviewOverview{}, ErrCourseNotFound{}
	}

	reviewsDetail := make([]ReviewDetail, 0)
	s.db.Limit(limit).
		Offset(offset).
		Where("course_id = ?", courseId).
		Find(&reviewsDetail)
	reviewsOverview := make([]ReviewOverview, 0)
	for _, reviewDetail := range reviewsDetail {
		reviewsOverview = append(reviewsOverview, reviewDetail.ReviewOverview)
	}
	return reviewsOverview, nil
}

func (s MySQLRepository) GetReviewsDetail(courseId string, limit int, offset int) ([]ReviewDetail, error) {
	if s.noCourse(courseId) {
		return []ReviewDetail{}, ErrCourseNotFound{}
	}

	reviewsDetail := make([]ReviewDetail, 0)
	result := s.db.Limit(limit).Offset(offset).Where("course_id = ?", courseId).Find(&reviewsDetail)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []ReviewDetail{}, ErrCourseNotFound{}
	}
	return reviewsDetail, nil
}

func (s MySQLRepository) CreateReview(courseId string, review ReviewDetail) error {
	if s.noCourse(courseId) {
		return ErrCourseNotFound{}
	}
	review.CourseID = courseId
	result := s.db.Create(&review)
	if result.Error == nil {
		var course CourseDetail
		s.db.First(&course)
		s.db.Model(&CourseDetail{}).Where("id = ?", courseId).Update("total_reviews", course.TotalReviews+1)
	}
	return result.Error
}

func (s MySQLRepository) noCourse(courseId string) bool {
	var course CourseDetail
	return errors.Is(s.db.Where("id = ?", courseId).First(&course).Error, gorm.ErrRecordNotFound)
}

func (s MySQLRepository) AddCourse(course CourseDetail) {
	s.db.Create(&course)
}

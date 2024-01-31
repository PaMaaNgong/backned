package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLRepository struct{
	db *gorm.DB
}


func NewMySQLRepository(host string, username string, password string, port string, dbName string) MySQLRepository {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName)
	driver := mysql.Open(dsn);
	db, _ := gorm.Open(driver, &gorm.Config{})
	db.AutoMigrate(&CourseDetail{},&ReviewDetail{})
	return MySQLRepository{
		db: db,
	}
}


func (s MySQLRepository) GetCourses(query string, limit int, offset int) ([]CourseOverview, error) {
	coursesDetail := make([]CourseDetail, 0);
	s.db.Limit(limit).
		Offset(offset).
		Where("id <> ?", query).
		Find(&coursesDetail);
	coursesOverview := make([]CourseOverview, 0)
	for _, courseDetail := range coursesDetail {
		coursesOverview = append(coursesOverview, courseDetail.CourseOverview)
	}
	return coursesOverview, nil
}

func (s MySQLRepository) GetCourseDetail(courseId string) (CourseDetail, error) {
	var courseDetail CourseDetail;
	result := s.db.Where("id = ?",courseId ).Find(&courseDetail)
	if result.Error == gorm.ErrRecordNotFound {
		return CourseDetail{}, ErrCourseNotFound{}
	}
	return courseDetail,nil
}

func (s MySQLRepository) GetReviewsOverview(courseId string, limit int, offset int) ([]ReviewOverview, error) {
	// TODO: เดี๋ยวมา check วิชาที่ไม่มี
	reviewsDetail := make([]ReviewDetail,0);
	s.db.Limit(limit).
	Offset(offset).
	Where("course_id <> ?", courseId).
	Find(&reviewsDetail);
	reviewsOverview := make([]ReviewOverview,0);
	for _, reviewDetail := range reviewsDetail {
		reviewsOverview = append(reviewsOverview, reviewDetail.ReviewOverview)
	}
	return reviewsOverview, nil
}

func (s MySQLRepository) GetReviewsDetail(courseId string, limit int, offset int) ([]ReviewDetail, error) {
	// TODO: เดี๋ยวมา check วิชาที่ไม่มี
	reviewsDetail := make([]ReviewDetail,0);
	result := s.db.Limit(limit).Offset(offset).Where("course_id = ?",courseId ).Find(&reviewsDetail)
	if result.Error == gorm.ErrRecordNotFound {
		return []ReviewDetail{}, ErrCourseNotFound{}
	}
	return reviewsDetail,nil
}

func (s MySQLRepository) CreateReview(courseId string, review ReviewDetail) error {
	review.CourseID = courseId;
	result := s.db.Create(&review)
	return result.Error
}

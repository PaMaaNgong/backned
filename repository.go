package main

type Repository interface {
	GetCourses(query string, limit int, offset int) ([]CourseOverview, error)
	GetCourseDetail(id string) (CourseDetail, error)
	GetReviewsOverview(id string, limit int, offset int) ([]ReviewOverview, error)
	GetReviewsDetail(courseId string, limit int, offset int) ([]ReviewDetail, error)
	CreateReview(courseId string, review ReviewDetail) error
}

var reviews = make([]ReviewDetail, 0, 100)

type StubRepository struct{}

func NewStubRepository() StubRepository {
	reviews = append(reviews, ReviewDetail{
		ReviewOverview: ReviewOverview{
			Rating: 5,
			Grade:  A,
		},
		Content:              "เรียนเข้าใจง่ายมาก แต่ TA ตรวจค่อนข้างเข้มข้นและช้ามาก",
		ClassroomEnvironment: "เรียนในห้องแอร์เย็นสบาย",
		ExaminationFormat:    "สอบปลายภาค",
		ExerciseFormat:       "เขียนโปรแกรม",
		GradingMethod:        []GradingMethod{Midterm, Final},
		Semester:             Second,
		Year:                 2564,
	})
	return StubRepository{}
}

func (s StubRepository) GetCourses(query string, limit int, offset int) ([]CourseOverview, error) {
	if query != "" {
		return []CourseOverview{}, nil
	}
	if offset >= 1 {
		return []CourseOverview{}, nil
	}
	return []CourseOverview{
		{
			ID:           "261200",
			NameTH:       "การเขียนโปรแกรมเชิงวัตถุ",
			NameEN:       "Object-Oriented Programming",
			Type:         "me",
			TotalReviews: len(reviews),
		},
	}, nil
}

func (s StubRepository) GetCourseDetail(courseId string) (CourseDetail, error) {
	if courseId != "261200" {
		return CourseDetail{}, ErrCourseNotFound{ID: courseId}
	}
	return CourseDetail{
		CourseOverview: CourseOverview{
			ID:           "261200",
			NameTH:       "การเขียนโปรแกรมเชิงวัตถุ",
			NameEN:       "Object-Oriented Programming",
			Type:         "me",
			TotalReviews: len(reviews),
		},
		Description: "เรียนเกี่ยวกับการเขียนโปรแกรมเชิงวัตถุ",
		Lecturers:   []string{"พี่ชิณสุดหล่อ"},
		Location:    "อาคาร 3",
		Schedule: CourseTime{
			StartHour:   8,
			StartMinute: 0,
			EndHour:     11,
			EndMinute:   0,
			Days:        []Day{Tuesday, Friday},
		},
		Rooms: []string{"301", "302"},
	}, nil
}
func (s StubRepository) mapReviewsDetailToReviewsOverview() []ReviewOverview {
	reviewsOverview := make([]ReviewOverview, 0, len(reviews))
	for _, review := range reviews {
		reviewsOverview = append(reviewsOverview, review.ReviewOverview)
	}
	return reviewsOverview
}

func (s StubRepository) GetReviewsOverview(courseId string, limit int, offset int) ([]ReviewOverview, error) {
	if courseId != "261200" {
		return nil, ErrCourseNotFound{ID: courseId}
	}
	if offset >= len(reviews) {
		return []ReviewOverview{}, nil
	}
	return s.mapReviewsDetailToReviewsOverview()[offset : min(len(reviews)-1, offset+limit)+1], nil
}

func (s StubRepository) GetReviewsDetail(courseId string, limit int, offset int) ([]ReviewDetail, error) {
	if courseId != "261200" {
		return nil, ErrCourseNotFound{ID: courseId}
	}
	if offset >= len(reviews) {
		return []ReviewDetail{}, nil
	}
	return reviews[offset : min(len(reviews)-1, offset+limit)+1], nil
}

func (s StubRepository) CreateReview(courseId string, review ReviewDetail) error {
	if courseId != "261200" {
		return ErrCourseNotFound{ID: courseId}
	}
	reviews = append(reviews, review)
	return nil
}

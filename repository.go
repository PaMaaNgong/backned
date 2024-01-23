package main

type Repository interface {
	GetCourses(query string, limit int, offset int) ([]CourseOverview, error)
	GetCourseDetail(id string) (CourseDetail, error)
	GetReviewsOverview(id string, limit int, offset int) ([]ReviewOverview, error)
	GetReviewsDetail(courseId string, limit int, offset int) ([]ReviewDetail, error)
	CreateReview(courseId string, review ReviewDetail) error
}

var courses = []CourseOverview{
	{
		ID:     "001101",
		NameTH: "ภาษาอังกฤษพื้นฐาน 1",
		NameEN: "Fundamental English 1",
		Type:   FreeElective,
	},
	{
		ID:     "001102",
		NameTH: "ภาษาอังกฤษพื้นฐาน 2",
		NameEN: "Fundamental English 2",
		Type:   FreeElective,
	},
	{
		ID:     "001201",
		NameTH: "การอ่านและการเขียนอย่างมีประสิทธิภาพ",
		NameEN: "Critical Reading and Effective Writing",
		Type:   FreeElective,
	},
	{
		ID:     "001225",
		NameTH: "ภาษาอังกฤษเพื่อวิทยาศาสตร์และเทคโนโลยี",
		NameEN: "English in Science and Technology Context",
		Type:   FreeElective,
	},
	{
		ID:     "206161",
		NameTH: "คณิตศาสตร์เพื่อวิศวกรรม 1",
		NameEN: "Calculus for Engineering 1",
		Type:   GeneralElective,
	},
	{
		ID:     "206162",
		NameTH: "คณิตศาสตร์เพื่อวิศวกรรม 2",
		NameEN: "Calculus for Engineering 2",
		Type:   GeneralElective,
	},
	{
		ID:     "206261",
		NameTH: "คณิตศาสตร์เพื่อวิศวกรรม 3",
		NameEN: "Calculus for Engineering 3",
		Type:   GeneralElective,
	},
	{
		ID:     "207105",
		NameTH: "ฟิสิกส์เพื่อนักศึกษาวิศวกรรมและอุตสาหกรรมเกษตร 1",
		NameEN: "Physics for Engineering and Agro-Industry Students 1",
		Type:   GeneralElective,
	},
	{
		ID:     "207115",
		NameTH: "การทดลองฟิสิกส์เพื่อนักศึกษาวิศวกรรมและอุตสาหกรรมเกษตร 1",
		NameEN: "Physics Laboratory for Engineering and Agro-Industry Students 1",
		Type:   GeneralElective,
	},
	{
		ID:     "261200",
		NameTH: "การเขียนโปรแกรมเชิงวัตถุ",
		NameEN: "Object-Oriented Programming",
		Type:   MajorElective,
	},
	{
		ID:     "261408",
		NameTH: "คอมพิวเตอร์ควอนตัม",
		NameEN: "Quantum Computing",
		Type:   MajorElective,
	},
	{
		ID:     "261434",
		NameTH: "การออกแบบและการจัดการเครือข่าย",
		NameEN: "Network Design and Management",
		Type:   MajorElective,
	},
	{
		ID:     "261453",
		NameTH: "การประมวลผลภาพดิจิตอล",
		NameEN: "Digital Image Processing",
		Type:   MajorElective,
	},
	{
		ID:     "261456",
		NameTH: "การประมวลผลด้วยเทคนิคทางปัญญาประดิษฐ์",
		NameEN: "Intro to Computational Intelligence",
		Type:   MajorElective,
	},
	{
		ID:     "261494-1",
		NameTH: "R Programming (selected topics in CPE)",
		NameEN: "R Programming (selected topics in CPE)",
		Type:   MajorElective,
	},
	{
		ID:     "261494-10",
		NameTH: "Blockchain Programming (selected topics in CPE)",
		NameEN: "Blockchain Programming (selected topics in CPE)",
		Type:   MajorElective,
	},
	{
		ID:     "261494-2",
		NameTH: "Adv Algorithms and Computation (selected topics in CPE)",
		NameEN: "Adv Algorithms and Computation (selected topics in CPE)",
		Type:   MajorElective,
	},
	{
		ID:     "261494-4",
		NameTH: "Up Skills Unlock Limits (selected topics in CPE)",
		NameEN: "Up Skills Unlock Limits (selected topics in CPE)",
		Type:   MajorElective,
	},
	{
		ID:     "261494-8",
		NameTH: "Startup Raid (selected topics in CPE)",
		NameEN: "Startup Raid (selected topics in CPE)",
		Type:   MajorElective,
	},
}
var reviews = make(map[string][]ReviewDetail)

type StubRepository struct{}

func NewStubRepository() StubRepository {
	for _, course := range courses {
		reviews[course.ID] = []ReviewDetail{}
	}
	reviews["261200"] = append(reviews["261200"], ReviewDetail{
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
	lastIndex := len(courses) - 1
	for i, course := range courses {
		courses[i].TotalReviews = len(reviews[course.ID])
	}
	return courses[min(offset, lastIndex):min(offset+limit, lastIndex)], nil
}

func (s StubRepository) GetCourseDetail(courseId string) (CourseDetail, error) {
	for _, course := range courses {
		if course.ID == courseId {
			return CourseDetail{
				CourseOverview: course,
				Description:    "เรียนเกี่ยวกับการเขียนโปรแกรมเชิงวัตถุทุกอัน",
				Lecturers:      []string{"พี่ชิณสุดหล่อ"},
				Location:       "อาคาร 3",
				Schedule: CourseTime{
					StartHour:   8,
					StartMinute: 0,
					EndHour:     11,
					EndMinute:   0,
					Days:        []Day{Tuesday, Friday},
				},
				Rooms: []string{"301", "302"},
				Credit: CourseCredit{
					Lecture: 3,
					Lab:     1,
				},
			}, nil
		}
	}
	return CourseDetail{}, ErrCourseNotFound{ID: courseId}
}

func mapReviewDetailToReviewOverview(details []ReviewDetail) []ReviewOverview {
	var overviews []ReviewOverview
	for _, detail := range details {
		overviews = append(overviews, detail.ReviewOverview)
	}
	return overviews
}

func (s StubRepository) GetReviewsOverview(courseId string, limit int, offset int) ([]ReviewOverview, error) {
	if _, ok := reviews[courseId]; !ok {
		return nil, ErrCourseNotFound{ID: courseId}
	}
	if offset >= len(reviews[courseId]) {
		return []ReviewOverview{}, nil
	}
	review := reviews[courseId]
	lastIndex := len(review) - 1
	return mapReviewDetailToReviewOverview(review)[min(offset, lastIndex):min(offset+limit, lastIndex)], nil
}

func (s StubRepository) GetReviewsDetail(courseId string, limit int, offset int) ([]ReviewDetail, error) {
	if _, ok := reviews[courseId]; !ok {
		return nil, ErrCourseNotFound{ID: courseId}
	}
	if offset >= len(reviews[courseId]) {
		return []ReviewDetail{}, nil
	}
	lastIndex := len(reviews[courseId]) - 1
	return reviews[courseId][min(offset, lastIndex):min(offset+limit, lastIndex)], nil
}

func (s StubRepository) CreateReview(courseId string, review ReviewDetail) error {
	if _, ok := reviews[courseId]; !ok {
		return ErrCourseNotFound{ID: courseId}
	}
	reviews[courseId] = append(reviews[courseId], review)
	return nil
}

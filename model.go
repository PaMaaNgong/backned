package main

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

type Enum interface {
	IsValid() bool
}

func ValidateEnum(fl validator.FieldLevel) bool {
	enum, ok := fl.Field().Interface().(Enum)
	if !ok {
		return false
	}
	return enum.IsValid()
}

func ValidEnumSlice(fl validator.FieldLevel) bool {
	slice, err := tryMapToEnumSlice(fl.Field().Interface())
	if err != nil {
		return false
	}
	for _, e := range slice {
		if !e.IsValid() {
			return false
		}
	}
	return true
}

func tryMapToEnumSlice(maybeSlice any) ([]Enum, error) {
	if reflect.TypeOf(maybeSlice).Kind() != reflect.Slice {
		return nil, ErrNotEnum{Value: maybeSlice}
	}
	slice := reflect.ValueOf(maybeSlice)
	enumSlice := make([]Enum, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		enum, ok := slice.Index(i).Interface().(Enum)
		if !ok {
			return nil, ErrNotEnum{Value: slice.Index(i).Interface()}
		}
		enumSlice = append(enumSlice, enum)
	}
	return enumSlice, nil
}

func RegisterEnumValidator(v *validator.Validate) {
	_ = v.RegisterValidation("enum", ValidateEnum)
	_ = v.RegisterValidation("enum_slice", ValidEnumSlice)
}

type CourseType string

const (
	MajorElective   CourseType = "me"
	FreeElective               = "fe"
	GeneralElective            = "ge"
)

func (c CourseType) IsValid() bool {
	return c == MajorElective || c == FreeElective || c == GeneralElective
}

type Day string

const (
	Monday    Day = "mon"
	Tuesday       = "tue"
	Wednesday     = "wed"
	Thursday      = "thu"
	Friday        = "fri"
	Saturday      = "sat"
	Sunday        = "sun"
)

func (d Day) IsValid() bool {
	return d == Monday || d == Tuesday || d == Wednesday || d == Thursday || d == Friday || d == Saturday || d == Sunday
}

type Grade string

const (
	A     Grade = "a"
	BPlus       = "b+"
	B           = "b"
	CPlus       = "c+"
	C           = "c"
	DPlus       = "d+"
	D           = "d"
	F           = "f"
	Blank       = "-"
)

func (g Grade) IsValid() bool {
	return g == A || g == BPlus || g == B || g == CPlus || g == C || g == DPlus || g == D || g == F || g == Blank
}

type GradingMethod string

const (
	Midterm GradingMethod = "midterm"
	Final                 = "final"
)

func (g GradingMethod) IsValid() bool {
	return g == Midterm || g == Final
}

type Semester string

const (
	First  Semester = "1"
	Second          = "2"
	Summer          = "3"
)

func (s Semester) IsValid() bool {
	return s == First || s == Second || s == Summer
}

type CourseTime struct {
	StartHour   int   `json:"start_hour"`
	StartMinute int   `json:"start_minute"`
	EndHour     int   `json:"end_hour"`
	EndMinute   int   `json:"end_minute"`
	Days        []Day `json:"days"`
}

type CourseOverview struct {
	ID           string     `json:"id"`
	NameTH       string     `json:"name_th"`
	NameEN       string     `json:"name_en"`
	Type         CourseType `json:"type"`
	TotalReviews int        `json:"total_reviews"`
}

type CourseDetail struct {
	CourseOverview
	Description string     `json:"description"`
	Lecturers   []string   `json:"lecturers"`
	Location    string     `json:"location"`
	Schedule    CourseTime `json:"schedule"`
	Rooms       []string   `json:"rooms"`
}

type ReviewOverview struct {
	Rating int   `json:"rating" binding:"required"`
	Grade  Grade `json:"grade" binding:"required,enum"`
}

type ReviewDetail struct {
	ReviewOverview
	Content              string          `json:"content" binding:"required"`
	ClassroomEnvironment string          `json:"classroom_environment" binding:"required"`
	ExaminationFormat    string          `json:"examination_format" binding:"required"`
	ExerciseFormat       string          `json:"exercise_format" binding:"required"`
	GradingMethod        []GradingMethod `json:"grading_method" binding:"required,enum_slice"`
	Semester             Semester        `json:"semester" binding:"required,enum"`
	Year                 int             `json:"year" binding:"required"`
}

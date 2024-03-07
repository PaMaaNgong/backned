package main

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type TrimString string

func (t *TrimString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = TrimString(strings.TrimSpace(s))
	return nil
}

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
	A     Grade = "A"
	BPlus       = "B+"
	B           = "B"
	CPlus       = "C+"
	C           = "C"
	DPlus       = "D+"
	D           = "D"
	F           = "F"
	W           = "W"
	Blank       = "-"
)

func (g Grade) IsValid() bool {
	return g == A || g == BPlus || g == B || g == CPlus || g == C || g == DPlus || g == D || g == F || g == W || g == Blank
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

type Difficult string

const (
	Easy   Difficult = "easy"
	Normal           = "normal"
	Hard             = "hard"
)

func (d Difficult) IsValid() bool {
	return d == Easy || d == Normal || d == Hard
}

type ExerciseFormat string

const (
	Individual ExerciseFormat = "individual"
	Group                     = "group"
)

func (e ExerciseFormat) IsValid() bool {
	return e == Individual || e == Group
}

type ExaminationFormat string

const (
	Objective  ExaminationFormat = "objective"
	Subjective                   = "subjective"
)

func (e ExaminationFormat) IsValid() bool {
	return e == Objective || e == Subjective
}

type ExaminationInfo struct {
	Format     []ExaminationFormat `json:"format" binding:"enum_slice" gorm:"serializer:json"`
	Difficulty Difficult           `json:"difficulty" binding:"enum"`
}

type ExerciseInfo struct {
	Format     []ExerciseFormat `json:"format" binding:"enum_slice" gorm:"serializer:json"`
	Difficulty Difficult        `json:"difficulty" binding:"enum"`
}

type CourseOverview struct {
	ID           string     `json:"id" gorm:"type:varchar(10);primaryKey"`
	NameTH       string     `json:"name_th"`
	NameEN       string     `json:"name_en"`
	Type         CourseType `json:"type"`
	TotalReviews int        `json:"total_reviews"`
	Rating       float64    `json:"rating"`
}

type CourseDetail struct {
	CourseOverview
	Description string         `json:"description"`
	Lecturers   []string       `json:"lecturers" gorm:"serializer:json"`
	Credit      CourseCredit   `json:"credit" gorm:"serializer:json"`
	Reviews     []ReviewDetail `json:"-" gorm:"foreignkey:CourseID"`
}

type CourseCredit struct {
	Lecture int `json:"lecture"`
	Lab     int `json:"lab"`
}

type ReviewOverview struct {
	ID      uint64     `json:"id" gorm:"primary_key"`
	Rating  int        `json:"rating" binding:"required"`
	Grade   Grade      `json:"grade" binding:"required,enum"`
	Content TrimString `json:"content" binding:"required"`
}

type ReviewDetail struct {
	ReviewOverview
	ClassroomEnvironment TrimString      `json:"classroom_environment" binding:"required"`
	Examination          ExaminationInfo `json:"examination_format" binding:"required" gorm:"serializer:json"`
	Exercise             ExerciseInfo    `json:"exercise_format" binding:"required" gorm:"serializer:json"`
	GradingMethod        []GradingMethod `json:"grading_method" binding:"required,enum_slice" gorm:"serializer:json"`
	Semester             Semester        `json:"semester" binding:"required,enum"`
	Year                 int             `json:"year" binding:"required"`
	Other                TrimString      `json:"other"`
	OwnerID              uint64          `json:"owner"`
	CourseID             string          `json:"course_id"`
}

type User struct {
	ID      uint64         `gorm:"primaryKey"`
	Reviews []ReviewDetail `json:"-" gorm:"foreignKey:OwnerID"`
}

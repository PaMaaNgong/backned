package main

type ErrCourseNotFound struct {
	ID string
}

func (e ErrCourseNotFound) Error() string {
	return "course " + e.ID + " not found"
}

type ErrNotEnum struct {
	Value any
}

func (e ErrNotEnum) Error() string {
	return "value is not enum"
}

type ErrInvalidAccessToken struct{}

func (e ErrInvalidAccessToken) Error() string {
	return "invalid access token"
}

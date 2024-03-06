package main

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func newCors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowOriginFunc = func(origin string) bool { return true }
	return cors.New(config)
}

func NewRouter(repo Repository, auth Auth) *gin.Engine {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		RegisterEnumValidator(v)
	}

	r := gin.Default()
	r.Use(newCors())

	r.GET("/v2/courses", func(c *gin.Context) {
		query := c.Query("query")
		courses, err := repo.GetCourses(query)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, courses)
	})

	r.GET("/v2/course/:id", func(c *gin.Context) {
		id := c.Param("id")
		course, err := repo.GetCourseDetail(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, course)
	})

	r.GET("/v2/course/:id/grades", func(c *gin.Context) {
		id := c.Param("id")
		grades, err := repo.GetCourseGrades(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, grades)
	})

	r.POST("/v2/course/:id", func(c *gin.Context) {
		if userId, err := auth.Verify(c.GetHeader("accessToken")); userId != 1 || err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		id := c.Param("id")
		var course CourseDetail
		err := c.BindJSON(&course)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		course.ID = id
		repo.AddCourse(course)
	})

	r.GET("/v2/course/:id/reviews", func(c *gin.Context) {
		id := c.Param("id")
		reviews, err := repo.GetReviewsOverview(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, reviews)
	})

	r.POST("/v2/course/:id/reviews", func(c *gin.Context) {
		id := c.Param("id")
		userId, err := auth.Verify(c.GetHeader("accessToken"))
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		var review ReviewDetail
		err = c.BindJSON(&review)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if review.Rating < 1 || review.Rating > 5 {
			c.Status(http.StatusBadRequest)
			return
		}
		err = repo.CreateReview(userId, id, review)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusCreated)
	})

	r.PATCH("/v2/course/:course_id/reviews/:review_id", func(c *gin.Context) {
		courseId := c.Param("course_id")
		userId, err := auth.Verify(c.GetHeader("accessToken"))
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		var review ReviewDetail
		err = c.BindJSON(&review)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if review.Rating < 1 || review.Rating > 5 {
			c.Status(http.StatusBadRequest)
			return
		}
		if reviewId, err := strconv.ParseUint(c.Param("review_id"), 10, 64); err != nil {
			err = repo.EditReview(userId, courseId, reviewId, review)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
		}
		c.Status(http.StatusOK)
	})

	r.DELETE("/v2/course/:course_id/reviews/:review_id", func(c *gin.Context) {
		courseId := c.Param("course_id")
		userId, err := auth.Verify(c.GetHeader("accessToken"))
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		if reviewId, err := strconv.ParseUint(c.Param("review_id"), 10, 64); err != nil {
			err = repo.DeleteReview(userId, courseId, reviewId)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
		}
		c.Status(http.StatusOK)
	})

	r.GET("/v2/course/:id/reviews/detail", func(c *gin.Context) {
		id := c.Param("id")
		reviews, err := repo.GetReviewsDetail(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, reviews)
	})

	return r
}

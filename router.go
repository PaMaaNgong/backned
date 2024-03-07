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
	config.AllowHeaders = []string{"accessToken"}
	return cors.New(config)
}

func NewRouter(repo Repository, auth Auth) *gin.Engine {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		RegisterEnumValidator(v)
	}

	r := gin.Default()
	r.Use(newCors())

	r.GET("/v1/courses", func(c *gin.Context) {
		query := c.Query("query")
		courses, err := repo.GetCourses(query)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, courses)
	})

	r.GET("/v1/course/:id", func(c *gin.Context) {
		id := c.Param("id")
		course, err := repo.GetCourseDetail(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, course)
	})

	r.GET("/v1/course/:id/grades", func(c *gin.Context) {
		id := c.Param("id")
		grades, err := repo.GetCourseGrades(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, grades)
	})

	r.POST("/v1/course/:id", func(c *gin.Context) {
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

	r.GET("/v1/course/:id/reviews", func(c *gin.Context) {
		id := c.Param("id")
		reviews, err := repo.GetReviewsOverview(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, reviews)
	})

	r.POST("/v1/course/:id/reviews", func(c *gin.Context) {
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

	r.PATCH("/v1/course/:course_id/reviews/:review_id", func(c *gin.Context) {
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
		if reviewId, err := strconv.ParseUint(c.Param("review_id"), 10, 64); err == nil {
			err = repo.EditReview(userId, courseId, reviewId, review)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
		} else {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	r.DELETE("/v1/course/:course_id/reviews/:review_id", func(c *gin.Context) {
		courseId := c.Param("course_id")
		userId, err := auth.Verify(c.GetHeader("accessToken"))
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		if reviewId, err := strconv.ParseUint(c.Param("review_id"), 10, 64); err == nil {
			err = repo.DeleteReview(userId, courseId, reviewId)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
		} else {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusOK)
	})

	r.GET("/v1/course/:id/reviews/detail", func(c *gin.Context) {
		id := c.Param("id")
		reviews, err := repo.GetReviewsDetail(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, reviews)
	})

	r.GET("/v1/auth", func(c *gin.Context) {
		c.Header("Location", "/v1/callback")
		c.Status(http.StatusMovedPermanently)
	})

	r.GET("/v1/callback", func(c *gin.Context) {
		c.Header("Location", "https://rococo-blini-2e875c.netlify.app/?accessToken=token-1")
		c.Status(http.StatusMovedPermanently)
	})

	r.GET("/v1/profile/reviews/courses", func(c *gin.Context) {
		userId, err := auth.Verify(c.GetHeader("accessToken"))
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		courses, _ := repo.GetCoursesByUser(userId)
		c.JSON(http.StatusOK, courses)
	})

	r.GET("/v1/profile/reviews/:course_id", func(c *gin.Context) {
		courseId := c.Param("course_id")
		userId, err := auth.Verify(c.GetHeader("accessToken"))
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		review, err := repo.GetReviewByUser(userId, courseId)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, review)
	})

	return r
}

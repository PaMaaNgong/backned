package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

func newCors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowOriginFunc = func(origin string) bool { return true }
	return cors.New(config)
}

func NewRouter(repo Repository) *gin.Engine {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		RegisterEnumValidator(v)
	}

	r := gin.Default()
	r.Use(newCors())

	r.GET("/v1/courses", func(c *gin.Context) {
		query := c.Query("query")
		limit, offset, err := getLimitAndOffset(c, 10, 0)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		courses, err := repo.GetCourses(query, limit, offset)
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

	r.GET("/v1/course/:id/reviews", func(c *gin.Context) {
		id := c.Param("id")
		limit, offset, err := getLimitAndOffset(c, 10, 0)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		reviews, err := repo.GetReviewsOverview(id, limit, offset)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, reviews)
	})

	r.POST("/v1/course/:id/reviews", func(c *gin.Context) {
		id := c.Param("id")
		var review ReviewDetail
		err := c.BindJSON(&review)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if review.Rating < 1 || review.Rating > 5 {
			c.Status(http.StatusBadRequest)
			return
		}
		err = repo.CreateReview(id, review)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusCreated)
	})

	r.GET("/v1/course/:id/reviews/detail", func(c *gin.Context) {
		id := c.Param("id")
		limit, offset, err := getLimitAndOffset(c, 10, 0)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		reviews, err := repo.GetReviewsDetail(id, limit, offset)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, reviews)
	})

	return r
}

func getLimitAndOffset(c *gin.Context, defaultLimit int, defaultOffset int) (int, int, error) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if c.Query("limit") != "" && err != nil {
		return 0, 0, err
	} else if c.Query("limit") == "" {
		limit = defaultLimit
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if c.Query("offset") != "" && err != nil {
		return 0, 0, err
	} else if c.Query("offset") == "" {
		offset = defaultOffset
	}
	return limit, offset, nil
}

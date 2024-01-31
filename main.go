package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func newCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	})
}

func main() {
	s := NewMySQLRepository("db-mysql-sgp1-18121-do-user-15505600-0.c.db.ondigitalocean.com","kritanon","AVNS_JyoxklF56B55Q1Y1MjU","25060","course-review-development")
	r := NewRouter(s)
	r.Use(newCors())
	_ = r.Run()
}

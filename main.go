package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	godotenv.Load()
	s := NewMySQLRepository(os.Getenv("MYSQL_HOST"),os.Getenv("MYSQL_USER"),os.Getenv("MYSQL_PASSWORD"),os.Getenv("MYSQL_PORT"),os.Getenv("MYSQL_DATABASE"))
	r := NewRouter(s)
	r.Use(newCors())
	_ = r.Run()
}

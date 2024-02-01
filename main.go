package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	s := NewMySQLRepository(os.Getenv("MYSQL_HOST"),os.Getenv("MYSQL_USER"),os.Getenv("MYSQL_PASSWORD"),os.Getenv("MYSQL_PORT"),os.Getenv("MYSQL_DATABASE"))
	r := NewRouter(s)
	_ = r.Run()
}

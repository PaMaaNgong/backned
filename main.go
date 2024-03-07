package main

import (
	"github.com/joho/godotenv"
	"os"
)

func main() {
	_ = godotenv.Load()
	s := NewMySQLRepository(
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	a := NewStubAuth()
	r := NewRouter(s, a, os.Getenv("OAUTH_REDIRECT_URI"))
	_ = r.Run()
}

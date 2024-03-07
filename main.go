package main

import (
	"os"

	"github.com/joho/godotenv"
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
	r := NewRouter(s, a)
	_ = r.Run()
}

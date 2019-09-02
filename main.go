package main

import "os"

func main() {
	os.Setenv("DB_USERNAME", "ashishrao")
	os.Setenv("DB_NAME", "meal_api")
	os.Setenv("DB_PASSWORD", "valeyforge")
	os.Setenv("DB_SSL_MODE", "disable")
	a := App{}
	a.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"))

	a.Run()
}

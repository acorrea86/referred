package main

import (
	"blumer-ms-refers/di"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	handler, err := di.Initialize()
	if err != nil {
		panic("fatal err: " + err.Error())
	}

	handler.Handler()
}

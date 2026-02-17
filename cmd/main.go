package main

import (
	"log"

	"WEB/internal/api"
)

func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application end!")
}

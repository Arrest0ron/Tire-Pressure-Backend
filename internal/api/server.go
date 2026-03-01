package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"web_backend/internal/app/handler"
	"web_backend/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализация репозитория")
	}

	handler1 := handler.NewHandler(repo)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/tires")
	})
	r.GET("/tires", handler1.GetTires)
	r.GET("/tire/:id", handler1.GetTire)
	r.GET("/tire_pressure/:id", handler1.GetTirePressure)

	r.Run("127.0.0.1:8081")
	log.Println("Server down")
}

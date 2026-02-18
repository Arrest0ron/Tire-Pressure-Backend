package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"web_backend/internal/app/h
	"web_backend/internal/app/handler"
	"web_backend/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализация репозитория")
	}

	handler := handler.NewHandler(repo)


	
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/tires")
	})
	r.GET("/tires", handler.GetTires)
	r.GET("/tire/:id", handler.GetTire)
	r.GET("/calculation/:id", handler.GetCalculation)

	r.Run()
	log.Println("Server down")
}

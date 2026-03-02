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
		logrus.Fatal("Ошибка инициализации репозитория: ", err)
	}

	h := handler.NewHandler(repo)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/tires")
	})
	
	// ✅ Маршруты: URL с подчёркиванием (как в ТЗ), методы — CamelCase (как в Go)
	r.GET("/tires", h.GetTires)                    // список шин
	r.GET("/tire/:id", h.GetTire)                  // деталь шины
	r.GET("/tire_pressure/:id", h.GetTirePressure) // заявка Tire_pressure

	r.Run("127.0.0.1:8081")
	log.Println("Server down")
}
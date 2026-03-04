package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"web_backend/internal/app/config"
	"web_backend/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
	Config     *config.Config
}

func NewHandler(r *repository.Repository, cfg *config.Config) *Handler {
	return &Handler{
		Repository: r,
		Config:     cfg,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	// Новые маршруты для шин
	router.GET("/", h.GetTires)
	router.GET("/tires", h.GetTires)
	router.GET("/tire/:id", h.GetTire)
	router.GET("/tire_pressure/:id", h.GetTirePressure)

	// Маршруты для управления заявками на расчёт давления
	router.POST("/tire_pressure/add", h.AddTireToPressure)
	router.POST("/tire_pressure/update_coeff", h.UpdateCoatingCoeff)
	router.POST("/tire_pressure/delete", h.DeleteTireEntry)

	// Маршруты для добавления услуг в заявку
	router.POST("/tire_pressure/add_tire", h.AddTireToPressure)

	// Новые маршруты для управления заявками
	router.POST("/tire_pressure/delete_request", h.DeleteTirePressure)
	router.POST("/tire_pressure/complete", h.CompleteTirePressure)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

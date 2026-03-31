package handler

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"metoda/internal/app/repository"
)

// ✅ Константа доступна всем файлам пакета handler
const minioBaseURL = "http://localhost:9090/tire-bucket"

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) getMinioURL() string {
	if url := os.Getenv("MINIO_URL"); url != "" {
		return url
	}
	return minioBaseURL
}

// ─── Register Routes ─────────────────────────────────────────────────────────

func (h *Handler) RegisterHandler(router *gin.Engine) {
	// 3 HTML страницы
	router.GET("/tires", h.Index)
	router.GET("/tire/:id", h.TirePage)
	router.GET("/tire-pressure/:id", h.TirePressurePage)

	// Tires (3 метода)
	router.GET("/api/tires", h.GetTires)
	router.GET("/api/tires/:id", h.GetTire)
	router.POST("/api/tires", h.CreateTire)

	// Tire Pressures (7 методов)
	router.GET("/api/tire-pressures/cart", h.GetTirePressureCart)
	router.GET("/api/tire-pressures", h.GetTirePressures)
	router.GET("/api/tire-pressures/:id", h.GetTirePressure)
	router.PUT("/api/tire-pressures/:id", h.UpdateTirePressure)
	router.PUT("/api/tire-pressures/:id/form", h.FormTirePressure)
	router.PUT("/api/tire-pressures/:id/finish", h.FinishTirePressure)
	router.DELETE("/api/tire-pressures/:id", h.DeleteTirePressure)

	// Tire Pressure Entries (3 метода)
	router.POST("/api/tire-pressure-entries/add/:tire_id", h.AddToTirePressure)
	router.DELETE("/api/tire-pressure-entries/:tire_id/:tire_pressure_id", h.DeleteTireFromPressure)
	router.PUT("/api/tire-pressure-entries/:tire_id/:tire_pressure_id", h.UpdateTirePressureEntry)

	// Users (3 метода)
	router.POST("/api/users/signup", h.APISignUp)
	router.POST("/api/users/signin", h.APISignIn)
	router.POST("/api/users/signout", h.APISignOut)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найдено"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещён"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}

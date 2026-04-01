package handler

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"metoda/internal/app/repository"
)

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
	// ─── HTML страницы (публичные) ──────────────────────────────────────
	router.GET("/tires", h.Index)
	router.GET("/tire/:id", h.TirePage)
	router.GET("/tire-pressure/:id", h.TirePressurePage)

	// ─── API группа ─────────────────────────────────────────────────────
	api := router.Group("/api")

	// === ПУБЛИЧНЫЕ эндпоинты (без авторизации) ===
	// Чтение каталога + регистрация и вход
	api.GET("/tires", h.GetTires)
	api.GET("/tires/:id", h.GetTire)
	api.POST("/users/signup", h.APISignUp)
	api.POST("/users/signin", h.APISignIn)

	// === ЗАЩИЩЁННЫЕ эндпоинты (требуется авторизация) ===
	needAuth := api.Group("")
	needAuth.Use(h.AuthMiddleware())

	// Выход из системы
	needAuth.POST("/users/signout", h.APISignOut)

	// Tire Pressure (черновики/заявки)
	needAuth.GET("/tire-pressures/cart", h.GetTirePressureCart)
	needAuth.GET("/tire-pressures", h.GetTirePressures)
	needAuth.GET("/tire-pressures/:id", h.GetTirePressure)
	needAuth.PUT("/tire-pressures/:id", h.UpdateTirePressure)
	needAuth.PUT("/tire-pressures/:id/form", h.FormTirePressure)
	needAuth.DELETE("/tire-pressures/:id", h.DeleteTirePressure)

	// Tire Pressure Entries (записи в заявке)
	needAuth.POST("/tire-pressure-entries/add/:tire_id", h.AddToTirePressure)
	needAuth.DELETE("/tire-pressure-entries/:tire_id/:tire_pressure_id", h.DeleteTireFromPressure)
	needAuth.PUT("/tire-pressure-entries/:tire_id/:tire_pressure_id", h.UpdateTirePressureEntry)

	// === ЭНДПОИНТЫ ДЛЯ МОДЕРАТОРОВ ===
	mod := api.Group("")
	mod.Use(h.AuthMiddleware())
	mod.Use(h.RequireModerator())

	// Создание шины (только модератор)
	mod.POST("/tires", h.CreateTire)

	// Завершение заявки (только модератор)
	mod.PUT("/tire-pressures/:id/finish", h.FinishTirePressure)

	// Swagger
	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// ─── Error Handler (расширенный) ─────────────────────────────────────────

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
	case errors.Is(err, repository.ErrUnauthorized):
		errorMessage = "Необходимо войти в систему"
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "description": errorMessage})
		return
	case errors.Is(err, repository.ErrForbidden):
		errorMessage = "Недостаточно прав"
		ctx.JSON(http.StatusForbidden, gin.H{"status": "error", "description": errorMessage})
		return
	case errors.Is(err, repository.ErrInvalidToken):
		errorMessage = "Неверный токен"
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "description": errorMessage})
		return
	case errors.Is(err, repository.ErrTokenExpired):
		errorMessage = "Срок токена истёк"
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "description": errorMessage})
		return
	default:
		errorMessage = err.Error()
	}

	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}
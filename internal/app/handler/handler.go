package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"metoda/internal/app/repository"
)

const minioBaseURL = "http://localhost:9090/tires"

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.GetTires)
	router.GET("/tire/:id", h.GetTire)
	// ✅ Убран отдельный роут для черновика — только /:id (как в оригинале)
	router.GET("/tire-pressure/:id", h.GetTirePressureByID)
	router.POST("/tire-pressure/update-item", h.UpdateTirePressureItem)
	router.POST("/form-tire-pressure", h.FormTirePressure)
	router.POST("/add-to-tire-pressure", h.AddToTirePressure)
	router.POST("/delete-tire-pressure", h.DeleteTirePressure)

	// домен Tire (услуги)
	router.GET("/api/tires", h.APIGetTires)
	router.GET("/api/tires/:id", h.APIGetTire)
	router.POST("/api/tires", h.APICreateTire)

	// домен М-М (tire_pressure_entries)
	router.POST("/api/tires/:id/add-to-tire-pressure", h.APIAddToCart)
	router.DELETE("/api/tire-pressure-entries/:tire_id/:tire_pressure_id", h.APIDeleteFromCart)
	router.PUT("/api/tire-pressure-entries/:tire_id/:tire_pressure_id", h.APIUpdateCartItem)

	// домен TirePressure (заявки)
	router.GET("/api/tire-pressures/cart", h.APIGetCart)
	router.GET("/api/tire-pressures", h.APIGetTirePressures)
	router.GET("/api/tire-pressures/:id", h.APIGetTirePressure)
	router.PUT("/api/tire-pressures/:id", h.APIUpdateTirePressure)
	router.PUT("/api/tire-pressures/:id/form", h.APIFormTirePressure)
	router.PUT("/api/tire-pressures/:id/finish", h.APIFinishTirePressure)
	router.DELETE("/api/tire-pressures/:id", h.APIDeleteTirePressure)

	// домен Users
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
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
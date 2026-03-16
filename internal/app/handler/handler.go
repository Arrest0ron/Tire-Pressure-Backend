package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"metoda/internal/app/repository"
)

const minioBaseURL = "http://localhost:9090/constructions"

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.GetConstructions)
	router.GET("/construction/:id", h.GetConstruction)
	router.GET("/dendrochronology", h.GetDendrochronology)
	router.POST("/dendrochronology/update-item", h.UpdateDendrochronologyItem)
	router.POST("/form-dendrochronology", h.FormDendrochronology)
	router.POST("/add-to-dendrochronology", h.AddToDendrochronology)
	router.POST("/delete-dendrochronology", h.DeleteDendrochronology)

	// домен Construction (услуги)
	router.GET("/api/constructions", h.APIGetConstructions)
	router.GET("/api/constructions/:id", h.APIGetConstruction)
	router.POST("/api/constructions", h.APICreateConstruction)

	// домен М-М (dendrochronology_constructions)
	router.POST("/api/constructions/:id/add-to-dendrochronology", h.APIAddToCart)
	router.DELETE("/api/dendrochronology-constructions/:construction_id/:dendrochronology_id", h.APIDeleteFromCart)
	router.PUT("/api/dendrochronology-constructions/:construction_id/:dendrochronology_id", h.APIUpdateCartItem)

	// домен Dendrochronology (заявки)
	router.GET("/api/dendrochronologies/cart", h.APIGetCart)
	router.GET("/api/dendrochronologies", h.APIGetDendrochronologies)
	router.GET("/api/dendrochronologies/:id", h.APIGetDendrochronology)
	router.PUT("/api/dendrochronologies/:id", h.APIUpdateDendrochronology)
	router.PUT("/api/dendrochronologies/:id/form", h.APIFormDendrochronology)
	router.PUT("/api/dendrochronologies/:id/finish", h.APIFinishDendrochronology)
	router.DELETE("/api/dendrochronologies/:id", h.APIDeleteDendrochronology)

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
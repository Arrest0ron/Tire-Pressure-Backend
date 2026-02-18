package handler

import (
	"web_backend/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetTires(ctx *gin.Context) {
	var tires []repository.Tire
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		tires, err = h.Repository.GetTires()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		tires, err = h.Repository.GetTireByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	cartTires, err := h.Repository.GetRequest()
	cartCount := 0
	if err == nil {
		cartCount = len(cartTires)
	}

	ctx.HTML(http.StatusOK, "tires.html", gin.H{
		"tires":     tires,
		"query":     searchQuery,
		"cartCount": cartCount,
	})
}

func (h *Handler) GetTire(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	tire, err := h.Repository.GetTire(id)
	if err != nil {
		logrus.Error(err)
	}

	// Get calculation parameters from query or use defaults
	temperature := ctx.DefaultQuery("temperature", "20")
	weight := ctx.DefaultQuery("weight", "1500")
	surface := ctx.DefaultQuery("surface", "1.0")

	tempVal, _ := strconv.ParseFloat(temperature, 64)
	weightVal, _ := strconv.ParseFloat(weight, 64)
	surfaceVal, _ := strconv.ParseFloat(surface, 64)

	// Calculate recommended pressure
	recommendedPressure := repository.CalculatePressure(tire.BasePressure, tempVal, weightVal, surfaceVal)

	ctx.HTML(http.StatusOK, "tire.html", gin.H{
		"tire":               tire,
		"temperature":        tempVal,
		"weight":             weightVal,
		"surface":            surfaceVal,
		"recommendedPressure": recommendedPressure,
	})
}

func (h *Handler) GetCalculation(ctx *gin.Context) {
	var tires []repository.Tire
	var err error

	tires, err = h.Repository.GetRequest()
	if err != nil {
		logrus.Error(err)
	}

	// Get form values for pressure calculation
	temperature := ctx.PostForm("temperature")
	weight := ctx.PostForm("weight")
	surface := ctx.PostForm("surface")

	tempVal := 20.0
	weightVal := 1500.0
	surfaceVal := 1.0

	if temperature != "" {
		if t, err := strconv.ParseFloat(temperature, 64); err == nil {
			tempVal = t
		}
	}
	if weight != "" {
		if w, err := strconv.ParseFloat(weight, 64); err == nil {
			weightVal = w
		}
	}
	if surface != "" {
		if s, err := strconv.ParseFloat(surface, 64); err == nil {
			surfaceVal = s
		}
	}

	// Calculate pressure for each tire in request
	type TireWithPressure struct {
		Tire       repository.Tire
		Pressure   float64
	}
	tiresWithPressure := make([]TireWithPressure, len(tires))
	for i, tire := range tires {
		tiresWithPressure[i] = TireWithPressure{
			Tire:     tire,
			Pressure: repository.CalculatePressure(tire.BasePressure, tempVal, weightVal, surfaceVal),
		}
	}

	ctx.HTML(http.StatusOK, "calculation.html", gin.H{
		"request_tires":      tiresWithPressure,
		"temperature":       tempVal,
		"weight":            weightVal,
		"surface":           surfaceVal,
		"recommendedPressure": repository.CalculatePressure(2.2, tempVal, weightVal, surfaceVal),
	})
}

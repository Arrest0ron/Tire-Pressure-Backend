package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"web_backend/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// GetTires - GET /tires - список шин с поиском
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

	// Подсчёт количества услуг в заявке
	request, err := h.Repository.GetRequestByID(1)
	cartCount := 0
	if err == nil {
		cartCount = len(request.TireIDs)
	}

	ctx.HTML(http.StatusOK, "tires.html", gin.H{
		"tires":     tires,
		"query":     searchQuery,
		"cartCount": cartCount,
	})
}

// GetTire - GET /tire/:id - детали шины
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

	temperature := ctx.DefaultQuery("temperature", "20")
	weight := ctx.DefaultQuery("weight", "1500")
	surface := ctx.DefaultQuery("surface", "1.0")

	tempVal, _ := strconv.ParseFloat(temperature, 64)
	weightVal, _ := strconv.ParseFloat(weight, 64)
	surfaceVal, _ := strconv.ParseFloat(surface, 64)

	recommendedPressure := repository.CalculatePressure(tire.TireCoefficient, tempVal, weightVal, surfaceVal)

	ctx.HTML(http.StatusOK, "tire.html", gin.H{
		"tire":                tire,
		"temperature":         tempVal,
		"weight":              weightVal,
		"surface":             surfaceVal,
		"recommendedPressure": recommendedPressure,
	})
}

// GetTirePressure - GET /tire_pressure/:id - заявка по ID
func (h *Handler) GetTirePressure(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	// Получаем заявку по ID из словаря
	request, err := h.Repository.GetRequestByID(id)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "Заявка не найдена")
		return
	}

	// Получаем шины для этой заявки
	tires, err := h.Repository.GetRequestTires(id)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "Заявка пуста")
		return
	}

	// Логируем количество шин для отладки
	logrus.Infof("Заявка %d: найдено %d шин", id, len(tires))
	for i, tire := range tires {
		logrus.Infof("Шина %d: ID=%d, Title=%s", i+1, tire.ID, tire.Title)
	}

	// Вычисляем давление для каждой шины
	type TireWithPressure struct {
		Tire     repository.Tire
		Pressure float64
	}
	tiresWithPressure := make([]TireWithPressure, len(tires))
	for i, tire := range tires {
		tiresWithPressure[i] = TireWithPressure{
			Tire:     tire,
			Pressure: repository.CalculatePressure(tire.TireCoefficient, request.AirTemperature, request.CarWeight, request.SurfaceCoefficient),
		}
	}

	ctx.HTML(http.StatusOK, "tire_pressure.html", gin.H{
		"request":             request,
		"request_tires":       tiresWithPressure,
		"temperature":         request.AirTemperature,
		"weight":              request.CarWeight,
		"surface":             request.SurfaceCoefficient,
		"recommendedPressure": request.PressureResult,
	})
}

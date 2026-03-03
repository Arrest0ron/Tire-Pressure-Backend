package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"web_backend/internal/app/repository"
)

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repo: r}
}

// calculatePressure — локальная функция расчёта давления (так как в repository она удалена).
func calculatePressure(tireCoeff, coatingCoeff float64, airTemp, carWeight float64) float64 {
	if coatingCoeff <= 0 {
		coatingCoeff = 1
	}
	tempFactor := (20.0 - airTemp) * 0.05
	weightFactor := (carWeight - 1500) * 0.001
	return (tempFactor + weightFactor + 1) * tireCoeff * coatingCoeff * 10
}

// getEntryCount — вспомогательный метод: возвращает количество записей в заявке Tire_pressure
func (h *Handler) getEntryCount() int {
	req, err := h.Repo.GetTirePressure(1)
	if err != nil || len(req.Entries) == 0 {
		return 0
	}
	return req.EntryCount // теперь работает, поле добавлено в repository
}

// GET /tires — список шин
func (h *Handler) GetTires(ctx *gin.Context) {
	query := ctx.Query("query")
	tires, err := h.Repo.GetTires()
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusInternalServerError, "Error")
		return
	}
	if query != "" {
		tires, _ = h.Repo.GetTiresByTitle(query)
	}
	ctx.HTML(http.StatusOK, "tires.html", gin.H{
		"tires":      tires,
		"query":      query,
		"entryCount": h.getEntryCount(),
	})
}

// GET /tire/:id — деталь шины
func (h *Handler) GetTire(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	tire, err := h.Repo.GetTire(id)
	if err != nil {
		ctx.String(http.StatusNotFound, "Not found")
		return
	}
	temp, _ := strconv.ParseFloat(ctx.DefaultQuery("temperature", "20"), 64)
	weight, _ := strconv.ParseFloat(ctx.DefaultQuery("weight", "1500"), 64)
	coating, _ := strconv.ParseFloat(ctx.DefaultQuery("coating", "1.0"), 64)

	ctx.HTML(http.StatusOK, "tire.html", gin.H{
		"tire":                tire,
		"temperature":         temp,
		"weight":              weight,
		"coating":             coating,
		"recommendedPressure": calculatePressure(tire.TireCoefficient, coating, temp, weight), // [ИСПРАВЛЕНО] локальная функция
		"entryCount":          h.getEntryCount(),
	})
}

// GET /tire_pressure/:id — заявка Tire_pressure
func (h *Handler) GetTirePressure(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	req, err := h.Repo.GetTirePressure(id)
	if err != nil || len(req.Entries) == 0 {
		ctx.String(http.StatusNotFound, "Not found")
		return
	}

	// [ИСПРАВЛЕНО] Гарантированно создаём мапу, даже если GetTires вернёт ошибку
	tires, err := h.Repo.GetTires()
	if err != nil {
		logrus.Warn("GetTires error in tire_pressure handler:", err)
	}
	tiresMap := make(map[int]repository.Tire) // всегда создаём, никогда не nil
	for _, t := range tires {
		tiresMap[t.ID] = t
	}

	// [ОТЛАДКА] Лог для проверки, что мапа заполнена
	logrus.Debugf("tire_pressure: passed %d tires to template", len(tiresMap))

	ctx.HTML(http.StatusOK, "tire_pressure.html", gin.H{
		"request":             req,
		"request_entries":     req.Entries,
		"tires":               tiresMap, // гарантированно не nil
		"temperature":         req.AirTemperature,
		"weight":              req.CarWeight,
		"coating":             req.Entries[0].CoatingCoeff,
		"recommendedPressure": req.Entries[0].Pressure,
		"entryCount":          req.EntryCount,
	})
}

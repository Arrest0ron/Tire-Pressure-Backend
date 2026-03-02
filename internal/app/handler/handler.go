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

// getEntryCount — вспомогательный метод: возвращает количество записей в заявке Tire_pressure
func (h *Handler) getEntryCount() int {
	req, err := h.Repo.GetTirePressure(1)
	if err != nil || len(req.Entries) == 0 {
		return 0
	}
	return
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
		"tires": tires,
		"query": query,
		"tires":      tires,
		"query":      query,
		"entryCount": h.getEntryCount(), // ✅ Количество шин в заявке

// GET /tire/:id — деталь шины + расчёт давления
func (h *Handler) GetTire(ctx *gin.Context) {
// GET /tire/:id — деталь шины
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
		"recommendedPressure": repository.CalculatePressure(tire.TireCoefficient, coating, temp, weight),
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
	ctx.HTML(http.StatusOK, "tire_pressure.html", gin.H{
		"request":             req,
		"request_entries":     req.Entries,
		"temperature":         req.Entries[0].AirTemperature,
		"weight":              req.Entries[0].CarWeight,
		"coating":             req.Entries[0].CoatingCoeff,
		"recommendedPressure": req.Entries[0].Pressure,
	})
}
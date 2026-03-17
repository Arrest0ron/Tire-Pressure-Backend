package handler

import (
	"metoda/internal/app/ds"
	"metoda/internal/app/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ✅ GET /tire-pressure — редирект на черновик по ID
func (h *Handler) GetTirePressure(ctx *gin.Context) {
	draft, err := h.Repository.GetDraftTirePressure(uint(repository.GetUserID()))
	if err != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	ctx.Redirect(http.StatusFound, "/tire-pressure/"+strconv.Itoa(int(draft.TirePressureID)))
}

// ✅ GET /tire-pressure/:id — просмотр заявки по ID
func (h *Handler) GetTirePressureByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	t, err := h.Repository.GetTirePressureByID(id)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	if t.Status == ds.StatusDraft && int(t.CreatorID) != repository.GetUserID() {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	if t.Status == ds.StatusDeleted {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	_, views, err := h.Repository.GetTirePressureWithEntries(t.TirePressureID)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.HTML(http.StatusOK, "tirepressurepage.html", gin.H{
		"tirePressure": t,
		"tires":        views,
		"minioBase":    minioBaseURL,
		"readonly":     t.Status != ds.StatusDraft,
	})
}

// ✅ Добавление шины → редирект на главную
func (h *Handler) AddToTirePressure(ctx *gin.Context) {
	strId := ctx.PostForm("tire_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Repository.AddTireToTirePressure(uint(id), uint(repository.GetUserID()))
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ НА ГЛАВНУЮ!
	ctx.Redirect(http.StatusFound, "/")
}

// ✅ Обновление коэффициента → редирект на заявку по ID
func (h *Handler) UpdateTirePressureItem(ctx *gin.Context) {
	strId := ctx.PostForm("item_id")
	itemID, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coatingStr := ctx.PostForm("coating_coefficient")
	coatingCoeff := 1.0
	if coatingStr != "" {
		if v, convErr := strconv.ParseFloat(coatingStr, 64); convErr == nil && v > 0 {
			coatingCoeff = v
		}
	}

	err = h.Repository.UpdateTirePressureItem(uint(itemID), coatingCoeff)
	if err != nil {
		logrus.Error(err)
	}

	// ✅ Получаем tire_pressure_id из записи м-м
	item, err := h.Repository.GetTirePressureEntryByID(uint(itemID))
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/tire-pressure/"+strconv.Itoa(int(item.TirePressureID)))
}

// ✅ Формирование → редирект на заявку по ID
func (h *Handler) FormTirePressure(ctx *gin.Context) {
	strId := ctx.PostForm("tire_pressure_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	err = h.Repository.FormTirePressure(uint(id))
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.Redirect(http.StatusFound, "/tire-pressure/"+strconv.Itoa(id))
}

// ✅ Удаление → редирект на главную
func (h *Handler) DeleteTirePressure(ctx *gin.Context) {
	strId := ctx.PostForm("tire_pressure_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	err = h.Repository.DeleteTirePressureBySQL(uint(id))
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

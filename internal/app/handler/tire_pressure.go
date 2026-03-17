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

// ✅ Удалена функция GetTirePressure — теперь только GetTirePressureByID (как в оригинале)

// ✅ GET /tire-pressure/:id — просмотр заявки по ID (черновик/сформирован/завершён)
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

	// ✅ Проверка доступа:
	// 1. Удалённые — никому
	// 2. Черновики — только создателю
	// 3. Сформированные/завершённые/отклонённые — создателю и модераторам
	if t.Status == ds.StatusDeleted {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	if t.Status == ds.StatusDraft && int(t.CreatorID) != repository.GetUserID() {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	// Проверка для модератора (может видеть все не-удалённые)
	moderator, _ := h.Repository.GetUserByID(repository.GetUserID())
	isModerator := moderator.IsModerator

	if t.Status != ds.StatusDraft && int(t.CreatorID) != repository.GetUserID() && !isModerator {
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
		"readonly":     t.Status != ds.StatusDraft, // ✅ Блокировка полей для не-черновиков
	})
}

func (h *Handler) AddToTirePressure(ctx *gin.Context) {
	strId := ctx.PostForm("tire_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ SINGLETON: замена creatorID → repository.GetUserID()
	err = h.Repository.AddTireToTirePressure(uint(id), uint(repository.GetUserID()))
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

// UpdateTirePressureItem обрабатывает обновление одной строки корзины:
// коэффициент типа покрытия для шины.
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
	// Редирект с 303 заставляет браузер сделать новый GET, без подстановки кэша
	ctx.Redirect(http.StatusSeeOther, "/tire-pressure/"+ctx.PostForm("item_id"))
}

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

	// ✅ Редирект на страницу просмотра по ID (как в оригинале)
	ctx.Redirect(http.StatusFound, "/tire-pressure/"+strconv.Itoa(id))
}

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

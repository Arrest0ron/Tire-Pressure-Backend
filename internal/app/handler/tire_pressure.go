package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"metoda/internal/app/ds"
	"metoda/internal/app/repository"
	"metoda/internal/app/serializer"
)

// ─── HTML Pages ────────────────────────────────────────────────────────────

// TirePressurePage Страница заявки (HTML)
// @Summary Страница заявки на проверку давления
// @Tags tire-pressure-html
// @Produce html
// @Param id path int true "ID заявки"
// @Success 200 {string} string "HTML-страница заявки"
// @Failure 302 "Редирект при ошибке или отсутствии доступа"
// @Router /tire-pressure/{id} [get]
func (h *Handler) TirePressurePage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/tires")
		return
	}

	t, err := h.Repository.GetTirePressureByID(id)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/tires")
		return
	}

	if t.Status == ds.StatusDeleted {
		ctx.Redirect(http.StatusFound, "/tires")
		return
	}

	// Проверка прав: черновик видит только создатель
	uid, err := authUserIDUint(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/tires")
		return
	}
	if t.Status == ds.StatusDraft && t.CreatorID != uint(uid) {
		ctx.Redirect(http.StatusFound, "/tires")
		return
	}

	items, err := h.Repository.GetTirePressureEntriesAPI(t.TirePressureID)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/tires")
		return
	}

	ctx.HTML(http.StatusOK, "tirepressurepage.html", gin.H{
		"tire_pressure_id": t.TirePressureID,
		"tires":            items,
		"minioBase":        h.getMinioURL(),
		"air_temperature":  t.AirTemperature,
		"car_weight":       t.CarWeight,
		"date_created":     t.DateCreate.Format("02.01.2006 15:04"),
		"status":           t.Status,
	})
}

// ─── API: Tire Pressures ───────────────────────────────────────────────────
// GetTirePressureCart — иконка корзины (без авторизации, всегда 200)
// @Summary Иконка корзины
// @Description Возвращает ID черновика и количество шин в нём. Работает без авторизации.
// @Tags tire-pressures
// @Produce json
// @Success 200 {object} serializer.CartJSON
// @Router /api/tire_pressure/tire_pressure-cart [get]
func (h *Handler) GetTirePressureCart(ctx *gin.Context) {
	// 🔥 Логируем начало запроса
	log.Printf("🔍 [GetTirePressureCart] Запрос корзины")

	// 🔥 Логируем заголовки (для отладки авторизации)
	authHeader := ctx.GetHeader("Authorization")
	log.Printf("🔍 [GetTirePressureCart] Authorization header: '%s'", authHeader)

	// 🔥 Вызываем функцию авторизации и логируем результат
	uid, err := authUserIDUint(ctx)
	log.Printf("🔍 [GetTirePressureCart] authUserIDUint → uid=%d, err=%v", uid, err)

	var cartID uint
	var count int64

	if err == nil && uid > 0 {
		log.Printf("🔍 [GetTirePressureCart] ✅ Пользователь авторизован (uid=%d), вызываю GetCartInfo...", uid)

		cartID, count, err = h.Repository.GetCartInfo(uid)
		log.Printf("🔍 [GetTirePressureCart] GetCartInfo(%d) → cartID=%d, count=%d, err=%v", uid, cartID, count, err)
	} else {
		log.Printf("🔍 [GetTirePressureCart] ❌ Не авторизован (err=%v, uid=%d) → возвращаю гостевую корзину", err, uid)
	}

	// 🔥 Логируем итоговый ответ
	log.Printf("🔍 [GetTirePressureCart] 📤 Ответ: {tire_pressure_id:%d, tires_count:%d}", cartID, count)

	ctx.JSON(http.StatusOK, serializer.CartJSON{
		TirePressureID: cartID,
		TiresCount:     count,
	})
}

// GetTirePressures Список заявок
// @Summary Список заявок на проверку давления
// @Description Фильтры по дате и статусу. Пользователь видит свои; модератор — все.
// @Tags tire-pressures
// @Produce json
// @Param from_date query string false "Начало диапазона даты (YYYY-MM-DD)"
// @Param to_date query string false "Конец диапазона даты (YYYY-MM-DD)"
// @Param status query string false "Фильтр по статусу"
// @Security ApiKeyAuth
// @Success 200 {array} serializer.TirePressureListJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tire-pressures [get]
func (h *Handler) GetTirePressures(ctx *gin.Context) {
	var from, to time.Time
	if s := ctx.Query("from_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный формат from_date, ожидается YYYY-MM-DD"))
			return
		}
		from = t
	}
	if s := ctx.Query("to_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный формат to_date, ожидается YYYY-MM-DD"))
			return
		}
		to = t
	}
	status := ctx.Query("status")

	uid, err := authUserIDUint(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	// Передаем isModerator в репозиторий для фильтрации
	list, err := h.Repository.GetAllTirePressures(from, to, status, uid, isModeratorFromCtx(ctx))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	resp := make([]serializer.TirePressureListJSON, 0, len(list))
	for _, t := range list {
		creatorLogin := h.Repository.GetCreatorLogin(t.CreatorID)
		moderatorLogin := h.Repository.GetModeratorLogin(t.ModeratorID)
		entriesCount := h.Repository.GetTireEntriesCount(t.TirePressureID)
		resp = append(resp, serializer.TirePressureToListJSON(t, creatorLogin, moderatorLogin, entriesCount))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetTirePressure Детальная заявка
// @Summary Заявка по ID
// @Description Состав шин и логины; доступ: создатель или модератор.
// @Tags tire-pressures
// @Produce json
// @Param id path int true "ID заявки"
// @Security ApiKeyAuth
// @Success 200 {object} serializer.TirePressureDetailJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tire-pressures/{id} [get]
func (h *Handler) GetTirePressure(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	t, err := h.Repository.GetTirePressureByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	uid, err := authUserIDUint(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	// Проверка прав доступа
	if t.CreatorID != uint(uid) && !isModeratorFromCtx(ctx) {
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("%w: нет доступа к чужой заявке", repository.ErrNotAllowed))
		return
	}

	entries, err := h.Repository.GetTirePressureEntriesAPI(t.TirePressureID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	creatorLogin := h.Repository.GetCreatorLogin(t.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(t.ModeratorID)

	var dateFormed, dateCompleted *time.Time
	if t.DateFormed.Valid {
		dateFormed = &t.DateFormed.Time
	}
	if t.DateCompleted.Valid {
		dateCompleted = &t.DateCompleted.Time
	}
	var modLogin *string
	if moderatorLogin != "" {
		modLogin = &moderatorLogin
	}

	ctx.JSON(http.StatusOK, serializer.TirePressureDetailJSON{
		TirePressureID: t.TirePressureID,
		Status:         t.Status,
		DateCreate:     t.DateCreate,
		DateFormed:     dateFormed,
		DateCompleted:  dateCompleted,
		CreatorLogin:   creatorLogin,
		ModeratorLogin: modLogin,
		AirTemperature: t.AirTemperature,
		CarWeight:      t.CarWeight,
		Entries:        entries,
	})
}

// UpdateTirePressure Обновление полей заявки
// @Summary Обновить заявку (температура, вес)
// @Tags tire-pressures
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param body body serializer.TirePressureUpdateJSON true "Поля для обновления"
// @Security ApiKeyAuth
// @Success 200 {object} serializer.TirePressureListJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tire-pressures/{id} [put]
func (h *Handler) UpdateTirePressure(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	var j serializer.TirePressureUpdateJSON
	if err := ctx.ShouldBindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	uid, err := authUserIDUint(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	t, err := h.Repository.UpdateTirePressureFields(id, j, int(uid))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	creatorLogin := h.Repository.GetCreatorLogin(t.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(t.ModeratorID)
	entriesCount := h.Repository.GetTireEntriesCount(t.TirePressureID)
	ctx.JSON(http.StatusOK, serializer.TirePressureToListJSON(t, creatorLogin, moderatorLogin, entriesCount))
}

// FormTirePressure Оформить заявку (из черновика)
// @Summary Оформить заявку
// @Tags tire-pressures
// @Produce json
// @Param id path int true "ID заявки"
// @Security ApiKeyAuth
// @Success 200 {object} serializer.TirePressureListJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tire-pressures/{id}/form [put]
func (h *Handler) FormTirePressure(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	uid, err := authUserIDUint(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	t, err := h.Repository.FormTirePressureAPI(id, int(uid))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusBadRequest, err)
		}
		return
	}

	creatorLogin := h.Repository.GetCreatorLogin(t.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(t.ModeratorID)
	entriesCount := h.Repository.GetTireEntriesCount(t.TirePressureID)
	ctx.JSON(http.StatusOK, serializer.TirePressureToListJSON(t, creatorLogin, moderatorLogin, entriesCount))
}

// FinishTirePressure Завершить заявку (модератор)
// @Summary Завершить заявку
// @Description Установка итогового статуса; только для модератора.
// @Tags tire-pressures
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param body body serializer.FinishJSON true "Статус завершения"
// @Security ApiKeyAuth
// @Success 200 {object} serializer.TirePressureListJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tire-pressures/{id}/finish [put]
func (h *Handler) FinishTirePressure(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	var j serializer.FinishJSON
	if err := ctx.ShouldBindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	uid, err := authUserIDUint(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	t, err := h.Repository.FinishTirePressureAPI(id, j.Status, int(uid))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusBadRequest, err)
		}
		return
	}

	creatorLogin := h.Repository.GetCreatorLogin(t.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(t.ModeratorID)
	entriesCount := h.Repository.GetTireEntriesCount(t.TirePressureID)
	ctx.JSON(http.StatusOK, serializer.TirePressureToListJSON(t, creatorLogin, moderatorLogin, entriesCount))
}

// DeleteTirePressure Удалить заявку
// @Summary Удалить заявку
// @Tags tire-pressures
// @Produce json
// @Param id path int true "ID заявки"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "status: deleted"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tire-pressures/{id} [delete]
func (h *Handler) DeleteTirePressure(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	uid, err := authUserIDUint(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	if err := h.Repository.DeleteTirePressureAPI(id, int(uid)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
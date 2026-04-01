package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"metoda/internal/app/repository"
	"metoda/internal/app/serializer"
)

// ─── API: Tire Pressure Entries (Связь многие-ко-многим) ───────────────────

// AddToTirePressure Добавить шину в заявку (черновик)
// @Summary Добавить шину в заявку
// @Description Добавляет шину в черновик заявки; при первом добавлении создаётся заявка (201 + Location). Поддерживает JSON и form-data.
// @Tags tire-pressure-entries
// @Produce json
// @Param tire_id path int true "ID шины"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Успешно добавлено"
// @Success 201 {object} map[string]interface{} "Заявка создана"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string "Уже в заявке"
// @Failure 500 {object} map[string]string
// @Router /api/tire-pressure-entries/add/{tire_id} [post]
func (h *Handler) AddToTirePressure(ctx *gin.Context) {
	acceptHeader := ctx.GetHeader("Accept")
	isJSON := strings.Contains(acceptHeader, "application/json")

	tireIDStr := ctx.Param("tire_id")
	if tireIDStr == "" {
		tireIDStr = ctx.PostForm("tire_id")
	}
	if tireIDStr == "" {
		if isJSON {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("tire_id is required"))
			return
		}
		ctx.Redirect(http.StatusSeeOther, "/tires")
		return
	}

	tireID, err := strconv.Atoi(tireIDStr)
	if err != nil {
		if isJSON {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		ctx.Redirect(http.StatusSeeOther, "/tires")
		return
	}

	uid, err := authUserIDUint(ctx)
	if err != nil {
		if isJSON {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
			return
		}
		ctx.Redirect(http.StatusSeeOther, "/tires")
		return
	}

	tp, created, err := h.Repository.AddTireToCartAPI(uint(tireID), uid)
	if err != nil {
		if isJSON {
			if errors.Is(err, repository.ErrNotFound) {
				h.errorHandler(ctx, http.StatusNotFound, err)
			} else if errors.Is(err, repository.ErrAlreadyExists) {
				h.errorHandler(ctx, http.StatusConflict, err)
			} else {
				h.errorHandler(ctx, http.StatusInternalServerError, err)
			}
			return
		}
		ctx.Redirect(http.StatusSeeOther, "/tires")
		return
	}

	if isJSON {
		status := http.StatusOK
		if created {
			ctx.Header("Location", fmt.Sprintf("/api/tire-pressures/%d", tp.TirePressureID))
			status = http.StatusCreated
		}
		ctx.JSON(status, gin.H{
			"status":           "success",
			"tire_pressure_id": tp.TirePressureID,
		})
		return
	}

	// HTML-редирект
	redirectTo := ctx.GetHeader("Referer")
	if redirectTo == "" {
		redirectTo = "/tires"
	}
	ctx.Redirect(http.StatusSeeOther, redirectTo)
}

// DeleteTireFromPressure Убрать шину из заявки
// @Summary Удалить шину из заявки
// @Tags tire-pressure-entries
// @Produce json
// @Param tire_id path int true "ID шины"
// @Param tire_pressure_id path int true "ID заявки"
// @Security ApiKeyAuth
// @Success 200 {object} serializer.TirePressureListJSON "Обновлённая заявка"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tire-pressure-entries/{tire_id}/{tire_pressure_id} [delete]
func (h *Handler) DeleteTireFromPressure(ctx *gin.Context) {
	tireID, err := strconv.Atoi(ctx.Param("tire_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный tire_id"))
		return
	}
	tpID, err := strconv.Atoi(ctx.Param("tire_pressure_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный tire_pressure_id"))
		return
	}

	uid, err := authUserIDUint(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	t, err := h.Repository.DeleteTireFromCartAPI(tireID, tpID, int(uid))
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

// UpdateTirePressureEntry Обновить параметры записи в заявке
// @Summary Обновить запись в заявке
// @Description Коэффициент покрытия, давление и другие параметры для пары шина–заявка.
// @Tags tire-pressure-entries
// @Accept json
// @Produce json
// @Param tire_id path int true "ID шины"
// @Param tire_pressure_id path int true "ID заявки"
// @Param body body serializer.TirePressureEntryUpdateJSON true "Поля для обновления"
// @Security ApiKeyAuth
// @Success 200 {object} serializer.TirePressureEntryUpdateJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tire-pressure-entries/{tire_id}/{tire_pressure_id} [put]
func (h *Handler) UpdateTirePressureEntry(ctx *gin.Context) {
	tireID, err := strconv.Atoi(ctx.Param("tire_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный tire_id"))
		return
	}
	tpID, err := strconv.Atoi(ctx.Param("tire_pressure_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный tire_pressure_id"))
		return
	}

	var j serializer.TirePressureEntryUpdateJSON
	if err := ctx.ShouldBindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	uid, err := authUserIDUint(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	item, err := h.Repository.UpdateTireInCartAPI(tireID, tpID, j, int(uid))
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

	ctx.JSON(http.StatusOK, serializer.TirePressureEntryUpdateJSON{
		CoatingCoefficient: item.CoatingCoefficient,
		Pressure:           item.Pressure,
		// добавьте другие поля при необходимости
	})
}
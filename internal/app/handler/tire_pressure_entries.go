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

// ─── API: Tire Pressure Entries (3 метода) ──────────────────────────────────

// ✅ POST /api/tire-pressure-entries/add/:tire_id
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
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	tireID, err := strconv.Atoi(tireIDStr)
	if err != nil {
		if isJSON {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	creatorID := uint(repository.GetUserID())
	tp, created, err := h.Repository.AddTireToCartAPI(uint(tireID), creatorID)
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
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	if isJSON {
		status := http.StatusOK
		if created {
			ctx.Header("Location", fmt.Sprintf("/api/tire-pressures/%d", tp.TirePressureID))
			status = http.StatusCreated
		}
		ctx.JSON(status, gin.H{"status": "success", "tire_pressure_id": tp.TirePressureID})
		return
	}

	redirectTo := ctx.GetHeader("Referer")
	if redirectTo == "" {
		redirectTo = "/"
	}
	ctx.Redirect(http.StatusSeeOther, redirectTo)
}

// ✅ PUT /api/tire-pressure-entries/:tire_id/:tire_pressure_id
func (h *Handler) UpdateTirePressureEntry(ctx *gin.Context) {
	tireID, err := strconv.Atoi(ctx.Param("tire_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	tpID, err := strconv.Atoi(ctx.Param("tire_pressure_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var j serializer.TirePressureEntryUpdateJSON
	if err := ctx.ShouldBindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	item, err := h.Repository.UpdateTireInCartAPI(tireID, tpID, j)
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

	ctx.JSON(http.StatusOK, gin.H{
		"coating_coefficient": item.CoatingCoefficient,
		"pressure":            item.Pressure,
	})
}

// ✅ DELETE /api/tire-pressure-entries/:tire_id/:tire_pressure_id
func (h *Handler) DeleteTireFromPressure(ctx *gin.Context) {
	tireID, err := strconv.Atoi(ctx.Param("tire_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	tpID, err := strconv.Atoi(ctx.Param("tire_pressure_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	t, err := h.Repository.DeleteTireFromCartAPI(tireID, tpID)
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

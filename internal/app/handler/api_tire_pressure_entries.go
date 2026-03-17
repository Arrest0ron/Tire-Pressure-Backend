package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"metoda/internal/app/repository"
	"metoda/internal/app/serializer"
)

// POST /api/tires/:id/add-to-tire-pressure
func (h *Handler) APIAddToCart(ctx *gin.Context) {
	tireID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id шины"))
		return
	}

	// ✅ SINGLETON: замена h.Repository.GetUserID() → repository.GetUserID()
	t, created, err := h.Repository.AddTireToCartAPI(uint(tireID), uint(repository.GetUserID()))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	status := http.StatusOK
	if created {
		ctx.Header("Location", fmt.Sprintf("/api/tire-pressures/%d", t.TirePressureID))
		status = http.StatusCreated
	}

	creatorLogin := h.Repository.GetCreatorLogin(t.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(t.ModeratorID)
	entriesCount := h.Repository.GetTireEntriesCount(t.TirePressureID)
	ctx.JSON(status, serializer.TirePressureToListJSON(t, creatorLogin, moderatorLogin, entriesCount))
}

// DELETE /api/tire-pressure-entries/:tire_id/:tire_pressure_id
func (h *Handler) APIDeleteFromCart(ctx *gin.Context) {
	tireID, err := strconv.Atoi(ctx.Param("tire_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный tire_id"))
		return
	}
	tirePressureID, err := strconv.Atoi(ctx.Param("tire_pressure_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный tire_pressure_id"))
		return
	}

	t, err := h.Repository.DeleteTireFromCartAPI(tireID, tirePressureID)
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

// PUT /api/tire-pressure-entries/:tire_id/:tire_pressure_id
func (h *Handler) APIUpdateCartItem(ctx *gin.Context) {
	tireID, err := strconv.Atoi(ctx.Param("tire_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный tire_id"))
		return
	}
	tirePressureID, err := strconv.Atoi(ctx.Param("tire_pressure_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный tire_pressure_id"))
		return
	}

	var j serializer.TirePressureEntryUpdateJSON
	if err := ctx.ShouldBindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	item, err := h.Repository.UpdateTireInCartAPI(tireID, tirePressureID, j)
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
	})
}

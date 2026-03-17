package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"metoda/internal/app/repository"
	"metoda/internal/app/serializer"
)

// GET /api/tire-pressures/cart
func (h *Handler) APIGetCart(ctx *gin.Context) {
	// ✅ SINGLETON: замена h.Repository.GetUserID() → repository.GetUserID()
	id, count, err := h.Repository.GetCartInfo(uint(repository.GetUserID()))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if id == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":      "no_draft",
			"tires_count": 0,
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.CartJSON{
		TirePressureID: id,
		TiresCount:     count,
	})
}

// GET /api/tire-pressures?from_date=YYYY-MM-DD&to_date=YYYY-MM-DD&status=... — список (кроме удалённых и черновика), логины создателя/модератора, фильтр по дате формирования и статусу.
func (h *Handler) APIGetTirePressures(ctx *gin.Context) {
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

	list, err := h.Repository.GetAllTirePressures(from, to, status)
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

// GET /api/tire-pressures/:id
func (h *Handler) APIGetTirePressure(ctx *gin.Context) {
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

// PUT /api/tire-pressures/:id
func (h *Handler) APIUpdateTirePressure(ctx *gin.Context) {
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

	t, err := h.Repository.UpdateTirePressureFields(id, j)
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

// PUT /api/tire-pressures/:id/form
func (h *Handler) APIFormTirePressure(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	t, err := h.Repository.FormTirePressureAPI(id)
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

// PUT /api/tire-pressures/:id/finish
func (h *Handler) APIFinishTirePressure(ctx *gin.Context) {
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

	t, err := h.Repository.FinishTirePressureAPI(id, j.Status)
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

// DELETE /api/tire-pressures/:id
func (h *Handler) APIDeleteTirePressure(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	if err := h.Repository.DeleteTirePressureAPI(id); err != nil {
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

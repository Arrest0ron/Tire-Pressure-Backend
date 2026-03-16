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

// GET /api/dendrochronologies/cart
func (h *Handler) APIGetCart(ctx *gin.Context) {
	id, count, err := h.Repository.GetCartInfo(uint(h.Repository.GetUserID()))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if id == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":              "no_draft",
			"constructions_count": 0,
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.CartJSON{
		DendrochronologyID: id,
		ConstructionsCount: count,
	})
}

// GET /api/dendrochronologies?from_date=YYYY-MM-DD&to_date=YYYY-MM-DD&status=... — список (кроме удалённых и черновика), логины создателя/модератора, фильтр по дате формирования и статусу.
func (h *Handler) APIGetDendrochronologies(ctx *gin.Context) {
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

	list, err := h.Repository.GetAllDendrochronologies(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	resp := make([]serializer.DendrochronologyListJSON, 0, len(list))
	for _, d := range list {
		creatorLogin := h.Repository.GetCreatorLogin(d.CreatorID)
		moderatorLogin := h.Repository.GetModeratorLogin(d.ModeratorID)
		datedCount := h.Repository.GetDatedConstructionsCount(d.ID)
		resp = append(resp, serializer.DendrochronologyToListJSON(d, creatorLogin, moderatorLogin, datedCount))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GET /api/dendrochronologies/:id
func (h *Handler) APIGetDendrochronology(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	d, err := h.Repository.GetDendrochronologyByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	constructions, err := h.Repository.GetDendrochronologyConstructionsAPI(d.ID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	creatorLogin := h.Repository.GetCreatorLogin(d.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(d.ModeratorID)

	var dateFormed, dateCompleted *time.Time
	if d.DateFormed.Valid {
		dateFormed = &d.DateFormed.Time
	}
	if d.DateCompleted.Valid {
		dateCompleted = &d.DateCompleted.Time
	}
	var modLogin *string
	if moderatorLogin != "" {
		modLogin = &moderatorLogin
	}

	ctx.JSON(http.StatusOK, serializer.DendrochronologyDetailJSON{
		ID:             d.ID,
		Status:         d.Status,
		DateCreate:     d.DateCreate,
		DateFormed:     dateFormed,
		DateCompleted:  dateCompleted,
		CreatorLogin:   creatorLogin,
		ModeratorLogin: modLogin,
		TotalSamples:   d.TotalSamples,
		BuildDate:      d.BuildDate,
		Constructions:  constructions,
	})
}

// PUT /api/dendrochronologies/:id
func (h *Handler) APIUpdateDendrochronology(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	var j serializer.DendrochronologyUpdateJSON
	if err := ctx.ShouldBindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	d, err := h.Repository.UpdateDendrochronologyFields(id, j)
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

	creatorLogin := h.Repository.GetCreatorLogin(d.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(d.ModeratorID)
	datedCount := h.Repository.GetDatedConstructionsCount(d.ID)
	ctx.JSON(http.StatusOK, serializer.DendrochronologyToListJSON(d, creatorLogin, moderatorLogin, datedCount))
}

// PUT /api/dendrochronologies/:id/form
func (h *Handler) APIFormDendrochronology(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	d, err := h.Repository.FormDendrochronologyAPI(id)
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

	creatorLogin := h.Repository.GetCreatorLogin(d.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(d.ModeratorID)
	datedCount := h.Repository.GetDatedConstructionsCount(d.ID)
	ctx.JSON(http.StatusOK, serializer.DendrochronologyToListJSON(d, creatorLogin, moderatorLogin, datedCount))
}

// PUT /api/dendrochronologies/:id/finish
func (h *Handler) APIFinishDendrochronology(ctx *gin.Context) {
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

	d, err := h.Repository.FinishDendrochronologyAPI(id, j.Status)
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

	creatorLogin := h.Repository.GetCreatorLogin(d.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(d.ModeratorID)
	datedCount := h.Repository.GetDatedConstructionsCount(d.ID)
	ctx.JSON(http.StatusOK, serializer.DendrochronologyToListJSON(d, creatorLogin, moderatorLogin, datedCount))
}

// DELETE /api/dendrochronologies/:id
func (h *Handler) APIDeleteDendrochronology(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	if err := h.Repository.DeleteDendrochronologyAPI(id); err != nil {
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

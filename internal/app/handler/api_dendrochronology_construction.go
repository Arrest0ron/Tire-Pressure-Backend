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

// POST /api/constructions/:id/add-to-dendrochronology
func (h *Handler) APIAddToCart(ctx *gin.Context) {
	constructionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id конструкции"))
		return
	}

	d, created, err := h.Repository.AddConstructionToCartAPI(uint(constructionID), uint(h.Repository.GetUserID()))
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
		ctx.Header("Location", fmt.Sprintf("/api/dendrochronologies/%d", d.ID))
		status = http.StatusCreated
	}

	creatorLogin := h.Repository.GetCreatorLogin(d.CreatorID)
	moderatorLogin := h.Repository.GetModeratorLogin(d.ModeratorID)
	datedCount := h.Repository.GetDatedConstructionsCount(d.ID)
	ctx.JSON(status, serializer.DendrochronologyToListJSON(d, creatorLogin, moderatorLogin, datedCount))
}

// DELETE /api/dendrochronology-constructions/:construction_id/:dendrochronology_id
func (h *Handler) APIDeleteFromCart(ctx *gin.Context) {
	constructionID, err := strconv.Atoi(ctx.Param("construction_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный construction_id"))
		return
	}
	dendrochronologyID, err := strconv.Atoi(ctx.Param("dendrochronology_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный dendrochronology_id"))
		return
	}

	d, err := h.Repository.DeleteConstructionFromCartAPI(constructionID, dendrochronologyID)
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

// PUT /api/dendrochronology-constructions/:construction_id/:dendrochronology_id
func (h *Handler) APIUpdateCartItem(ctx *gin.Context) {
	constructionID, err := strconv.Atoi(ctx.Param("construction_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный construction_id"))
		return
	}
	dendrochronologyID, err := strconv.Atoi(ctx.Param("dendrochronology_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный dendrochronology_id"))
		return
	}

	var j serializer.DendroConstructionUpdateJSON
	if err := ctx.ShouldBindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	item, err := h.Repository.UpdateConstructionInCartAPI(constructionID, dendrochronologyID, j)
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

	ctx.JSON(http.StatusOK, serializer.DendroConstructionUpdateJSON{
		SamplesCount:   item.SamplesCount,
		CuttingDate:    item.CuttingDate,
		DateCorrection: item.DateCorrection,
	})
}

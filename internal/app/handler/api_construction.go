package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"metoda/internal/app/ds"
	"metoda/internal/app/repository"
	"metoda/internal/app/serializer"
)

// GET /api/constructions?query=...
func (h *Handler) APIGetConstructions(ctx *gin.Context) {
	query := ctx.Query("query")
	var (
		constructions []ds.Construction
		err           error
	)
	if query == "" {
		constructions, err = h.Repository.GetAllConstructions()
	} else {
		constructions, err = h.Repository.SearchConstructionsByTitle(query)
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	resp := make([]serializer.ConstructionJSON, 0, len(constructions))
	for _, c := range constructions {
		resp = append(resp, serializer.ConstructionToJSON(c))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GET /api/constructions/:id
func (h *Handler) APIGetConstruction(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	c, err := h.Repository.GetConstructionByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.ConstructionToJSON(*c))
}

// POST /api/constructions
// Accepts multipart/form-data: construction_title, use_life, description, image (file), video (file)
func (h *Handler) APICreateConstruction(ctx *gin.Context) {
	j := serializer.ConstructionJSON{
		ConstructionTitle: ctx.PostForm("construction_title"),
		UseLife:           ctx.PostForm("use_life"),
		Description:       ctx.PostForm("description"),
	}

	c, err := h.Repository.CreateConstruction(j)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Upload image if provided
	imageFile, imgErr := ctx.FormFile("image")
	if imgErr == nil && imageFile != nil {
		updated, uploadErr := h.Repository.UploadConstructionImage(ctx, int(c.ID), imageFile)
		if uploadErr != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, uploadErr)
			return
		}
		c = updated
	}

	// Upload video if provided
	videoFile, vidErr := ctx.FormFile("video")
	if vidErr == nil && videoFile != nil {
		updated, uploadErr := h.Repository.UploadConstructionVideo(ctx, int(c.ID), videoFile)
		if uploadErr != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, uploadErr)
			return
		}
		c = updated
	}

	ctx.Header("Location", fmt.Sprintf("/api/constructions/%d", c.ID))
	ctx.JSON(http.StatusCreated, serializer.ConstructionToJSON(c))
}

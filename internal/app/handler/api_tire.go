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

// GET /api/tires?query=...
func (h *Handler) APIGetTires(ctx *gin.Context) {
	query := ctx.Query("query")
	var (
		tires []ds.Tire
		err   error
	)
	if query == "" {
		tires, err = h.Repository.GetAllTires()
	} else {
		tires, err = h.Repository.SearchTiresByTitle(query)
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	resp := make([]serializer.TireJSON, 0, len(tires))
	for _, t := range tires {
		resp = append(resp, serializer.TireToJSON(t))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GET /api/tires/:id
func (h *Handler) APIGetTire(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный id"))
		return
	}

	t, err := h.Repository.GetTireByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.TireToJSON(*t))
}

// POST /api/tires
// Accepts multipart/form-data: tire_title, tire_material_coefficient, tire_thickness_coefficient, description, photo (file), video (file)
func (h *Handler) APICreateTire(ctx *gin.Context) {
	materialCoeff, err := strconv.ParseFloat(ctx.PostForm("tire_material_coefficient"), 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный коэффициент материала"))
		return
	}

	thicknessCoeff, err := strconv.ParseFloat(ctx.PostForm("tire_thickness_coefficient"), 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный коэффициент толщины"))
		return
	}

	j := serializer.TireJSON{
		TireTitle:                ctx.PostForm("tire_title"),
		TireMaterialCoefficient:  materialCoeff,
		TireThicknessCoefficient: thicknessCoeff,
		Description:              ctx.PostForm("description"),
	}

	t, err := h.Repository.CreateTire(j)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Upload photo if provided
	photoFile, photoErr := ctx.FormFile("photo")
	if photoErr == nil && photoFile != nil {
		updated, uploadErr := h.Repository.UploadTirePhoto(ctx, int(t.TireID), photoFile)
		if uploadErr != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, uploadErr)
			return
		}
		t = updated
	}

	// Upload video if provided
	videoFile, vidErr := ctx.FormFile("video")
	if vidErr == nil && videoFile != nil {
		updated, uploadErr := h.Repository.UploadTireVideo(ctx, int(t.TireID), videoFile)
		if uploadErr != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, uploadErr)
			return
		}
		t = updated
	}

	ctx.Header("Location", fmt.Sprintf("/api/tires/%d", t.TireID))
	ctx.JSON(http.StatusCreated, serializer.TireToJSON(t))
}

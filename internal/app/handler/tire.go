package handler

import (
	"net/http"
	"strconv"
	"web_backend/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetTires(ctx *gin.Context) {
	var tires []ds.Tire
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		tires, err = h.Repository.GetTires()
	} else {
		tires, err = h.Repository.GetTiresByTitle(searchQuery)
	}
	if err != nil {
		logrus.Error(err)
	}

	creatorID := uint(1)
	appCount := h.Repository.GetTirePressureCount(creatorID)

	ctx.HTML(http.StatusOK, "tires.html", gin.H{
		"tires":        tires,
		"query":        searchQuery,
		"entryCount":   appCount,
		"minioUrl":     h.Config.MinioURL,
	})
}

func (h *Handler) GetTire(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tire, err := h.Repository.GetTire(id)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "tire.html", gin.H{
		"tire":     tire,
		"minioUrl": h.Config.MinioURL,
	})
}

func (h *Handler) GetTirePressure(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	creatorID := uint(1)
	entries, err := h.Repository.GetTirePressure(id, creatorID)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Get the main request data (TirePressure) using the new repository function
	request, err := h.Repository.GetTirePressureByID(id, creatorID)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Get all tires for the dropdown
	tires, err := h.Repository.GetTires()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "tire_pressure.html", gin.H{
		"request_entries": entries,
		"request":         request,
		"tires":           tires,
		"minioUrl":        h.Config.MinioURL,
	})
}

func (h *Handler) AddToTirePressure(ctx *gin.Context) {
	tireIDStr := ctx.PostForm("tire_id")
	tireID, err := strconv.Atoi(tireIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(1)

	err = h.Repository.AddTireToPressure(tireID, creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
}

func (h *Handler) DeleteTirePressure(ctx *gin.Context) {
	tirePressureIDStr := ctx.PostForm("tire_pressure_id")
	tirePressureID, err := strconv.Atoi(tirePressureIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteTirePressure(uint(tirePressureID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) CompleteTirePressure(ctx *gin.Context) {
	tirePressureIDStr := ctx.PostForm("tire_pressure_id")
	tirePressureID, err := strconv.Atoi(tirePressureIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.CompleteTirePressure(uint(tirePressureID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Stay on the current page instead of redirecting to main page
	ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
}
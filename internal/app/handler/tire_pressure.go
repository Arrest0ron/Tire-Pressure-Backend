package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AddTireToPressure - Add tire to request (using ORM)
func (h *Handler) AddTireToPressure(ctx *gin.Context) {
	tireIDStr := ctx.PostForm("tire_id")
	tireID, err := strconv.Atoi(tireIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	creatorID := uint(1)
	err = h.Repository.AddTireToPressure(tireID, creatorID)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/tire_pressure/1")
}

// UpdateCoatingCoeff - Update coating coefficient for a tire
func (h *Handler) UpdateCoatingCoeff(ctx *gin.Context) {
	appIDStr := ctx.PostForm("app_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tireIDStr := ctx.PostForm("tire_id")
	tireID, err := strconv.Atoi(tireIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	coeffStr := ctx.PostForm("coating_coeff")
	coeff, err := strconv.ParseFloat(coeffStr, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Repository.UpdateCoatingCoeff(uint(appID), uint(tireID), coeff)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/tire_pressure/1")
}

// DeleteTireEntry - Delete specific tire entry from request
func (h *Handler) DeleteTireEntry(ctx *gin.Context) {
	appIDStr := ctx.PostForm("app_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tireIDStr := ctx.PostForm("tire_id")
	tireID, err := strconv.Atoi(tireIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Repository.DeleteTireEntry(uint(appID), uint(tireID))
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/tire_pressure/1")
}
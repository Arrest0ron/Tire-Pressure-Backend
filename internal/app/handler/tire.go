package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"metoda/internal/app/ds"
	"metoda/internal/app/repository"
	"metoda/internal/app/serializer"
)

// ─── HTML Pages (2 страницы: главная + шина) ────────────────────────────────

// ✅ GET / — Главная страница (список шин)
func (h *Handler) Index(ctx *gin.Context) {
	var tires []ds.Tire
	var err error
	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		tires, err = h.Repository.GetAllTires()
	} else {
		tires, err = h.Repository.SearchTiresByTitle(searchQuery)
	}
	if err != nil {
		tires = []ds.Tire{}
	}

	cartCount := h.Repository.GetCartCount(uint(repository.GetUserID()))
	tirePressureID := h.Repository.GetDraftTirePressureID(uint(repository.GetUserID()))

	ctx.HTML(http.StatusOK, "mainpage.html", gin.H{
		"tires":          tires,
		"query":          searchQuery,
		"cartCount":      cartCount,
		"tirePressureID": tirePressureID,
		"minioBase":      minioBaseURL,
	})
}

// ✅ GET /tire/:id — Страница одной шины
func (h *Handler) TirePage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	tire, err := h.Repository.GetTireByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "productpage.html", gin.H{
		"tire":      tire,
		"minioBase": minioBaseURL,
	})
}

// ─── API: Tires (3 метода) ──────────────────────────────────────────────────

// ✅ GET /api/tires?query=...
func (h *Handler) GetTires(ctx *gin.Context) {
	query := ctx.Query("query")
	var tires []ds.Tire
	var err error
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

// ✅ GET /api/tires/:id
func (h *Handler) GetTire(ctx *gin.Context) {
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

// ✅ POST /api/tires
func (h *Handler) CreateTire(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")
	var j serializer.TireJSON
	if strings.HasPrefix(contentType, "application/json") {
		if err := ctx.BindJSON(&j); err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	} else {
		materialCoeff, _ := strconv.ParseFloat(ctx.PostForm("tire_material_coefficient"), 64)
		thicknessCoeff, _ := strconv.ParseFloat(ctx.PostForm("tire_thickness_coefficient"), 64)
		j = serializer.TireJSON{
			TireTitle:                ctx.PostForm("tire_title"),
			TireMaterialCoefficient:  materialCoeff,
			TireThicknessCoefficient: thicknessCoeff,
			Description:              ctx.PostForm("description"),
		}
	}
	t, err := h.Repository.CreateTire(j)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if photoFile, err := ctx.FormFile("photo"); err == nil {
		t, _ = h.Repository.UploadTirePhoto(ctx, int(t.TireID), photoFile)
	}
	if videoFile, err := ctx.FormFile("video"); err == nil {
		t, _ = h.Repository.UploadTireVideo(ctx, int(t.TireID), videoFile)
	}
	ctx.Header("Location", fmt.Sprintf("/api/tires/%d", t.TireID))
	ctx.JSON(http.StatusCreated, serializer.TireToJSON(t))
}

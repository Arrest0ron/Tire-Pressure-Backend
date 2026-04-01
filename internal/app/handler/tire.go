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

// ─── HTML Pages ────────────────────────────────────────────────────────────

// Index Главная страница со списком шин
// @Summary Главная страница (список шин)
// @Tags tires-html
// @Produce html
// @Param query query string false "Поиск по названию шины"
// @Success 200 {string} string "HTML-страница"
// @Router /tires [get]
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

	cartCount := h.Repository.GetCartCount(uint(h.Repository.GetUserID()))
	tirePressureID := h.Repository.GetDraftTirePressureID(uint(h.Repository.GetUserID()))

	ctx.HTML(http.StatusOK, "mainpage.html", gin.H{
		"tires":          tires,
		"query":          searchQuery,
		"cartCount":      cartCount,
		"tirePressureID": tirePressureID,
		"minioBase":      h.getMinioURL(),
	})
}

// TirePage Страница одной шины
// @Summary Страница шины по ID
// @Tags tires-html
// @Produce html
// @Param id path int true "ID шины"
// @Success 200 {string} string "HTML-страница"
// @Failure 404 {object} map[string]string
// @Router /tire/{id} [get]
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
		"minioBase": h.getMinioURL(),
	})
}

// ─── API: Tires ────────────────────────────────────────────────────────────

// GetTires Список шин (API)
// @Summary Список шин
// @Tags tires
// @Produce json
// @Param query query string false "Поиск по названию шины"
// @Success 200 {array} serializer.TireJSON
// @Failure 500 {object} map[string]string
// @Router /api/tires [get]
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

// GetTire Одна шина по ID (API)
// @Summary Шина по ID
// @Tags tires
// @Produce json
// @Param id path int true "ID шины"
// @Success 200 {object} serializer.TireJSON
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/tires/{id} [get]
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

// CreateTire Создание шины (только модератор)
// @Summary Создать шину
// @Description multipart/form-data или application/json; только для модераторов
// @Tags tires
// @Accept mpfd
// @Produce json
// @Param tire_title formData string false "Название шины"
// @Param tire_material_coefficient formData number false "Коэффициент материала"
// @Param tire_thickness_coefficient formData number false "Коэффициент толщины"
// @Param description formData string false "Описание"
// @Param photo formData file false "Фото шины"
// @Param video formData file false "Видео шины"
// @Success 201 {object} serializer.TireJSON
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string "Требуется роль модератора"
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/tires [post]
func (h *Handler) CreateTire(ctx *gin.Context) {
	// ✅ Получаем userID из контекста (мидлвар уже проверил токен)
	uid, err := authUserIDUint(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	contentType := ctx.GetHeader("Content-Type")
	var j serializer.TireJSON

	// Поддержка и JSON, и form-data
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

	// ✅ Передаем userID в репозиторий для проверки прав модератора
	t, err := h.Repository.CreateTire(j, int(uid))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Загрузка фото (если есть)
	if photoFile, err := ctx.FormFile("photo"); err == nil && photoFile != nil {
		updated, uploadErr := h.Repository.UploadTirePhoto(ctx, int(t.TireID), photoFile)
		if uploadErr != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, uploadErr)
			return
		}
		t = updated
	}

	// Загрузка видео (если есть)
	if videoFile, err := ctx.FormFile("video"); err == nil && videoFile != nil {
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
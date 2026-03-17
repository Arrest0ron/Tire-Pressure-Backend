package handler

import (
	"metoda/internal/app/ds"
	"metoda/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetTires(ctx *gin.Context) {
	var tires []ds.Tire
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		tires, err = h.Repository.GetAllTires()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		tires, err = h.Repository.SearchTiresByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	// ✅ SINGLETON: замена creatorID → repository.GetUserID()
	cartCount := h.Repository.GetCartCount(uint(repository.GetUserID()))
	// ✅ Как в оригинале: передаём ID активной заявки (черновика)
	tirePressureID := h.Repository.GetDraftTirePressureID(uint(repository.GetUserID()))

	// ✅ Логирование для отладки
	logrus.Info("cartCount:", cartCount, "tirePressureID:", tirePressureID)

	ctx.HTML(http.StatusOK, "mainpage.html", gin.H{
		"tires":          tires,
		"query":          searchQuery,
		"cartCount":      cartCount,
		"tirePressureID": tirePressureID, // ✅ Для иконки корзины
		"minioBase":      minioBaseURL,
	})
}

func (h *Handler) GetTire(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	tire, err := h.Repository.GetTireByID(id)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	if tire == nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.HTML(http.StatusOK, "productpage.html", gin.H{
		"tire":      tire,
		"minioBase": minioBaseURL,
	})
}

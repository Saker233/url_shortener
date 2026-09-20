package handler

import (
	"net/http"
	"os"

	db "url_shortener/internal/database/sqlc"
	"url_shortener/internal/service"
	"url_shortener/internal/util"

	"github.com/gin-gonic/gin"
)

type createURLRequest struct {
	Original_URL string `json:"original_url" binding:"required,url"`
}

type createURLResponse struct {
	Short_Code string `json:"short_code"`
}

func InitServer() {
	r := gin.Default()
	r.GET("/health", getHealth)
	r.POST("/url", createURL)

	_ = r.Run(os.Getenv("PORT"))
}

func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Healthy",
	})
}

func createURL(c *gin.Context) {
	var req createURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	shortCode, err := service.CreateURL()
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateURLParams{
		ShortCode:   shortCode,
		OriginalUrl: req.Original_URL,
	}

	url, err := util.GetQueries().CreateURL(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := db.Url(url)
	c.JSON(http.StatusOK, rsp)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

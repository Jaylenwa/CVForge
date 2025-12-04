package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type polishReq struct {
	Text string `json:"text"`
	Tone string `json:"tone"`
}
type summaryReq struct {
	JobTitle string `json:"jobTitle"`
	Skills   string `json:"skills"`
}

func RegisterAIRoutes(r *gin.RouterGroup) {
	r.POST("/polish", func(c *gin.Context) {
		var req polishReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		out := strings.TrimSpace(req.Text)
		if out != "" {
			out = strings.ToUpper(out[:1]) + out[1:]
		}
		c.JSON(http.StatusOK, gin.H{"text": out})
	})

	r.POST("/summary", func(c *gin.Context) {
		var req summaryReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		text := "Experienced " + req.JobTitle + " with skills in " + req.Skills + ". Delivers impact through collaboration and ownership."
		c.JSON(http.StatusOK, gin.H{"text": text})
	})
}

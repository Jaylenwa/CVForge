package middleware

import (
    "github.com/gin-gonic/gin"
)

// helper to inject Auth from main
func AuthFromGroup(_ *gin.RouterGroup) gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"openresume/cache"
	"openresume/config"
	"openresume/db"
	"openresume/handlers"
	"openresume/metrics"
	"openresume/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	sqlDB, err := db.InitMySQL(cfg)
	if err != nil {
		log.Fatalf("mysql init error: %v", err)
	}
	defer sqlDB.Close()

	rdb := cache.InitRedis(cfg)
	defer func() { _ = rdb.Close() }()

	router := gin.New()
	router.Use(middleware.RequestID(), middleware.Logger(), gin.Recovery(), metrics.Middleware())
	router.Use(middleware.CORS(cfg))

	api := router.Group("/api/v1")

	handlers.RegisterAuthRoutes(api, cfg, rdb, db.Gorm())
	handlers.RegisterUserRoutesReal(api, db.Gorm(), middleware.Auth(cfg))
	handlers.RegisterTemplateRoutes(api, db.Gorm(), rdb)
	api.Use(middleware.RateLimitUser(rdb, 120, time.Minute))
	handlers.RegisterResumeRoutes(api, db.Gorm(), middleware.Auth(cfg))
	handlers.RegisterUploadRoutes(api)
	handlers.RegisterShareRoutes(api, db.Gorm(), rdb, middleware.Auth(cfg))
	admin := api.Group("/admin")
	admin.Use(middleware.Auth(cfg), middleware.RequireRole("admin"))
	handlers.RegisterAdminUserRoutes(api, db.Gorm(), middleware.Auth(cfg), middleware.RequireRole("admin"))
	handlers.RegisterAdminResumeRoutes(api, db.Gorm(), middleware.Auth(cfg), middleware.RequireRole("admin"), rdb)
	handlers.RegisterAdminTemplateRoutes(api, db.Gorm(), middleware.Auth(cfg), middleware.RequireRole("admin"), rdb)
	handlers.RegisterAdminShareRoutes(api, db.Gorm(), middleware.Auth(cfg), middleware.RequireRole("admin"), rdb)
	air := api.Group("/ai")
	air.Use(middleware.RateLimit(rdb, 10, time.Minute))
	handlers.RegisterAIRoutes(air)
	handlers.RegisterHealthRoutes(api)
	api.GET("/metrics", metrics.Handler())
	handlers.RegisterPDFRoutes(api, db.Gorm(), middleware.Auth(cfg))

	router.Static("/public/uploads", "./uploads")

	addr := ":" + cfg.Port
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
	}
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

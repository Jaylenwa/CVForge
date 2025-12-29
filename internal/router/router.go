package router

import (
	"log"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/config"
	"openresume/internal/middleware"
	"openresume/internal/module/ai"
	"openresume/internal/module/auth"
	conf "openresume/internal/module/config"
	"openresume/internal/module/health"
	"openresume/internal/module/pdf"
	"openresume/internal/module/resume"
	"openresume/internal/module/share"
	"openresume/internal/module/stats"
	"openresume/internal/module/template"
	"openresume/internal/module/upload"
	"openresume/internal/module/user"
	"openresume/internal/pkg/metrics"

	"github.com/gin-gonic/gin"
)

func Init(cfg config.Config) *gin.Engine {
	router := gin.New()
	router.Use(middleware.RequestID(), middleware.Logger(), gin.Recovery(), metrics.Middleware())

	api := router.Group("/api/v1")

	confService := conf.NewService()
	confHandler := conf.NewHandler()
	_ = confService.EnsureDefaults(cfg)
	router.Use(middleware.CORS(cfg, confService))

	authH := auth.NewHandler(cfg, confService)
	userAuth := middleware.Auth(cfg)
	userAdmin := middleware.RequireRole(common.RoleAdmin)
	templateH := template.NewHandler()
	api.GET("/templates", templateH.ListAll)
	api.GET("/templates/:id", templateH.GetByID)

	authR := api.Group("/auth")
	authR.GET("/config", confHandler.GetPublic)
	authR.GET("/wechat/redirect", middleware.RateLimit(10, time.Minute), authH.WeChatRedirect(cfg))
	authR.GET("/wechat/callback", middleware.RateLimit(30, time.Minute), authH.WeChatCallback(cfg))
	authR.GET("/github/redirect", middleware.RateLimit(10, time.Minute), authH.GithubRedirect(cfg))
	authR.GET("/github/callback", middleware.RateLimit(30, time.Minute), authH.GithubCallback(cfg))
	authR.POST("/wechat/consume-ott", authH.ConsumeOTT)
	authR.POST("/send-code", middleware.RateLimit(3, time.Minute), authH.SendCode)
	authR.POST("/register", middleware.RateLimit(5, time.Minute), authH.Register)
	authR.POST("/login", middleware.RateLimit(5, time.Minute), authH.Login)
	authR.POST("/refresh", authH.Refresh)
	authR.POST("/logout", authH.Logout)

	api.Use(middleware.RateLimitUser(120, time.Minute))
	api.Use(middleware.DailyUV("/api/v1/healthz", "/api/v1/metrics"))
	g := api.Group("")
	g.Use(userAuth)

	resumeH := resume.NewHandler()
	g.GET("/resumes", resumeH.List)
	g.POST("/resumes", resumeH.Create)
	g.GET("/resumes/:id", resumeH.Get)
	g.PUT("/resumes/:id", resumeH.Update)
	g.DELETE("/resumes/:id", resumeH.Delete)

	uploadH := upload.NewHandler(cfg, confService)
	g.POST("/upload/avatar", uploadH.UploadAvatar)

	shareH := share.NewHandler()
	g.POST("/resumes/:id/publish", shareH.PublishResume)
	api.GET("/public/resumes/:slug", shareH.GetPublic)

	userH := user.NewHandler()
	g.GET("/users/me", userH.Me)
	g.PUT("/users/profile", userH.UpdateProfile)
	g.PUT("/users/password", userH.UpdatePassword)

	adm := api.Group("/admin")
	adm.Use(userAuth, userAdmin)
	userAdmH := user.NewAdminHandler()
	adm.GET("/users", userAdmH.AdminList)
	adm.GET("/users/:id", userAdmH.AdminGet)
	adm.PATCH("/users/:id", userAdmH.AdminPatch)
	adm.POST("/users/:id/reset-password", userAdmH.AdminResetPassword)
	adm.POST("/users/:id/ban", userAdmH.AdminBan)
	adm.POST("/users/:id/unban", userAdmH.AdminUnban)

	resumeAdmH := resume.NewAdminHandler()
	adm.GET("/resumes", resumeAdmH.AdminList)
	adm.GET("/resumes/:id", resumeAdmH.AdminGet)
	adm.DELETE("/resumes/:id", resumeAdmH.AdminDelete)
	adm.PATCH("/resumes/:id/visibility", resumeAdmH.AdminUpdateVisibility)

	templateAdmH := template.NewAdminHandler()
	adm.POST("/templates", templateAdmH.AdminCreate)
	adm.PATCH("/templates/:id", templateAdmH.AdminPatch)
	adm.DELETE("/templates/:id", templateAdmH.AdminDelete)

	shareAdmH := share.NewAdminHandler()
	adm.GET("/share-links", shareAdmH.AdminList)
	adm.PATCH("/share-links/:slug", shareAdmH.AdminUpdate)
	adm.DELETE("/share-links/:slug", shareAdmH.AdminDelete)

	statsH := stats.NewHandler()
	adm.GET("/stats", statsH.AdminStats)

	adm.GET("/configs", confHandler.AdminList)
	adm.PUT("/configs", confHandler.AdminUpdate)

	air := api.Group("/ai")
	air.Use(middleware.RateLimit(10, time.Minute))
	aiH := ai.NewHandler()
	air.POST("/polish", aiH.Polish)
	air.POST("/summary", aiH.Summary)
	healthH := &health.Handler{}
	api.GET("/healthz", healthH.Healthz)
	api.GET("/metrics", metrics.Handler())
	pdfH := pdf.NewHandler(cfg, confService)
	g.POST("/resumes/:id/pdf", pdfH.GeneratePDF)
	g.POST("/resumes/:id/image", pdfH.GenerateImage)

	router.Static("/public/uploads", "./uploads")
	log.Println("router initialized")
	return router
}

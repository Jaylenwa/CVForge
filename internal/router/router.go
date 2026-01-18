package router

import (
	"time"

	"openresume/internal/common"
	"openresume/internal/middleware"
	"openresume/internal/module/ai"
	"openresume/internal/module/auth"
	"openresume/internal/module/catalog"
	conf "openresume/internal/module/config"
	"openresume/internal/module/health"
	"openresume/internal/module/pdf"
	"openresume/internal/module/resume"
	"openresume/internal/module/share"
	"openresume/internal/module/stats"
	"openresume/internal/module/template"
	"openresume/internal/module/upload"
	"openresume/internal/module/user"
	"openresume/internal/pkg/logger"
	"openresume/internal/pkg/metrics"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.New()
	router.Use(middleware.RequestID(), middleware.Logger(), gin.Recovery(), metrics.Middleware())
	api := router.Group("/api/v1")
	router.Use(middleware.CORS())

	confHandler := conf.NewHandler()
	authH := auth.NewHandler()
	userAuth := middleware.Auth()
	userAdmin := middleware.RequireRole(common.RoleAdmin)
	templateH := template.NewHandler()
	api.GET("/templates", templateH.ListAll)
	api.GET("/templates/:id", templateH.GetByID)

	catalogH := catalog.NewHandler()
	api.GET("/job-categories", catalogH.ListJobCategories)
	api.GET("/job-roles", catalogH.ListJobRoles)
	api.GET("/template-variants", catalogH.ListTemplateVariants)
	api.GET("/content-presets/:id", catalogH.GetContentPreset)

	authR := api.Group("/auth")
	authR.GET("/config", confHandler.GetPublic)
	authR.GET("/wechat/redirect", middleware.RateLimit(10, time.Minute), authH.WeChatRedirect())
	authR.GET("/wechat/callback", middleware.RateLimit(30, time.Minute), authH.WeChatCallback())
	authR.GET("/github/redirect", middleware.RateLimit(10, time.Minute), authH.GithubRedirect())
	authR.GET("/github/callback", middleware.RateLimit(30, time.Minute), authH.GithubCallback())
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
	g.POST("/resumes/from-variant", resumeH.CreateFromVariant)
	g.GET("/resumes/:id", resumeH.Get)
	g.PUT("/resumes/:id", resumeH.Update)
	g.DELETE("/resumes/:id", resumeH.Delete)

	uploadH := upload.NewHandler()
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

	catalogAdmH := catalog.NewAdminHandler()
	adm.POST("/catalog/import-seed", catalogAdmH.AdminImportSeed)
	adm.GET("/catalog/job-categories", catalogAdmH.AdminListJobCategories)
	adm.POST("/catalog/job-categories", catalogAdmH.AdminCreateJobCategory)
	adm.PATCH("/catalog/job-categories/:id", catalogAdmH.AdminPatchJobCategory)
	adm.DELETE("/catalog/job-categories/:id", catalogAdmH.AdminDeleteJobCategory)
	adm.GET("/catalog/job-roles", catalogAdmH.AdminListJobRoles)
	adm.POST("/catalog/job-roles", catalogAdmH.AdminCreateJobRole)
	adm.PATCH("/catalog/job-roles/:id", catalogAdmH.AdminPatchJobRole)
	adm.DELETE("/catalog/job-roles/:id", catalogAdmH.AdminDeleteJobRole)
	adm.GET("/catalog/content-presets", catalogAdmH.AdminListContentPresets)
	adm.POST("/catalog/content-presets", catalogAdmH.AdminCreateContentPreset)
	adm.PATCH("/catalog/content-presets/:id", catalogAdmH.AdminPatchContentPreset)
	adm.DELETE("/catalog/content-presets/:id", catalogAdmH.AdminDeleteContentPreset)
	adm.GET("/catalog/template-variants", catalogAdmH.AdminListTemplateVariants)
	adm.POST("/catalog/template-variants", catalogAdmH.AdminCreateTemplateVariant)
	adm.PATCH("/catalog/template-variants/:id", catalogAdmH.AdminPatchTemplateVariant)
	adm.DELETE("/catalog/template-variants/:id", catalogAdmH.AdminDeleteTemplateVariant)
	adm.POST("/catalog/template-variants/generate", catalogAdmH.AdminGenerateTemplateVariants)

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
	pdfH := pdf.NewHandler()
	g.POST("/resumes/:id/image", pdfH.GenerateImage)
	g.POST("/pdf/exports", pdfH.SubmitExport)
	g.GET("/pdf/exports/:job_id", pdfH.ExportStatus)
	api.GET("/pdf/exports/:job_id/download", pdfH.ExportDownload)

	router.Static("/public/uploads", "./uploads")
	logger.WithCtx(nil).Info("router initialized")
	return router
}

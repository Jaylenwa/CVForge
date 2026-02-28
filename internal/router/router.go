package router

import (
	"time"

	"cvforge/internal/common"
	"cvforge/internal/middleware"
	"cvforge/internal/module/ai"
	"cvforge/internal/module/auth"
	conf "cvforge/internal/module/config"
	confadmin "cvforge/internal/module/config/admin"
	"cvforge/internal/module/health"
	"cvforge/internal/module/library"
	"cvforge/internal/module/pdf"
	"cvforge/internal/module/preset"
	presetadmin "cvforge/internal/module/preset/admin"
	"cvforge/internal/module/resume"
	resumeadmin "cvforge/internal/module/resume/admin"
	"cvforge/internal/module/share"
	shareadmin "cvforge/internal/module/share/admin"
	statsadmin "cvforge/internal/module/stats/admin"
	"cvforge/internal/module/taxonomy"
	taxonomyadmin "cvforge/internal/module/taxonomy/admin"
	"cvforge/internal/module/template"
	templateadmin "cvforge/internal/module/template/admin"
	"cvforge/internal/module/upload"
	"cvforge/internal/module/user"
	useradmin "cvforge/internal/module/user/admin"
	"cvforge/internal/pkg/logger"
	"cvforge/internal/pkg/metrics"

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

	taxH := taxonomy.NewHandler()
	api.GET("/taxonomy/categories", taxH.ListCategories)
	api.GET("/taxonomy/roles", taxH.ListRoles)
	libraryH := library.NewHandler()
	api.GET("/library/templates", libraryH.ListTemplates)
	presetH := preset.NewHandler()
	api.GET("/presets/:id", presetH.GetByID)

	authR := api.Group("/auth")
	authR.GET("/config", confHandler.GetPublic)
	authR.GET("/github/redirect", middleware.RateLimit(10, time.Minute), authH.GithubRedirect())
	authR.GET("/github/callback", middleware.RateLimit(30, time.Minute), authH.GithubCallback())
	authR.POST("/wechat/consume-ott", authH.ConsumeOTT)
	authR.POST("/wechat-mp/scene/create", middleware.RateLimit(30, time.Minute), authH.WeChatMPCreateScene)
	authR.GET("/wechat-mp/scene/:scene/status", middleware.RateLimit(300, time.Minute), authH.WeChatMPSceneStatus)
	authR.POST("/send-code", middleware.RateLimit(3, time.Minute), authH.SendCode)
	authR.POST("/register", middleware.RateLimit(5, time.Minute), authH.Register)
	authR.POST("/login", middleware.RateLimit(5, time.Minute), authH.Login)
	authR.POST("/refresh", authH.Refresh)
	authR.POST("/logout", authH.Logout)

	api.GET("/wechat/mp/callback", middleware.RateLimit(600, time.Minute), authH.WeChatMPCallbackGet)
	api.POST("/wechat/mp/callback", middleware.RateLimit(600, time.Minute), authH.WeChatMPCallbackPost)

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

	uploadH := upload.NewHandler()
	g.POST("/upload/avatar", uploadH.UploadAvatar)

	shareH := share.NewHandler()
	g.POST("/resumes/:id/publish", shareH.PublishResume)
	g.GET("/resumes/:id/share", shareH.GetSettings)
	g.PATCH("/resumes/:id/share", shareH.UpdateSettings)
	api.GET("/public/resumes/:slug", shareH.GetPublic)
	api.POST("/public/resumes/:slug/auth", middleware.RateLimit(10, time.Minute), shareH.AuthPublic)

	userH := user.NewHandler()
	g.GET("/users/me", userH.Me)
	g.PUT("/users/profile", userH.UpdateProfile)
	g.PUT("/users/password", userH.UpdatePassword)

	adm := api.Group("/admin")
	adm.Use(userAuth, userAdmin)
	userAdmH := useradmin.NewHandler()
	adm.GET("/users", userAdmH.AdminList)
	adm.GET("/users/:id", userAdmH.AdminGet)
	adm.PATCH("/users/:id", userAdmH.AdminPatch)
	adm.POST("/users/:id/reset-password", userAdmH.AdminResetPassword)
	adm.POST("/users/:id/ban", userAdmH.AdminBan)
	adm.POST("/users/:id/unban", userAdmH.AdminUnban)

	resumeAdmH := resumeadmin.NewHandler()
	adm.GET("/resumes", resumeAdmH.AdminList)
	adm.GET("/resumes/:id", resumeAdmH.AdminGet)
	adm.DELETE("/resumes/:id", resumeAdmH.AdminDelete)
	adm.PATCH("/resumes/:id/visibility", resumeAdmH.AdminUpdateVisibility)

	templateAdmH := templateadmin.NewHandler()
	adm.POST("/templates", templateAdmH.AdminCreate)
	adm.PATCH("/templates/:id", templateAdmH.AdminPatch)
	adm.DELETE("/templates/:id", templateAdmH.AdminDelete)

	shareAdmH := shareadmin.NewHandler()
	adm.GET("/share-links", shareAdmH.AdminList)
	adm.PATCH("/share-links/:slug", shareAdmH.AdminUpdate)
	adm.DELETE("/share-links/:slug", shareAdmH.AdminDelete)

	taxAdmH := taxonomyadmin.NewHandler()
	adm.GET("/taxonomy/categories", taxAdmH.AdminListCategories)
	adm.POST("/taxonomy/categories", taxAdmH.AdminCreateCategory)
	adm.PATCH("/taxonomy/categories/:id", taxAdmH.AdminPatchCategory)
	adm.DELETE("/taxonomy/categories/:id", taxAdmH.AdminDeleteCategory)
	adm.GET("/taxonomy/roles", taxAdmH.AdminListRoles)
	adm.POST("/taxonomy/roles", taxAdmH.AdminCreateRole)
	adm.PATCH("/taxonomy/roles/:id", taxAdmH.AdminPatchRole)
	adm.DELETE("/taxonomy/roles/:id", taxAdmH.AdminDeleteRole)

	presetAdmH := presetadmin.NewHandler()
	adm.GET("/presets", presetAdmH.AdminListPresets)
	adm.POST("/presets", presetAdmH.AdminCreatePreset)
	adm.PATCH("/presets/:id", presetAdmH.AdminPatchPreset)
	adm.DELETE("/presets/:id", presetAdmH.AdminDeletePreset)

	statsH := statsadmin.NewHandler()
	adm.GET("/stats", statsH.AdminStats)

	confAdmH := confadmin.NewHandler()
	adm.GET("/configs", confAdmH.AdminList)
	adm.PUT("/configs", confAdmH.AdminUpdate)

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

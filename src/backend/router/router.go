package router

import (
	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/config"
	"github.com/oversteplab/oversteplab/internal/handler"
	"github.com/oversteplab/oversteplab/internal/middleware"
	"github.com/oversteplab/oversteplab/internal/service"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())

	// Initialize services
	authSvc := service.NewAuthService(cfg)
	userSvc := service.NewUserService()
	vpsSvc := service.NewVPSService(cfg.JWTSecret)
	companySvc := service.NewCompanyService()
	orderSvc := service.NewOrderService()
	ticketSvc := service.NewTicketService()
	apiKeySvc := service.NewAPIKeyService()
	auditLogSvc := service.NewAuditLogService()
	billSvc := service.NewBillService()
	adminSvc := service.NewAdminService()
	announceSvc := service.NewAnnouncementService()

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(authSvc, userSvc)
	vpsHandler := handler.NewVPSHandler(vpsSvc)
	companyHandler := handler.NewCompanyHandler(companySvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	ticketHandler := handler.NewTicketHandler(ticketSvc)
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeySvc)
	auditLogHandler := handler.NewAuditLogHandler(auditLogSvc)
	billHandler := handler.NewBillHandler(billSvc)
	adminHandler := handler.NewAdminHandler(adminSvc)
	challengeHandler := handler.NewChallengeHandler(cfg.DBPath)
	announceHandler := handler.NewAnnouncementHandler(announceSvc)
	configHandler := handler.NewConfigHandler()

	// Auth routes (public)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// API v1 routes
	api := r.Group("/api/v1")

	// Public endpoints within API group
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Authenticated routes
	authGroup := api.Group("")
	authGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	authGroup.Use(middleware.AuditMiddleware())
	{
		authGroup.POST("/logout", authHandler.Logout)

		// User routes
		user := authGroup.Group("/user")
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)
		}
		authGroup.GET("/users/:id", userHandler.GetUserByID)

		// Company routes
		company := authGroup.Group("/company")
		{
			company.GET("/members", companyHandler.ListMembers)
			company.POST("/members", companyHandler.AddMember)
			company.PUT("/members/:id", companyHandler.UpdateMember)
			company.DELETE("/members/:id", companyHandler.DeleteMember)
			company.PUT("/members/:id/role", companyHandler.ChangeRole)
		}

		// VPS routes
		vps := authGroup.Group("/vps")
		{
			vps.GET("", vpsHandler.List)
			vps.POST("", vpsHandler.Create)
			vps.GET("/:id", vpsHandler.GetByID)
			vps.POST("/:id/start", vpsHandler.Start)
			vps.POST("/:id/stop", vpsHandler.Stop)
			vps.POST("/:id/restart", vpsHandler.Restart)
			vps.POST("/:id/reinstall", vpsHandler.Reinstall)
			vps.DELETE("/:id", vpsHandler.Delete)
			vps.GET("/:id/console", vpsHandler.Console)
		}

		// Public console view route (uses signed token, no JWT required)
		api.GET("/vps/:id/console/view", vpsHandler.ConsoleView)

		// Order routes
		orders := authGroup.Group("/orders")
		{
			orders.GET("", orderHandler.List)
			orders.GET("/:id", orderHandler.GetByID)
		}

		// Bill routes
		bills := authGroup.Group("/bills")
		{
			bills.GET("", billHandler.List)
			bills.POST("/recharge", billHandler.Recharge)
			bills.GET("/export", billHandler.Export)
		}

		// Ticket routes
		tickets := authGroup.Group("/tickets")
		{
			tickets.GET("", ticketHandler.List)
			tickets.POST("", ticketHandler.Create)
			tickets.GET("/:id", ticketHandler.GetByID)
			tickets.POST("/:id/reply", ticketHandler.Reply)
			tickets.PUT("/:id/close", ticketHandler.Close)
		}

		// API Key routes
		apikeys := authGroup.Group("/apikeys")
		{
			apikeys.GET("", apiKeyHandler.List)
			apikeys.POST("", apiKeyHandler.Create)
			apikeys.DELETE("/:id", apiKeyHandler.Delete)
		}

		// Audit log routes
		authGroup.GET("/audit-logs", auditLogHandler.List)

		// Announcement routes (public for authenticated users)
		authGroup.GET("/announcements", announceHandler.ListPublished)

		// Admin routes
		admin := authGroup.Group("/admin")
		admin.Use(middleware.RequireAdmin())
		{
			admin.GET("/users", adminHandler.ListUsers)
			admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
			admin.PUT("/users/:id/password", adminHandler.ResetUserPassword)
			admin.GET("/companies", adminHandler.ListCompanies)
			admin.GET("/vps", adminHandler.ListAllVPS)
			admin.GET("/audit-logs", adminHandler.ListAllLogs)
			admin.POST("/reset", challengeHandler.ResetDatabase)
			admin.POST("/announcements", announceHandler.Create)
			admin.PUT("/announcements/:id", announceHandler.Update)
			admin.DELETE("/announcements/:id", announceHandler.Delete)
			admin.GET("/config", configHandler.GetConfig)
			admin.PUT("/config", configHandler.UpdateConfig)
		}

		// Challenge routes
		challenges := authGroup.Group("/challenges")
		{
			challenges.GET("", challengeHandler.List)
			challenges.GET("/:id", challengeHandler.Detail)
			challenges.GET("/:id/hints/:level", challengeHandler.GetHint)
			challenges.POST("/:id/complete", challengeHandler.MarkComplete)
		}

		authGroup.GET("/security-mode", challengeHandler.GetSecurityMode)
		authGroup.PUT("/security-mode", challengeHandler.SetSecurityMode)
	}

	return r
}

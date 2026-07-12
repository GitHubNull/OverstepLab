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
	encodedHandler := handler.NewEncodedHandler(vpsSvc, userSvc, orderSvc, ticketSvc, apiKeySvc, companySvc)

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
	authGroup.Use(middleware.EncodingMiddleware())
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
		authGroup.GET("/users", userHandler.GetUserByID)

		// Company routes
		company := authGroup.Group("/company")
		{
			company.GET("/members", companyHandler.ListMembers)
			company.POST("/members", companyHandler.AddMember)
			company.PUT("/members", companyHandler.UpdateMember)
			company.DELETE("/members", companyHandler.DeleteMember)
			company.PUT("/members/role", companyHandler.ChangeRole)
		}

		// VPS routes
		vps := authGroup.Group("/vps")
		{
			vps.GET("", vpsHandler.List)
			vps.POST("", vpsHandler.Create)
			vps.GET("/detail", vpsHandler.GetByID)
			vps.POST("/start", vpsHandler.Start)
			vps.POST("/stop", vpsHandler.Stop)
			vps.POST("/restart", vpsHandler.Restart)
			vps.POST("/reinstall", vpsHandler.Reinstall)
			vps.DELETE("", vpsHandler.Delete)
			vps.GET("/console", vpsHandler.Console)
		}

		// Public console view route (uses signed token, no JWT required)
		api.GET("/vps/console/view", vpsHandler.ConsoleView)

		// Order routes
		orders := authGroup.Group("/orders")
		{
			orders.GET("", orderHandler.List)
			orders.GET("/detail", orderHandler.GetByID)
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
			tickets.GET("/detail", ticketHandler.GetByID)
			tickets.POST("/reply", ticketHandler.Reply)
			tickets.PUT("/close", ticketHandler.Close)
		}

		// API Key routes
		apikeys := authGroup.Group("/apikeys")
		{
			apikeys.GET("", apiKeyHandler.List)
			apikeys.POST("", apiKeyHandler.Create)
			apikeys.DELETE("", apiKeyHandler.Delete)
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
			admin.PUT("/users/status", adminHandler.UpdateUserStatus)
			admin.PUT("/users/password", adminHandler.ResetUserPassword)
			admin.GET("/companies", adminHandler.ListCompanies)
			admin.GET("/vps", adminHandler.ListAllVPS)
			admin.GET("/audit-logs", adminHandler.ListAllLogs)
			admin.POST("/reset", challengeHandler.ResetDatabase)
			admin.POST("/announcements", announceHandler.Create)
			admin.PUT("/announcements", announceHandler.Update)
			admin.DELETE("/announcements", announceHandler.Delete)
			admin.GET("/config", configHandler.GetConfig)
			admin.PUT("/config", configHandler.UpdateConfig)
		}

		// Challenge routes
		challenges := authGroup.Group("/challenges")
		{
			challenges.GET("", challengeHandler.List)
			challenges.GET("/detail", challengeHandler.Detail)
			challenges.GET("/hints", challengeHandler.GetHint)
			challenges.POST("/complete", challengeHandler.MarkComplete)
		}

		authGroup.GET("/security-mode", challengeHandler.GetSecurityMode)
		authGroup.PUT("/security-mode", challengeHandler.SetSecurityMode)

		// Encoding challenge state (any authenticated user can activate/deactivate)
		authGroup.GET("/encoding-challenge-state", challengeHandler.GetEncodingChallengeState)
		authGroup.PUT("/encoding-challenge-state", challengeHandler.SetEncodingChallengeState)

		// Encoded challenge routes (encoding/encryption wrapper endpoints)
		encoded := authGroup.Group("/encoded")
		{
			encoded.GET("/vps", encodedHandler.GetVPSByEncodedID)
			encoded.POST("/vps/start", encodedHandler.StartVPSEncoded)
			encoded.POST("/vps/stop", encodedHandler.StopVPSEncoded)
			encoded.POST("/vps/reinstall", encodedHandler.ReinstallVPSEncoded)
			encoded.GET("/users", encodedHandler.GetUserByEncodedID)
			encoded.GET("/orders", encodedHandler.GetOrderByEncodedID)
			encoded.GET("/tickets", encodedHandler.GetTicketByEncodedID)
			encoded.DELETE("/apikeys", encodedHandler.DeleteAPIKeyEncoded)
			encoded.POST("/company/members", encodedHandler.AddMemberEncoded)
			encoded.PUT("/company/members/role", encodedHandler.ChangeRoleEncoded)
		}

		// Crypto utility routes (for encoding/decoding practice)
		cryptoGroup := authGroup.Group("/crypto")
		{
			cryptoGroup.POST("/encode", encodedHandler.EncodeValue)
			cryptoGroup.POST("/decode", encodedHandler.DecodeValue)
			cryptoGroup.GET("/keys", encodedHandler.GetCryptoKeys)
		}
	}

	return r
}

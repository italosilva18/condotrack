package http

import (
	"log"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/internal/delivery/http/handler"
	"github.com/condotrack/api/internal/delivery/http/middleware"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/internal/infrastructure/database"
	"github.com/condotrack/api/internal/infrastructure/auth"
	"github.com/condotrack/api/internal/infrastructure/external/asaas"
	infraRepo "github.com/condotrack/api/internal/infrastructure/repository"
	"github.com/condotrack/api/internal/infrastructure/storage"
	"github.com/condotrack/api/internal/usecase/agenda"
	authUseCase "github.com/condotrack/api/internal/usecase/auth"
	"github.com/condotrack/api/internal/usecase/audit"
	"github.com/condotrack/api/internal/usecase/certificado"
	"github.com/condotrack/api/internal/usecase/checkout"
	"github.com/condotrack/api/internal/usecase/contrato"
	"github.com/condotrack/api/internal/usecase/course"
	"github.com/condotrack/api/internal/usecase/gestor"
	"github.com/condotrack/api/internal/usecase/inspection"
	"github.com/condotrack/api/internal/usecase/matricula"
	"github.com/condotrack/api/internal/usecase/payment"
	"github.com/condotrack/api/internal/usecase/revenue"
	"github.com/condotrack/api/internal/usecase/supplier"
	"github.com/condotrack/api/internal/usecase/task"
	"github.com/condotrack/api/internal/usecase/team"
	"github.com/gin-gonic/gin"
)

// Router holds all the handlers and configuration
type Router struct {
	cfg     *config.Config
	db      *database.MySQL
	storage *storage.StorageService

	// Handlers
	healthHandler         *handler.HealthHandler
	gestorHandler         *handler.GestorHandler
	contratoHandler       *handler.ContratoHandler
	auditHandler          *handler.AuditHandler
	auditCategoryHandler  *handler.AuditCategoryHandler
	matriculaHandler      *handler.MatriculaHandler
	paymentHandler        *handler.PaymentHandler
	checkoutHandler       *handler.CheckoutHandler
	webhookHandler        *handler.WebhookHandler
	certificadoHandler    *handler.CertificadoHandler
	imageHandler          *handler.ImageHandler
	portalHandler         *handler.PortalHandler
	notificationHandler   *handler.NotificationHandler
	statsHandler          *handler.StatsHandler
	revenueHandler        *handler.RevenueHandler
	supplierHandler       *handler.SupplierHandler
	courseHandler         *handler.CourseHandler
	taskHandler           *handler.TaskHandler
	teamHandler       *handler.TeamHandler
	agendaHandler     *handler.AgendaHandler
	inspectionHandler *handler.InspectionHandler
	authHandler       *handler.AuthHandler
	jwtManager        *auth.JWTManager
}

// NewRouter creates a new router with all dependencies
func NewRouter(cfg *config.Config, db *database.MySQL) *Router {
	// Initialize Asaas client
	asaasClient := asaas.NewClient(cfg.AsaasAPIKey, cfg.AsaasAPIURL)

	// Initialize storage service (MinIO)
	storageService, err := storage.NewStorageService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize storage service: %v", err)
		log.Printf("Portal features requiring storage will be unavailable")
	}

	// Initialize repositories
	gestorRepo := infraRepo.NewGestorMySQLRepository(db.DB)
	contratoRepo := infraRepo.NewContratoMySQLRepository(db.DB)
	auditRepo := infraRepo.NewAuditMySQLRepository(db.DB)
	auditItemRepo := infraRepo.NewAuditItemMySQLRepository(db.DB)
	auditCategoryRepo := infraRepo.NewAuditCategoryMySQLRepository(db.DB)
	matriculaRepo := infraRepo.NewMatriculaMySQLRepository(db.DB)
	certificadoRepo := infraRepo.NewCertificadoMySQLRepository(db.DB)
	notificacaoRepo := infraRepo.NewNotificacaoMySQLRepository(db.DB)
	revenueSplitRepo := infraRepo.NewRevenueSplitMySQLRepository(db.DB)
	supplierRepo := infraRepo.NewSupplierMySQLRepository(db.DB)
	courseRepo := infraRepo.NewCourseMySQLRepository(db.DB)
	taskRepo := infraRepo.NewTaskMySQLRepository(db.DB)
	teamRepo := infraRepo.NewTeamMySQLRepository(db.DB)
	agendaRepo := infraRepo.NewAgendaMySQLRepository(db.DB)
	inspectionRepo := infraRepo.NewInspectionMySQLRepository(db.DB)
	userRepo := infraRepo.NewUserMySQLRepository(db.DB)

	// Initialize use cases
	gestorUC := gestor.NewUseCase(gestorRepo)
	contratoUC := contrato.NewUseCase(contratoRepo, gestorRepo)
	auditUC := audit.NewUseCase(auditRepo, auditItemRepo, contratoRepo, db)
	auditCategoryUC := audit.NewCategoryUseCase(auditCategoryRepo)
	matriculaUC := matricula.NewUseCase(matriculaRepo)
	paymentUC := payment.NewUseCase(asaasClient, cfg)
	checkoutUC := checkout.NewUseCase(asaasClient, matriculaRepo, db, cfg)
	certificadoUC := certificado.NewUseCase(certificadoRepo, matriculaRepo)
	revenueUC := revenue.NewUseCase(revenueSplitRepo)
	supplierUC := supplier.NewUseCase(supplierRepo)
	courseUC := course.NewUseCase(courseRepo)
	taskUC := task.NewUseCase(taskRepo, contratoRepo, gestorRepo)
	teamUC := team.NewUseCase(teamRepo, gestorRepo, contratoRepo)
	agendaUC := agenda.NewUseCase(agendaRepo, contratoRepo, gestorRepo)
	inspectionUC := inspection.NewUseCase(inspectionRepo, contratoRepo, gestorRepo)

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
	authUC := authUseCase.NewUseCase(userRepo, jwtManager)

	// Initialize handlers
	return &Router{
		cfg:                  cfg,
		db:                   db,
		storage:              storageService,
		healthHandler:        handler.NewHealthHandler(db),
		gestorHandler:        handler.NewGestorHandler(gestorUC),
		contratoHandler:      handler.NewContratoHandler(contratoUC),
		auditHandler:         handler.NewAuditHandler(auditUC),
		auditCategoryHandler: handler.NewAuditCategoryHandler(auditCategoryUC),
		matriculaHandler:     handler.NewMatriculaHandler(matriculaUC),
		paymentHandler:       handler.NewPaymentHandler(paymentUC, matriculaRepo),
		checkoutHandler:      handler.NewCheckoutHandler(checkoutUC),
		webhookHandler:       handler.NewWebhookHandler(cfg, db, matriculaRepo),
		certificadoHandler:   handler.NewCertificadoHandler(certificadoUC),
		imageHandler:         handler.NewImageHandler(cfg),
		portalHandler:        handler.NewPortalHandler(storageService, cfg),
		notificationHandler:  handler.NewNotificationHandler(notificacaoRepo),
		statsHandler:         handler.NewStatsHandler(db.DB, matriculaRepo, auditRepo, contratoRepo, gestorRepo),
		revenueHandler:       handler.NewRevenueHandler(revenueUC),
		supplierHandler:      handler.NewSupplierHandler(supplierUC),
		courseHandler:        handler.NewCourseHandler(courseUC),
		taskHandler:          handler.NewTaskHandler(taskUC),
		teamHandler:       handler.NewTeamHandler(teamUC),
		agendaHandler:     handler.NewAgendaHandler(agendaUC),
		inspectionHandler: handler.NewInspectionHandler(inspectionUC),
		authHandler:       handler.NewAuthHandler(authUC),
		jwtManager:        jwtManager,
	}
}

// Setup configures the Gin router with all routes
func (r *Router) Setup() *gin.Engine {
	// Set Gin mode
	if r.cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine
	engine := gin.New()

	// Apply global middlewares
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// Serve static files (uploads)
	engine.Static("/uploads", r.cfg.UploadDir)

	// Health check routes
	engine.GET("/ping", r.healthHandler.Ping)

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Health
		v1.GET("/health", r.healthHandler.HealthCheck)

		// Gestores
		gestores := v1.Group("/gestores")
		{
			gestores.GET("", r.gestorHandler.ListGestores)
			gestores.GET("/:id", r.gestorHandler.GetGestorByID)
			gestores.POST("", r.gestorHandler.CreateGestor)
			gestores.PUT("/:id", r.gestorHandler.UpdateGestor)
			gestores.DELETE("/:id", r.gestorHandler.DeleteGestor)
		}

		// Contratos
		contratos := v1.Group("/contratos")
		{
			contratos.GET("", r.contratoHandler.ListContratos)
			contratos.GET("/:id", r.contratoHandler.GetContratoByID)
			contratos.POST("", r.contratoHandler.CreateContrato)
			contratos.PUT("/:id", r.contratoHandler.UpdateContrato)
			contratos.DELETE("/:id", r.contratoHandler.DeleteContrato)
		}

		// Audits
		audits := v1.Group("/audits")
		{
			audits.GET("", r.auditHandler.ListAudits)
			audits.GET("/meta", r.auditHandler.GetAuditMeta)
			audits.GET("/:id", r.auditHandler.GetAuditByID)
			audits.POST("", r.auditHandler.CreateAudit)
			audits.PUT("/:id", r.auditHandler.UpdateAudit)
			audits.DELETE("/:id", r.auditHandler.DeleteAudit)
		}

		// Audit Categories
		auditCategories := v1.Group("/audit-categories")
		{
			auditCategories.GET("", r.auditCategoryHandler.ListCategories)
			auditCategories.GET("/:id", r.auditCategoryHandler.GetCategoryByID)
			auditCategories.POST("", r.auditCategoryHandler.CreateCategory)
			auditCategories.PUT("/:id", r.auditCategoryHandler.UpdateCategory)
			auditCategories.DELETE("/:id", r.auditCategoryHandler.DeleteCategory)
		}

		// Enrollments
		enrollments := v1.Group("/enrollments")
		{
			enrollments.GET("", r.matriculaHandler.ListEnrollments)
			enrollments.GET("/:id", r.matriculaHandler.GetEnrollmentByID)
			enrollments.POST("", r.matriculaHandler.CreateEnrollment)
			enrollments.PATCH("/:id/payment-status", r.matriculaHandler.UpdatePaymentStatus)
			enrollments.PATCH("/:id/progress", r.matriculaHandler.UpdateProgress)
		}

		// Payments
		payments := v1.Group("/payments")
		{
			payments.GET("", r.paymentHandler.ListPayments)
			payments.GET("/enrollment/:id", r.paymentHandler.GetPaymentsByEnrollment)
			payments.POST("/customer", r.paymentHandler.CreateCustomer)
			payments.POST("/pix", r.paymentHandler.CreatePixPayment)
			payments.POST("/boleto", r.paymentHandler.CreateBoletoPayment)
			payments.POST("/card", r.paymentHandler.CreateCardPayment)
			payments.GET("/:id/status", r.paymentHandler.GetPaymentStatus)
			payments.GET("/simulate-split", r.paymentHandler.SimulateRevenueSplit)
		}

		// Checkout
		checkout := v1.Group("/checkout")
		{
			checkout.POST("", r.checkoutHandler.CreateCheckout)
			checkout.GET("/:id/status", r.checkoutHandler.GetCheckoutStatus)
		}

		// Webhooks
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/asaas", r.webhookHandler.HandleAsaasWebhook)
		}

		// Certificados
		certificados := v1.Group("/certificados")
		{
			certificados.GET("/:aluno_id", r.certificadoHandler.GetCertificatesByStudent)
			certificados.GET("/detail/:id", r.certificadoHandler.GetCertificateByID)
			certificados.GET("/validate/:code", r.certificadoHandler.ValidateCertificate)
			certificados.POST("/generate", r.certificadoHandler.GenerateCertificate)
		}

		// Notifications
		notifications := v1.Group("/notifications")
		{
			notifications.GET("", r.notificationHandler.ListNotifications)
			notifications.GET("/unread", r.notificationHandler.GetUnreadNotifications)
			notifications.GET("/count", r.notificationHandler.GetUnreadCount)
			notifications.POST("", r.notificationHandler.CreateNotification)
			notifications.PATCH("/:id/read", r.notificationHandler.MarkAsRead)
			notifications.PATCH("/mark-all-read", r.notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", r.notificationHandler.DeleteNotification)
		}

		// Images
		images := v1.Group("/images")
		{
			images.GET("", r.imageHandler.ListImages)
			images.POST("", r.imageHandler.UploadImage)
			images.DELETE("/:filename", r.imageHandler.DeleteImage)
		}

		// Stats/Dashboard
		stats := v1.Group("/stats")
		{
			stats.GET("/overview", r.statsHandler.GetOverview)
			stats.GET("/enrollments", r.statsHandler.GetEnrollmentStats)
			stats.GET("/payments", r.statsHandler.GetPaymentStats)
			stats.GET("/audits", r.statsHandler.GetAuditStats)
		}

		// Revenue Splits
		revenueSplits := v1.Group("/revenue-splits")
		{
			revenueSplits.GET("", r.revenueHandler.ListRevenueSplits)
			revenueSplits.GET("/:id", r.revenueHandler.GetRevenueSplitByID)
			revenueSplits.GET("/enrollment/:id", r.revenueHandler.GetRevenueSplitByEnrollment)
			revenueSplits.GET("/instructor/:id", r.revenueHandler.GetInstructorEarnings)
			revenueSplits.GET("/instructor/:id/total", r.revenueHandler.GetInstructorTotalEarnings)
			revenueSplits.PATCH("/:id/status", r.revenueHandler.UpdateStatus)
		}

		// Suppliers
		suppliers := v1.Group("/suppliers")
		{
			suppliers.GET("", r.supplierHandler.ListSuppliers)
			suppliers.GET("/:id", r.supplierHandler.GetSupplierByID)
			suppliers.POST("", r.supplierHandler.CreateSupplier)
			suppliers.PUT("/:id", r.supplierHandler.UpdateSupplier)
			suppliers.DELETE("/:id", r.supplierHandler.DeleteSupplier)
		}

		// Courses
		courses := v1.Group("/courses")
		{
			courses.GET("", r.courseHandler.ListCourses)
			courses.GET("/:id", r.courseHandler.GetCourseByID)
			courses.POST("", r.courseHandler.CreateCourse)
			courses.PUT("/:id", r.courseHandler.UpdateCourse)
			courses.DELETE("/:id", r.courseHandler.DeleteCourse)
		}

		// Tasks
		tasks := v1.Group("/tasks")
		{
			tasks.GET("", r.taskHandler.ListTasks)
			tasks.GET("/overdue", r.taskHandler.GetOverdueTasks)
			tasks.GET("/contract/:id", r.taskHandler.GetTasksByContract)
			tasks.GET("/assignee/:id", r.taskHandler.GetTasksByAssignee)
			tasks.GET("/:id", r.taskHandler.GetTaskByID)
			tasks.POST("", r.taskHandler.CreateTask)
			tasks.PUT("/:id", r.taskHandler.UpdateTask)
			tasks.PATCH("/:id/status", r.taskHandler.UpdateTaskStatus)
			tasks.DELETE("/:id", r.taskHandler.DeleteTask)
		}

		// Team Management
		teamGroup := v1.Group("/team")
		{
			teamGroup.GET("", r.teamHandler.ListTeamMembers)
			teamGroup.GET("/:id", r.teamHandler.GetTeamMemberByID)
			teamGroup.POST("", r.teamHandler.CreateTeamMember)
			teamGroup.PUT("/:id", r.teamHandler.UpdateTeamMember)
			teamGroup.DELETE("/:id", r.teamHandler.DeleteTeamMember)
			teamGroup.GET("/contract/:id", r.teamHandler.GetTeamByContract)
			teamGroup.GET("/user/:id", r.teamHandler.GetContractsByUser)
		}

		// Agenda (Calendar)
		agendaRoutes := v1.Group("/agenda")
		{
			agendaRoutes.GET("", r.agendaHandler.ListEvents)
			agendaRoutes.GET("/:id", r.agendaHandler.GetEventByID)
			agendaRoutes.POST("", r.agendaHandler.CreateEvent)
			agendaRoutes.PUT("/:id", r.agendaHandler.UpdateEvent)
			agendaRoutes.DELETE("/:id", r.agendaHandler.DeleteEvent)
		}

		// Inspections
		inspections := v1.Group("/inspections")
		{
			inspections.GET("", r.inspectionHandler.ListInspections)
			inspections.GET("/scheduled", r.inspectionHandler.GetScheduledInspections)
			inspections.GET("/:id", r.inspectionHandler.GetInspectionByID)
			inspections.POST("", r.inspectionHandler.CreateInspection)
			inspections.PUT("/:id", r.inspectionHandler.UpdateInspection)
			inspections.DELETE("/:id", r.inspectionHandler.DeleteInspection)
		}

		// Authentication
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", r.authHandler.Login)
			authRoutes.POST("/register", r.authHandler.Register)
			authRoutes.POST("/logout", r.authHandler.Logout)

			// Protected routes
			authProtected := authRoutes.Group("")
			authProtected.Use(middleware.AuthMiddleware(r.jwtManager))
			{
				authProtected.GET("/me", r.authHandler.GetCurrentUser)
				authProtected.PUT("/me", r.authHandler.UpdateUser)
				authProtected.POST("/change-password", r.authHandler.ChangePassword)
			}

			// Admin routes
			adminRoutes := authRoutes.Group("/users")
			adminRoutes.Use(middleware.AuthMiddleware(r.jwtManager))
			adminRoutes.Use(middleware.RequireRole("admin"))
			{
				adminRoutes.GET("", r.authHandler.ListUsers)
				adminRoutes.GET("/:id", r.authHandler.GetUserByID)
				adminRoutes.PUT("/:id", r.authHandler.AdminUpdateUser)
				adminRoutes.DELETE("/:id", r.authHandler.DeleteUser)
			}
		}

		// Portal-specific endpoints
		portal := v1.Group("/portal")
		{
			// Portal images (with system overwrite support)
			portal.GET("/images", r.portalHandler.ListPortalImages)
			portal.POST("/images", r.portalHandler.UploadPortalImage)
			portal.DELETE("/images/:filename", r.portalHandler.DeletePortalImage)

			// Evidence uploads for audits
			portal.GET("/evidence", r.portalHandler.ListEvidence)
			portal.POST("/evidence", r.portalHandler.UploadEvidence)
			portal.DELETE("/evidence/:filename", r.portalHandler.DeleteEvidence)

			// AI Proxy (Gemini)
			portal.POST("/ai", r.portalHandler.ProxyGeminiAI)
		}

		// Backend integration compatibility routes (for legacy PHP API compatibility)
		backendIntegration := engine.Group("/backend_integration")
		{
			backendIntegration.GET("/api_router.php", r.handleLegacyAPIRouter)
			backendIntegration.POST("/api_router.php", r.handleLegacyAPIRouter)
			backendIntegration.DELETE("/api_router.php", r.handleLegacyAPIRouter)
		}
	}

	return engine
}

// handleLegacyAPIRouter handles legacy PHP API compatibility
func (r *Router) handleLegacyAPIRouter(c *gin.Context) {
	endpoint := c.Query("endpoint")

	switch endpoint {
	case "login":
		r.authHandler.Login(c)
	case "register":
		r.authHandler.Register(c)
	case "me", "current_user":
		// Check if token is provided
		r.authHandler.GetCurrentUser(c)
	case "logout":
		r.authHandler.Logout(c)
	case "images":
		if c.Request.Method == "GET" {
			r.portalHandler.ListPortalImages(c)
		} else if c.Request.Method == "DELETE" {
			// Extract filename from query
			filename := c.Query("file")
			c.Params = append(c.Params, gin.Param{Key: "filename", Value: filename})
			r.portalHandler.DeletePortalImage(c)
		}
	case "upload":
		r.portalHandler.UploadPortalImage(c)
	case "health", "ping":
		r.healthHandler.HealthCheck(c)
	case "gestores", "managers":
		r.gestorHandler.ListGestores(c)
	case "contratos", "contracts":
		if c.Request.Method == "GET" {
			r.contratoHandler.ListContratos(c)
		} else if c.Request.Method == "POST" {
			r.contratoHandler.CreateContrato(c)
		}
	case "audits", "auditorias":
		if c.Request.Method == "GET" {
			r.auditHandler.ListAudits(c)
		} else if c.Request.Method == "POST" {
			r.auditHandler.CreateAudit(c)
		}
	case "enrollments", "matriculas":
		if c.Request.Method == "GET" {
			r.matriculaHandler.ListEnrollments(c)
		} else if c.Request.Method == "POST" {
			r.matriculaHandler.CreateEnrollment(c)
		}
	case "courses":
		if c.Request.Method == "GET" {
			r.courseHandler.ListCourses(c)
		} else if c.Request.Method == "POST" {
			r.courseHandler.CreateCourse(c)
		}
	case "tasks":
		if c.Request.Method == "GET" {
			r.taskHandler.ListTasks(c)
		} else if c.Request.Method == "POST" {
			r.taskHandler.CreateTask(c)
		}
	case "suppliers":
		if c.Request.Method == "GET" {
			r.supplierHandler.ListSuppliers(c)
		} else if c.Request.Method == "POST" {
			r.supplierHandler.CreateSupplier(c)
		}
	case "team":
		if c.Request.Method == "GET" {
			r.teamHandler.ListTeamMembers(c)
		} else if c.Request.Method == "POST" {
			r.teamHandler.CreateTeamMember(c)
		}
	case "agenda":
		if c.Request.Method == "GET" {
			r.agendaHandler.ListEvents(c)
		} else if c.Request.Method == "POST" {
			r.agendaHandler.CreateEvent(c)
		}
	case "inspections":
		if c.Request.Method == "GET" {
			r.inspectionHandler.ListInspections(c)
		} else if c.Request.Method == "POST" {
			r.inspectionHandler.CreateInspection(c)
		}
	case "stats":
		r.statsHandler.GetOverview(c)
	case "notifications":
		if c.Request.Method == "GET" {
			r.notificationHandler.ListNotifications(c)
		} else if c.Request.Method == "POST" {
			r.notificationHandler.CreateNotification(c)
		}
	default:
		c.JSON(200, gin.H{
			"status":  "CondoTrack API Online (Go)",
			"actions": []string{"login", "register", "me", "images", "upload", "health", "gestores", "contratos", "audits", "enrollments", "courses", "tasks", "suppliers", "team", "agenda", "inspections", "stats", "notifications"},
		})
	}
}

// GetMatriculaRepository returns the matricula repository (for webhook handler)
func (r *Router) GetMatriculaRepository() repository.MatriculaRepository {
	return infraRepo.NewMatriculaMySQLRepository(r.db.DB)
}

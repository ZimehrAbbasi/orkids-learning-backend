package api

import (
	"orkidslearning/src/router"
	"orkidslearning/src/services"

	"github.com/gin-gonic/gin"
)

// InitializeRoutes sets up all application routes
func InitializeRoutes(router *gin.Engine, contextService *services.ContextService) {

	router.Use(LoggerMiddleware())

	// Public routes
	public := router.Group("/api/public")
	public.Use(InjectContextService(contextService))
	initializePublicRoutes(public)

	// Auth routes
	auth := router.Group("/api/auth")
	auth.Use(InjectContextService(contextService))
	initializeAuthRoutes(auth)

	// Protected routes with JWT middleware
	protected := router.Group("api")
	protected.Use(JWTAuthMiddleware(contextService.GetJWTService()))
	protected.Use(InjectContextService(contextService))
	initializeProtectedRoutes(protected)
}

// initializePublicRoutes defines public routes
func initializePublicRoutes(public *gin.RouterGroup) {
	public.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Welcome to the Gin server with MongoDB!"})
	})
	public.GET("/courses", router.GetAllCourses)
}

// initializeAuthRoutes defines authentication routes
func initializeAuthRoutes(auth *gin.RouterGroup) {
	auth.POST("/signup", router.SignupHandler)
	auth.POST("/login", router.LoginHandler)
}

// initializeProtectedRoutes defines protected routes
func initializeProtectedRoutes(protected *gin.RouterGroup) {
	protected.POST("/courses", router.AddCourse)
	protected.POST("/courses/:id", router.GetCourseById)
	protected.POST("/courses/enroll/:id", router.EnrollInCourse)
	protected.POST("/courses/unenroll/:id", router.UnenrollFromCourse)
}

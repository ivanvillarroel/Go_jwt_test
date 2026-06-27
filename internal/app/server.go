package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ivanvillarroel/go_jwt_gin/internal/auth"
	"github.com/ivanvillarroel/go_jwt_gin/internal/config"
	"github.com/ivanvillarroel/go_jwt_gin/internal/database"
	"github.com/ivanvillarroel/go_jwt_gin/internal/users"
)

type Server struct {
	Engine *gin.Engine
	DB     *sql.DB
}

func NewServer(cfg config.Config) (*Server, error) {
	db, err := database.NewInMemory()
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(db); err != nil {
		return nil, err
	}

	if err := database.SeedUsers(db); err != nil {
		return nil, err
	}

	engine, err := NewRouter(db, cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		Engine: engine,
		DB:     db,
	}, nil
}

func NewRouter(db *sql.DB, cfg config.Config) (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	userRepository := users.NewSQLiteRepository(db)
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)

	tokenService := auth.NewJWTService(cfg.JWTSecret)
	authHandler := auth.NewHandler(tokenService, cfg.AuthUser, cfg.AuthPassword)
	authMiddleware := auth.Middleware(tokenService)

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	api.POST("/auth/token", authHandler.Login)

	protected := api.Group("")
	protected.Use(authMiddleware)
	protected.GET("/auth/valid", authHandler.Valid)
	protected.GET("/auth/read", authHandler.Read)
	protected.GET("/users", auth.RequirePermission("users:read"), userHandler.List)

	return router, nil
}

func (s *Server) Run(address string) error {
	return s.Engine.Run(address)
}

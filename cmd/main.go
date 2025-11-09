package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/config"
	"github.com/rafaelc-rb/geekery-api/internal/database"
	"github.com/rafaelc-rb/geekery-api/internal/logger"
	"github.com/rafaelc-rb/geekery-api/internal/middleware"
	"github.com/rafaelc-rb/geekery-api/internal/routes"
)

func main() {
	// Banner
	printBanner()

	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Inicializar logger estruturado
	logger.Init(cfg.LogLevel)
	logger.Info().Msg("Configuration loaded successfully")

	// Conectar ao banco de dados
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Configurar modo do Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Criar router (sem logger padrão do Gin)
	router := gin.New()
	router.Use(gin.Recovery()) // Manter recovery middleware

	// Middlewares customizados
	router.Use(middleware.Logger()) // Logger estruturado
	router.Use(corsMiddleware())    // CORS

	// Configurar rotas
	routes.SetupRoutes(router, db)

	// Iniciar servidor
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Info().
		Str("address", fmt.Sprintf("http://localhost%s", serverAddr)).
		Str("environment", cfg.Environment).
		Msg("Server starting")

	// Graceful shutdown
	go func() {
		if err := router.Run(serverAddr); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Esperar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")
	logger.Info().Msg("Server stopped gracefully")
}

// corsMiddleware adiciona headers CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// printBanner exibe o banner da aplicação
func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════╗
║                                                   ║
║     ██████╗ ███████╗███████╗██╗  ██╗███████╗██████╗██╗   ██╗    ║
║    ██╔════╝ ██╔════╝██╔════╝██║ ██╔╝██╔════╝██╔══██╗╚██╗ ██╔╝    ║
║    ██║  ███╗█████╗  █████╗  █████╔╝ █████╗  ██████╔╝ ╚████╔╝     ║
║    ██║   ██║██╔══╝  ██╔══╝  ██╔═██╗ ██╔══╝  ██╔══██╗  ╚██╔╝      ║
║    ╚██████╔╝███████╗███████╗██║  ██╗███████╗██║  ██║   ██║       ║
║     ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝       ║
║                                                   ║
║          Your Personal Geek Media Tracker        ║
║                    API v1.0.0                    ║
║                                                   ║
╚═══════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

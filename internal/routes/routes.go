package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/auth"
	"github.com/rafaelc-rb/geekery-api/internal/config"
	"github.com/rafaelc-rb/geekery-api/internal/handlers"
	"github.com/rafaelc-rb/geekery-api/internal/repositories"
	"github.com/rafaelc-rb/geekery-api/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Geekery API is running!",
		})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	// ========================================
	// JWT Manager
	// ========================================
	cfg := config.AppConfig
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, 24*time.Hour) // Token válido por 24 horas

	// ========================================
	// Repositórios
	// ========================================
	itemRepo := repositories.NewItemRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	userItemRepo := repositories.NewUserItemRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// ========================================
	// Serviços
	// ========================================
	itemService := services.NewItemService(itemRepo, tagRepo)
	tagService := services.NewTagService(tagRepo)
	userItemService := services.NewUserItemService(userItemRepo, itemRepo)
	authService := services.NewAuthService(userRepo, jwtManager)

	// ========================================
	// Handlers
	// ========================================
	itemHandler := handlers.NewItemHandler(itemService)
	tagHandler := handlers.NewTagHandler(tagService)
	userItemHandler := handlers.NewUserItemHandler(userItemService)
	authHandler := handlers.NewAuthHandler(authService)

	// ========================================
	// Rotas Públicas - Catálogo de Items
	// ========================================
	itemsRoutes := api.Group("/items")
	{
		itemsRoutes.GET("", itemHandler.GetAllItems)           // GET /api/items?type=anime
		itemsRoutes.GET("/search", itemHandler.SearchItems)    // GET /api/items/search?q=attack
		itemsRoutes.GET("/:id", itemHandler.GetItemByID)       // GET /api/items/1

		// Rotas admin (futuramente protegidas)
		itemsRoutes.POST("", itemHandler.CreateItem)           // POST /api/items
		itemsRoutes.PUT("/:id", itemHandler.UpdateItem)        // PUT /api/items/1
		itemsRoutes.DELETE("/:id", itemHandler.DeleteItem)     // DELETE /api/items/1

		// Import endpoints
		itemsRoutes.POST("/import/anime", itemHandler.ImportAnime)     // POST /api/items/import/anime
		itemsRoutes.POST("/import/comic", itemHandler.ImportComic)     // POST /api/items/import/comic
		itemsRoutes.POST("/import/novel", itemHandler.ImportNovel)     // POST /api/items/import/novel
		itemsRoutes.POST("/import/movie", itemHandler.ImportMovie)     // POST /api/items/import/movie
		itemsRoutes.POST("/import/series", itemHandler.ImportSeries)   // POST /api/items/import/series
		itemsRoutes.POST("/import/game", itemHandler.ImportGame)       // POST /api/items/import/game
		itemsRoutes.POST("/import/book", itemHandler.ImportBook)       // POST /api/items/import/book
	}

	// ========================================
	// Rotas de Tags
	// ========================================
	tagsRoutes := api.Group("/tags")
	{
		tagsRoutes.POST("", tagHandler.CreateTag)          // POST /api/tags
		tagsRoutes.GET("", tagHandler.GetAllTags)          // GET /api/tags
		tagsRoutes.GET("/:id", tagHandler.GetTagByID)      // GET /api/tags/1
		tagsRoutes.PUT("/:id", tagHandler.UpdateTag)       // PUT /api/tags/1
		tagsRoutes.DELETE("/:id", tagHandler.DeleteTag)    // DELETE /api/tags/1
	}

	// ========================================
	// Rotas de Autenticação (Públicas)
	// ========================================
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register) // POST /api/auth/register
		authRoutes.POST("/login", authHandler.Login)       // POST /api/auth/login
	}

	// ========================================
	// Rotas Protegidas - Lista Pessoal do Usuário
	// Requer autenticação JWT
	// ========================================
	myListRoutes := api.Group("/my-list")
	myListRoutes.Use(auth.AuthMiddleware(jwtManager)) // Proteger todas as rotas deste grupo
	{
		myListRoutes.POST("", userItemHandler.AddToList)              // POST /api/my-list
		myListRoutes.GET("", userItemHandler.GetMyList)               // GET /api/my-list?status=watching&favorite=true
		myListRoutes.GET("/stats", userItemHandler.GetStatistics)     // GET /api/my-list/stats
		myListRoutes.GET("/:id", userItemHandler.GetMyListItem)       // GET /api/my-list/1
		myListRoutes.PUT("/:id", userItemHandler.UpdateListItem)      // PUT /api/my-list/1
		myListRoutes.DELETE("/:id", userItemHandler.RemoveFromList)   // DELETE /api/my-list/1
	}
}

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafaelc-rb/geekery-api/internal/logger"
)

// Logger middleware para log estruturado de requisições HTTP
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capturar tempo de início
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Processar requisição
		c.Next()

		// Calcular latência
		latency := time.Since(start)

		// Log estruturado da requisição
		logEvent := logger.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Int("size", c.Writer.Size())

		// Adicionar query se existir
		if query != "" {
			logEvent.Str("query", query)
		}

		// Adicionar errors se existir
		if len(c.Errors) > 0 {
			logEvent.Str("errors", c.Errors.String())
		}

		logEvent.Msg("HTTP request")
	}
}

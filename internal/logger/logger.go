package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Log é a instância global do logger estruturado
var Log zerolog.Logger

// Init inicializa o logger com o nível configurado
func Init(level string) {
	// Formato de timestamp Unix para performance
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Configurar nível de log baseado no ambiente
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Criar logger com timestamp
	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()

	Log.Info().
		Str("level", level).
		Msg("Logger initialized")
}

// Debug log com nível debug
func Debug() *zerolog.Event {
	return Log.Debug()
}

// Info log com nível info
func Info() *zerolog.Event {
	return Log.Info()
}

// Warn log com nível warn
func Warn() *zerolog.Event {
	return Log.Warn()
}

// Error log com nível error
func Error() *zerolog.Event {
	return Log.Error()
}

// Fatal log com nível fatal (termina o programa)
func Fatal() *zerolog.Event {
	return Log.Fatal()
}

package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger para exponer una interfaz limpia y contextualizable.
type Logger struct {
	zap *zap.Logger
}

// Config define la configuración del logger.
type Config struct {
	// Level: debug | info | warn | error
	Level string
	// Env: development | production
	Env string
	// ServiceName se incluye en cada log como campo fijo.
	ServiceName string
}

// New crea una instancia del logger según el entorno.
// En development: salida legible en consola con colores.
// En production:  salida JSON estructurada para ingestión por herramientas como Datadog, Loki, etc.
func New(cfg Config) *Logger {
	level := parseLevel(cfg.Level)

	var zapLogger *zap.Logger

	if strings.ToLower(cfg.Env) == "production" {
		zapLogger = buildProductionLogger(level, cfg.ServiceName)
	} else {
		zapLogger = buildDevelopmentLogger(level, cfg.ServiceName)
	}

	return &Logger{zap: zapLogger}
}

// -------------------------
// Métodos de logging
// -------------------------

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

// With retorna un nuevo Logger con campos fijos adicionales.
// Ideal para propagar contexto (request_id, user_id, trace_id, etc.)
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{zap: l.zap.With(fields...)}
}

// Sync flushea los buffers del logger. Llamar en defer en el main.
func (l *Logger) Sync() {
	_ = l.zap.Sync()
}

// Zap expone el logger subyacente para integraciones que lo requieran directamente.
func (l *Logger) Zap() *zap.Logger {
	return l.zap
}

// -------------------------
// Builders internos
// -------------------------

func buildProductionLogger(level zapcore.Level, serviceName string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0)).
		With(zap.String("service", serviceName))
}

func buildDevelopmentLogger(level zapcore.Level, serviceName string) *zap.Logger {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05.000"))
	}
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)

	return zap.New(core, zap.AddCaller(), zap.Development()).
		With(zap.String("service", serviceName))
}

// -------------------------
// Helpers
// -------------------------

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

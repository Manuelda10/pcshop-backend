package logger

import (
	"go.uber.org/zap"
)

// Este archivo centraliza los campos de log más comunes para garantizar
// consistencia de nombres a través de todos los servicios.
// Uso: logger.RequestID("abc-123"), logger.UserID("uuid")

// ---- Request / HTTP ----

func RequestID(id string) zap.Field {
	return zap.String("request_id", id)
}

func Method(method string) zap.Field {
	return zap.String("http_method", method)
}

func Path(path string) zap.Field {
	return zap.String("http_path", path)
}

func StatusCode(code int) zap.Field {
	return zap.Int("status_code", code)
}

func Latency(ms int64) zap.Field {
	return zap.Int64("latency_ms", ms)
}

func IP(ip string) zap.Field {
	return zap.String("client_ip", ip)
}

func UserAgent(ua string) zap.Field {
	return zap.String("user_agent", ua)
}

// ---- Dominio ----

func UserID(id string) zap.Field {
	return zap.String("user_id", id)
}

func Email(email string) zap.Field {
	return zap.String("email", email)
}

// ---- Capas de aplicación ----

// Layer identifica desde qué capa se emite el log: handler | service | repository
func Layer(layer string) zap.Field {
	return zap.String("layer", layer)
}

// Operation describe la operación en curso: "register_user", "login", etc.
func Operation(op string) zap.Field {
	return zap.String("operation", op)
}

// Input loguea lo que recibe una función/capa (evitar datos sensibles).
func Input(data any) zap.Field {
	return zap.Any("input", data)
}

// Output loguea lo que retorna una función/capa.
func Output(data any) zap.Field {
	return zap.Any("output", data)
}

// ---- Errores ----

func Err(err error) zap.Field {
	return zap.Error(err)
}

func ErrDetail(detail string) zap.Field {
	return zap.String("error_detail", detail)
}

// ---- Base de datos ----

func DBQuery(query string) zap.Field {
	return zap.String("db_query", query)
}

func DBTable(table string) zap.Field {
	return zap.String("db_table", table)
}

// ---- Infraestructura ----

func TraceID(id string) zap.Field {
	return zap.String("trace_id", id)
}

func Component(name string) zap.Field {
	return zap.String("component", name)
}

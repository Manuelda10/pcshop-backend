package logger_test

// Este archivo muestra cómo usar el logger en cada capa de la arquitectura hexagonal.
// No es código de producción, es documentación ejecutable.

// ============================================================
// 1. INICIALIZACIÓN EN main.go DE CADA MICROSERVICIO
// ============================================================
//
// import "github.com/tu-org/pc-store/shared/logger"
//
// func main() {
//     log := logger.New(logger.Config{
//         Level:       os.Getenv("LOG_LEVEL"),   // "debug" en dev, "info" en prod
//         Env:         os.Getenv("APP_ENV"),      // "development" | "production"
//         ServiceName: "auth-service",
//     })
//     defer log.Sync()
//
//     app := fiber.New()
//     app.Use(logger.FiberMiddleware(log))  // middleware global
//
//     log.Info("server starting", zap.Int("port", 8080))
//     app.Listen(":8080")
// }

// ============================================================
// 2. EN UN HANDLER (adapter/input/http)
// ============================================================
//
// func (h *AuthHandler) Register(c *fiber.Ctx) error {
//     log := logger.FromCtx(c)  // logger ya tiene request_id, method, path, ip
//
//     var req RegisterRequest
//     if err := c.BodyParser(&req); err != nil {
//         log.Warn("invalid request body",
//             logger.Layer("handler"),
//             logger.Operation("register"),
//             logger.Err(err),
//         )
//         return c.Status(400).JSON(...)
//     }
//
//     log.Debug("handler received input",
//         logger.Layer("handler"),
//         logger.Operation("register"),
//         logger.Input(req),  // password ya fue sanitizado por el middleware
//     )
//
//     result, err := h.service.Register(c.Context(), req.toCommand(), log)
//     if err != nil { ... }
//
//     log.Info("handler sending response",
//         logger.Layer("handler"),
//         logger.Operation("register"),
//         logger.Output(result),
//     )
//     return c.Status(201).JSON(result)
// }

// ============================================================
// 3. EN UN SERVICE (core/services)
// ============================================================
//
// func (s *AuthService) Register(ctx context.Context, cmd RegisterCommand, log *logger.Logger) (*User, error) {
//     log.Debug("service processing register",
//         logger.Layer("service"),
//         logger.Operation("register"),
//         logger.Input(cmd),
//     )
//
//     // ... lógica de negocio ...
//
//     user, err := s.userRepo.Save(ctx, newUser)
//     if err != nil {
//         log.Error("failed to save user",
//             logger.Layer("service"),
//             logger.Operation("register"),
//             logger.DBTable("users"),
//             logger.Err(err),
//         )
//         return nil, ErrInternalServer
//     }
//
//     log.Info("user registered successfully",
//         logger.Layer("service"),
//         logger.Operation("register"),
//         logger.UserID(user.ID.String()),
//         logger.Email(user.Email),
//     )
//     return user, nil
// }

// ============================================================
// 4. EN UN REPOSITORY (adapter/output/postgres)
// ============================================================
//
// func (r *UserRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
//     r.log.Debug("repository executing query",
//         logger.Layer("repository"),
//         logger.Operation("save_user"),
//         logger.DBTable("users"),
//         logger.DBQuery("INSERT INTO users ..."),
//     )
//     // ... pgx query ...
// }

// ============================================================
// EJEMPLO DE OUTPUT EN DEVELOPMENT (consola)
// ============================================================
//
// 10:32:01.123  DEBUG  auth-service  → request received    {request_id: "abc-123", method: "POST", path: "/auth/register", body: {"email":"user@test.com","password":"***"}}
// 10:32:01.124  DEBUG  auth-service  handler received input {layer: "handler", operation: "register", input: {email: "user@test.com"}}
// 10:32:01.125  DEBUG  auth-service  service processing     {layer: "service", operation: "register"}
// 10:32:01.130  DEBUG  auth-service  repository executing   {layer: "repository", operation: "save_user", db_table: "users"}
// 10:32:01.145  INFO   auth-service  user registered        {layer: "service", user_id: "uuid-xyz", email: "user@test.com"}
// 10:32:01.146  INFO   auth-service  ← response sent        {request_id: "abc-123", status_code: 201, latency_ms: 23}

// ============================================================
// EJEMPLO DE OUTPUT EN PRODUCTION (JSON)
// ============================================================
//
// {"timestamp":"2024-01-15T10:32:01.146Z","level":"info","service":"auth-service","msg":"← response sent","request_id":"abc-123","status_code":201,"latency_ms":23}

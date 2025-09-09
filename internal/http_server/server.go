// Package http_server
package http_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/fsd_server/packet"
	"github.com/half-nothing/simple-fsd/internal/http_server/controller"
	mid "github.com/half-nothing/simple-fsd/internal/http_server/middleware"
	impl "github.com/half-nothing/simple-fsd/internal/http_server/service"
	"github.com/half-nothing/simple-fsd/internal/http_server/service/store"
	. "github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/samber/slog-echo"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type HttpServerShutdownCallback struct {
	serverHandler *echo.Echo
}

func NewHttpServerShutdownCallback(serverHandler *echo.Echo) *HttpServerShutdownCallback {
	return &HttpServerShutdownCallback{
		serverHandler: serverHandler,
	}
}

func (hc *HttpServerShutdownCallback) Invoke(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return hc.serverHandler.Shutdown(timeoutCtx)
}

func StartHttpServer(applicationContent *ApplicationContent) {
	config := applicationContent.ConfigManager().Config()
	logger := applicationContent.Logger()

	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	e.Logger.SetLevel(log.OFF)
	httpConfig := config.Server.HttpServer

	switch httpConfig.ProxyType {
	case 0:
		e.IPExtractor = echo.ExtractIPDirect()
	case 1:
		e.IPExtractor = echo.ExtractIPFromXFFHeader()
	case 2:
		e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	default:
		logger.WarnF("Invalid proxy type %d, using default (direct)", httpConfig.ProxyType)
		e.IPExtractor = echo.ExtractIPDirect()
	}

	if config.Server.HttpServer.SSL.ForceSSL {
		e.Use(middleware.HTTPSRedirect())
	}

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: 30 * time.Second}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(ctx echo.Context, err error, stack []byte) error {
			logger.ErrorF("Recovered from a fatal error: %v, stack: %s", err, string(stack))
			return err
		},
	}))

	loggerConfig := slogecho.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,
	}
	e.Use(slogecho.NewWithConfig(logger.LogHandler(), loggerConfig))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            httpConfig.SSL.HstsExpiredTime,
		HSTSExcludeSubdomains: !httpConfig.SSL.IncludeDomain,
	}))
	e.Use(middleware.CORS())
	if httpConfig.BodyLimit != "" {
		e.Use(middleware.BodyLimit(httpConfig.BodyLimit))
	}
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	if httpConfig.Limits.RateLimit <= 0 {
		logger.WarnF("Invalid rate limit value %d, using default 15", httpConfig.Limits.RateLimit)
		httpConfig.Limits.RateLimit = 15
	}

	if httpConfig.Limits.RateLimitDuration <= 0 {
		logger.WarnF("Invalid rate limit duration %v, using default 1m", httpConfig.Limits.RateLimitDuration)
		httpConfig.Limits.RateLimitDuration = time.Minute
	}

	ipPathLimiter := mid.NewSlidingWindowLimiter(
		httpConfig.Limits.RateLimitDuration,
		httpConfig.Limits.RateLimit,
	)
	cleanupInterval := httpConfig.Limits.RateLimitDuration * 2
	if cleanupInterval > time.Hour {
		cleanupInterval = time.Hour
		logger.InfoF("Limiting cleanup interval to 1 hour for efficiency")
	}
	ipPathLimiter.StartCleanup(cleanupInterval)

	whazzupContent := fmt.Sprintf("url0=%s/api/clients", httpConfig.ServerAddress)

	e.Use(mid.RateLimitMiddleware(ipPathLimiter, mid.CombinedKeyFunc))

	jwtConfig := echojwt.Config{
		SigningKey:    []byte(httpConfig.JWT.Secret),
		TokenLookup:   "header:Authorization:Bearer ",
		SigningMethod: "HS512",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(service.Claims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			var data *service.ApiResponse[any]
			switch {
			case errors.Is(err, echojwt.ErrJWTMissing):
				data = service.NewApiResponse[any](&service.ErrMissingOrMalformedJwt, service.Unsatisfied, nil)
			case errors.Is(err, echojwt.ErrJWTInvalid):
				data = service.NewApiResponse[any](&service.ErrInvalidOrExpiredJwt, service.Unsatisfied, nil)
			default:
				data = service.NewApiResponse[any](&service.ErrUnknown, service.Unsatisfied, nil)
			}
			return data.Response(c)
		},
	}

	jwtMiddleware := echojwt.WithConfig(jwtConfig)

	emailService := impl.NewEmailService(logger, config.Server.HttpServer.Email)
	impl.InitValidator(config.Server.HttpServer.Limits)

	var storeService service.StoreServiceInterface
	storeService = store.NewLocalStoreService(logger, httpConfig.Store)
	switch httpConfig.Store.StoreType {
	case 1:
		storeService = store.NewALiYunOssStoreService(logger, httpConfig.Store, storeService)
	case 2:
		storeService = store.NewTencentCosStoreService(logger, httpConfig.Store, storeService)
	}

	userOperation := applicationContent.Operations().UserOperation()
	historyOperation := applicationContent.Operations().HistoryOperation()
	auditLogOperation := applicationContent.Operations().AuditLogOperation()
	activityOperation := applicationContent.Operations().ActivityOperation()

	userService := impl.NewUserService(logger, httpConfig, userOperation, historyOperation, auditLogOperation, storeService, emailService)
	clientManager := packet.NewClientManager(applicationContent)
	clientService := impl.NewClientService(logger, httpConfig, userOperation, auditLogOperation, clientManager, emailService)
	serverService := impl.NewServerService(logger, config.Server, userOperation, activityOperation)
	activityService := impl.NewActivityService(logger, httpConfig, userOperation, activityOperation, auditLogOperation, storeService)
	auditLogService := impl.NewAuditService(logger, auditLogOperation)

	userController := controller.NewUserHandler(logger, userService)
	emailController := controller.NewEmailController(logger, emailService)
	clientController := controller.NewClientController(logger, clientService)
	serverController := controller.NewServerController(logger, serverService)
	activityController := controller.NewActivityController(logger, activityService)
	fileController := controller.NewFileController(logger, storeService)
	auditLogController := controller.NewAuditLogController(logger, auditLogService)

	apiGroup := e.Group("/api")
	apiGroup.POST("/sessions", userController.UserLogin)
	apiGroup.GET("/sessions", userController.GetToken, jwtMiddleware)
	apiGroup.POST("/codes", emailController.SendVerifyEmail)
	apiGroup.GET("/profile", userController.GetCurrentUserProfile, jwtMiddleware)
	apiGroup.PATCH("/profile", userController.EditCurrentProfile, jwtMiddleware)
	apiGroup.GET("/history", userController.GetUserHistory, jwtMiddleware)

	userGroup := apiGroup.Group("/users")
	userGroup.POST("", userController.UserRegister)
	userGroup.GET("", userController.GetUsers, jwtMiddleware)
	userGroup.GET("/controllers", userController.GetControllers, jwtMiddleware)
	userGroup.GET("/availability", userController.CheckUserAvailability)
	userGroup.GET("/:uid/profile", userController.GetUserProfile, jwtMiddleware)
	userGroup.PATCH("/:uid/profile", userController.EditProfile, jwtMiddleware)
	userGroup.PATCH("/:uid/permission", userController.EditUserPermission, jwtMiddleware)
	userGroup.PUT("/:uid/rating", userController.EditUserRating, jwtMiddleware)

	clientGroup := apiGroup.Group("/clients")
	clientGroup.GET("/status", func(c echo.Context) error { return c.String(http.StatusOK, whazzupContent) })
	clientGroup.GET("", clientController.GetOnlineClients)
	clientGroup.GET("/paths", clientController.GetClientPath, jwtMiddleware)
	clientGroup.POST("/:callsign/message", clientController.SendMessageToClient, jwtMiddleware)
	clientGroup.DELETE("/:callsign", clientController.KillClient, jwtMiddleware)

	serverGroup := apiGroup.Group("/server")
	serverGroup.GET("/config", serverController.GetServerConfig)
	serverGroup.GET("/info", serverController.GetServerInfo, jwtMiddleware)
	serverGroup.GET("/rating", serverController.GetServerOnlineTime, jwtMiddleware)

	activityGroup := apiGroup.Group("/activities")
	activityGroup.GET("", activityController.GetActivities, jwtMiddleware)
	activityGroup.GET("/list", activityController.GetActivitiesPage, jwtMiddleware)
	activityGroup.GET("/:id", activityController.GetActivityInfo, jwtMiddleware)
	activityGroup.POST("", activityController.AddActivity, jwtMiddleware)
	activityGroup.DELETE("/:id", activityController.DeleteActivity, jwtMiddleware)
	activityGroup.POST("/:id/controllers/:facility_id", activityController.ControllerJoin, jwtMiddleware)
	activityGroup.DELETE("/:id/controllers/:facility_id", activityController.ControllerLeave, jwtMiddleware)
	activityGroup.POST("/:id/pilots", activityController.PilotJoin, jwtMiddleware)
	activityGroup.DELETE("/:id/pilots", activityController.PilotLeave, jwtMiddleware)
	activityGroup.PUT("/:id/status", activityController.EditActivityStatus, jwtMiddleware)
	activityGroup.PUT("/:id/pilots/:pilot_id/status", activityController.EditPilotStatus, jwtMiddleware)
	activityGroup.PUT("/:id", activityController.EditActivity, jwtMiddleware)

	fileGroup := apiGroup.Group("/files")
	fileGroup.POST("/images", fileController.UploadImages, jwtMiddleware)

	auditLogGroup := apiGroup.Group("/audits")
	auditLogGroup.GET("", auditLogController.GetAuditLogs, jwtMiddleware)
	auditLogGroup.POST("/unlawful_overreach", auditLogController.LogUnlawfulOverreach, jwtMiddleware)

	apiGroup.Use(middleware.Static(httpConfig.Store.LocalStorePath))

	applicationContent.Cleaner().Add(NewHttpServerShutdownCallback(e))

	protocol := "http"
	if httpConfig.SSL.Enable {
		protocol = "https"
	}
	logger.InfoF("Starting %s server on %s", protocol, httpConfig.Address)
	logger.InfoF("Rate limit: %d requests per %v",
		httpConfig.Limits.RateLimit,
		httpConfig.Limits.RateLimitDuration)

	var err error
	if httpConfig.SSL.Enable {
		err = e.StartTLS(
			httpConfig.Address,
			httpConfig.SSL.CertFile,
			httpConfig.SSL.KeyFile,
		)
	} else {
		err = e.Start(httpConfig.Address)
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Http fsd_server error: %v", err)
	}
}

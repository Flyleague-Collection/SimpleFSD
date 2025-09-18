// Package http_server
package http_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/half-nothing/simple-fsd/internal/http_server/controller"
	mid "github.com/half-nothing/simple-fsd/internal/http_server/middleware"
	impl "github.com/half-nothing/simple-fsd/internal/http_server/service"
	"github.com/half-nothing/simple-fsd/internal/http_server/service/store"
	. "github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
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
	logger := applicationContent.Logger().HttpLogger()

	logger.Info("Http server initializing...")
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
				data = service.NewApiResponse[any](service.ErrMissingOrMalformedJwt, nil)
			case errors.Is(err, echojwt.ErrJWTInvalid):
				data = service.NewApiResponse[any](service.ErrInvalidOrExpiredJwt, nil)
			default:
				data = service.NewApiResponse[any](service.ErrUnknownJwtError, nil)
			}
			return data.Response(c)
		},
	}

	jwtMiddleware := echojwt.WithConfig(jwtConfig)

	jwtVerifyMiddleWare := func(flushToken bool) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {
				token := ctx.Get("user").(*jwt.Token)
				claim := token.Claims.(*service.Claims)
				if flushToken == claim.FlushToken {
					return next(ctx)
				}
				return service.NewApiResponse[any](service.ErrInvalidJwtType, nil).Response(ctx)
			}
		}
	}

	requireNoFlushToken := jwtVerifyMiddleWare(false)
	requireFlushToken := jwtVerifyMiddleWare(true)

	logger.Info("Service initializing...")

	impl.InitValidator(config.Server.HttpServer.Limits)

	userOperation := applicationContent.Operations().UserOperation()
	controllerOperation := applicationContent.Operations().ControllerOperation()
	controllerRecordOperation := applicationContent.Operations().ControllerRecordOperation()
	historyOperation := applicationContent.Operations().HistoryOperation()
	auditLogOperation := applicationContent.Operations().AuditLogOperation()
	activityOperation := applicationContent.Operations().ActivityOperation()
	ticketOperation := applicationContent.Operations().TicketOperation()
	flightPlanOperation := applicationContent.Operations().FlightPlanOperation()

	emailService := impl.NewEmailService(logger, config.Server.HttpServer.Email)

	messageQueue := applicationContent.MessageQueue()
	messageQueue.Subscribe(queue.SendVerifyEmail, emailService.HandleSendVerifyEmailMessage)
	messageQueue.Subscribe(queue.SendRatingChangeEmail, emailService.HandleSendRatingChangeEmailMessage)
	messageQueue.Subscribe(queue.SendPermissionChangeEmail, emailService.HandleSendPermissionChangeEmailMessage)
	messageQueue.Subscribe(queue.SendPasswordChangeEmail, emailService.HandleSendPermissionChangeEmailMessage)
	messageQueue.Subscribe(queue.SendKickedFromServerEmail, emailService.HandleSendKickedFromServerEmailMessage)

	auditLogService := impl.NewAuditService(logger, auditLogOperation)
	messageQueue.Subscribe(queue.AuditLog, auditLogService.HandleAuditLogMessage)
	messageQueue.Subscribe(queue.AuditLogs, auditLogService.HandleAuditLogsMessage)

	clientManager := applicationContent.ClientManager()

	var storeService service.StoreServiceInterface
	storeService = store.NewLocalStoreService(logger, httpConfig.Store, messageQueue, auditLogOperation)
	switch httpConfig.Store.StoreType {
	case 1:
		storeService = store.NewALiYunOssStoreService(logger, httpConfig.Store, storeService, messageQueue, auditLogOperation)
	case 2:
		storeService = store.NewTencentCosStoreService(logger, httpConfig.Store, storeService, messageQueue, auditLogOperation)
	}

	userService := impl.NewUserService(logger, httpConfig, messageQueue, userOperation, historyOperation, auditLogOperation, storeService, emailService)
	clientService := impl.NewClientService(logger, httpConfig, userOperation, auditLogOperation, clientManager, messageQueue)
	serverService := impl.NewServerService(logger, config.Server, userOperation, controllerOperation, activityOperation)
	activityService := impl.NewActivityService(logger, httpConfig, messageQueue, userOperation, activityOperation, auditLogOperation, storeService)
	controllerService := impl.NewControllerService(logger, httpConfig, messageQueue, userOperation, controllerOperation, controllerRecordOperation, auditLogOperation)
	ticketService := impl.NewTicketService(logger, messageQueue, userOperation, ticketOperation, auditLogOperation)
	flightPlanService := impl.NewFlightPlanService(logger, messageQueue, userOperation, flightPlanOperation, auditLogOperation)

	logger.Info("Controller initializing...")

	userController := controller.NewUserHandler(logger, userService)
	emailController := controller.NewEmailController(logger, emailService)
	clientController := controller.NewClientController(logger, clientService)
	serverController := controller.NewServerController(logger, serverService)
	activityController := controller.NewActivityController(logger, activityService)
	fileController := controller.NewFileController(logger, storeService)
	auditLogController := controller.NewAuditLogController(logger, auditLogService)
	controllerController := controller.NewATCController(logger, controllerService)
	ticketController := controller.NewTicketController(logger, ticketService)
	flightPlanController := controller.NewFlightPlanController(logger, flightPlanService)

	logger.Info("Applying router...")

	apiGroup := e.Group("/api")
	apiGroup.POST("/codes", emailController.SendVerifyEmail)

	userGroup := apiGroup.Group("/users")
	userGroup.POST("", userController.UserRegister)
	userGroup.GET("", userController.GetUsers, jwtMiddleware, requireNoFlushToken)
	userGroup.POST("/sessions", userController.UserLogin)
	userGroup.GET("/sessions", userController.GetToken, jwtMiddleware, requireFlushToken)
	userGroup.GET("/availability", userController.CheckUserAvailability)
	userGroup.GET("/histories/self", userController.GetUserHistory, jwtMiddleware, requireNoFlushToken)
	userGroup.GET("/profiles/self", userController.GetCurrentUserProfile, jwtMiddleware, requireNoFlushToken)
	userGroup.PATCH("/profiles/self", userController.EditCurrentProfile, jwtMiddleware, requireNoFlushToken)
	userGroup.GET("/profiles/:uid", userController.GetUserProfile, jwtMiddleware, requireNoFlushToken)
	userGroup.PATCH("/profiles/:uid", userController.EditProfile, jwtMiddleware, requireNoFlushToken)
	userGroup.PATCH("/profiles/:uid/permission", userController.EditUserPermission, jwtMiddleware, requireNoFlushToken)

	controllerGroup := apiGroup.Group("/controllers")
	controllerGroup.GET("", controllerController.GetControllers, jwtMiddleware, requireNoFlushToken)
	controllerGroup.GET("/records/self", controllerController.GetCurrentControllerRecord, jwtMiddleware, requireNoFlushToken)
	controllerGroup.GET("/records/:uid", controllerController.GetControllerRecord, jwtMiddleware, requireNoFlushToken)
	controllerGroup.POST("/records/:uid", controllerController.AddControllerRecord, jwtMiddleware, requireNoFlushToken)
	controllerGroup.DELETE("/records/:uid/:rid", controllerController.DeleteControllerRecord, jwtMiddleware, requireNoFlushToken)
	controllerGroup.PUT("/:uid/rating", controllerController.UpdateControllerRating, jwtMiddleware, requireNoFlushToken)
	controllerGroup.PUT("/:uid/um", controllerController.SetControllerUnderMonitor, jwtMiddleware, requireNoFlushToken)
	controllerGroup.DELETE("/:uid/um", controllerController.UnsetControllerUnderMonitor, jwtMiddleware, requireNoFlushToken)
	controllerGroup.PUT("/:uid/solo", controllerController.SetControllerUnderSolo, jwtMiddleware, requireNoFlushToken)
	controllerGroup.DELETE("/:uid/solo", controllerController.UnsetControllerUnderSolo, jwtMiddleware, requireNoFlushToken)
	controllerGroup.PUT("/:uid/guest", controllerController.SetControllerGuest, jwtMiddleware, requireNoFlushToken)
	controllerGroup.DELETE("/:uid/guest", controllerController.UnsetControllerGuest, jwtMiddleware, requireNoFlushToken)

	clientGroup := apiGroup.Group("/clients")
	clientGroup.GET("", clientController.GetOnlineClients)
	clientGroup.GET("/status", func(c echo.Context) error { return c.String(http.StatusOK, whazzupContent) })
	clientGroup.GET("/paths/:callsign", clientController.GetClientPath, jwtMiddleware, requireNoFlushToken)
	clientGroup.POST("/messages/:callsign", clientController.SendMessageToClient, jwtMiddleware, requireNoFlushToken)
	clientGroup.DELETE("/:callsign", clientController.KillClient, jwtMiddleware, requireNoFlushToken)

	serverGroup := apiGroup.Group("/server")
	serverGroup.GET("/config", serverController.GetServerConfig)
	serverGroup.GET("/info", serverController.GetServerInfo, jwtMiddleware, requireNoFlushToken)
	serverGroup.GET("/rating", serverController.GetServerOnlineTime, jwtMiddleware, requireNoFlushToken)

	activityGroup := apiGroup.Group("/activities")
	activityGroup.GET("", activityController.GetActivities, jwtMiddleware, requireNoFlushToken)
	activityGroup.GET("/pages", activityController.GetActivitiesPage, jwtMiddleware, requireNoFlushToken)
	activityGroup.GET("/:activity_id", activityController.GetActivityInfo, jwtMiddleware, requireNoFlushToken)
	activityGroup.POST("", activityController.AddActivity, jwtMiddleware, requireNoFlushToken)
	activityGroup.DELETE("/:activity_id", activityController.DeleteActivity, jwtMiddleware, requireNoFlushToken)
	activityGroup.POST("/:activity_id/controllers/:facility_id", activityController.ControllerJoin, jwtMiddleware, requireNoFlushToken)
	activityGroup.DELETE("/:activity_id/controllers/:facility_id", activityController.ControllerLeave, jwtMiddleware, requireNoFlushToken)
	activityGroup.POST("/:activity_id/pilots", activityController.PilotJoin, jwtMiddleware, requireNoFlushToken)
	activityGroup.DELETE("/:activity_id/pilots", activityController.PilotLeave, jwtMiddleware, requireNoFlushToken)
	activityGroup.PUT("/:activity_id/status", activityController.EditActivityStatus, jwtMiddleware, requireNoFlushToken)
	activityGroup.PUT("/:activity_id/pilots/:user_id/status", activityController.EditPilotStatus, jwtMiddleware, requireNoFlushToken)
	activityGroup.PUT("/:activity_id", activityController.EditActivity, jwtMiddleware, requireNoFlushToken)

	ticketGroup := apiGroup.Group("/tickets")
	ticketGroup.GET("", ticketController.GetTickets, jwtMiddleware, requireNoFlushToken)
	ticketGroup.GET("/self", ticketController.GetUserTickets, jwtMiddleware, requireNoFlushToken)
	ticketGroup.POST("", ticketController.CreateTicket, jwtMiddleware, requireNoFlushToken)
	ticketGroup.PUT("/:tid", ticketController.CloseTicket, jwtMiddleware, requireNoFlushToken)
	ticketGroup.DELETE("/:tid", ticketController.DeleteTicket, jwtMiddleware, requireNoFlushToken)

	flightPlanGroup := apiGroup.Group("/plans")
	flightPlanGroup.POST("", flightPlanController.SubmitFlightPlan, jwtMiddleware, requireNoFlushToken)
	flightPlanGroup.GET("", flightPlanController.GetFlightPlans, jwtMiddleware, requireNoFlushToken)
	flightPlanGroup.GET("/self", flightPlanController.GetFlightPlan, jwtMiddleware, requireNoFlushToken)
	flightPlanGroup.PUT("/:cid/lock", flightPlanController.LockFlightPlan, jwtMiddleware, requireNoFlushToken)
	flightPlanGroup.DELETE("/:cid/lock", flightPlanController.UnlockFlightPlan, jwtMiddleware, requireNoFlushToken)
	flightPlanGroup.DELETE("/:cid", flightPlanController.DeleteFlightPlan, jwtMiddleware, requireNoFlushToken)

	fileGroup := apiGroup.Group("/files")
	fileGroup.POST("/images", fileController.UploadImages, jwtMiddleware, requireNoFlushToken)

	auditLogGroup := apiGroup.Group("/audits")
	auditLogGroup.GET("", auditLogController.GetAuditLogs, jwtMiddleware, requireNoFlushToken)
	auditLogGroup.POST("/unlawful_overreach", auditLogController.LogUnlawfulOverreach, jwtMiddleware, requireNoFlushToken)

	apiGroup.Use(middleware.Static(httpConfig.Store.LocalStorePath))

	applicationContent.Cleaner().Add(NewHttpServerShutdownCallback(e))

	protocol := "http"
	if httpConfig.SSL.Enable {
		protocol = "https"
	}
	logger.InfoF("Starting %s server on %s", protocol, httpConfig.Address)
	logger.InfoF("Rate limit: %d requests per %v", httpConfig.Limits.RateLimit, httpConfig.Limits.RateLimitDuration)

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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	serverhttp "github.com/vtievsky/auth-id/gen/httpserver/auth-id"
	"github.com/vtievsky/auth-id/internal/conf"
	"github.com/vtievsky/auth-id/internal/httptransport"
	otelclient "github.com/vtievsky/auth-id/internal/otel/client"
	otelmetrics "github.com/vtievsky/auth-id/internal/otel/metrics"
	oteltracing "github.com/vtievsky/auth-id/internal/otel/tracing"
	clienttarantool "github.com/vtievsky/auth-id/internal/repositories/db/client/tarantool"
	tarantoolprivileges "github.com/vtievsky/auth-id/internal/repositories/db/privileges"
	tarantoolroles "github.com/vtievsky/auth-id/internal/repositories/db/roles"
	tarantoolusers "github.com/vtievsky/auth-id/internal/repositories/db/users"
	clientredis "github.com/vtievsky/auth-id/internal/repositories/sessions/client/redis"
	reposessions "github.com/vtievsky/auth-id/internal/repositories/sessions/sessions"
	"github.com/vtievsky/auth-id/internal/services"
	privilegesvc "github.com/vtievsky/auth-id/internal/services/privileges"
	roleprivilegesvc "github.com/vtievsky/auth-id/internal/services/role-privileges"
	roleusersvc "github.com/vtievsky/auth-id/internal/services/role-users"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
	sessionsvc "github.com/vtievsky/auth-id/internal/services/sessions"
	userprivilegesvc "github.com/vtievsky/auth-id/internal/services/user-privileges"
	userrolesvc "github.com/vtievsky/auth-id/internal/services/user-roles"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
	"github.com/vtievsky/golibs/runtime/logger"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	conf := conf.New()
	logger := logger.CreateZapLogger(conf.Debug, conf.Log.EnableStacktrace)
	ctx := context.Background()

	otelClient, err := otelclient.New(&otelclient.OtelOpts{
		URL: conf.Metrics.URL,
	})
	if err != nil {
		log.Fatal(err)
	}

	otelResource, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(conf.ServiceName),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	shutdownMeterProvider, err := otelmetrics.InitMeterProvider(ctx, otelResource, otelClient)
	if err != nil {
		log.Fatal(err)
	}

	shutdownTracerProvider, err := oteltracing.InitTracerProvider(ctx, otelResource, otelClient)
	if err != nil {
		log.Fatal(err)
	}

	dbClient, err := clienttarantool.New(&clienttarantool.ClientOpts{
		URL:       conf.DB.URL,
		RateLimit: 25, //nolint:mnd
	})
	if err != nil {
		log.Fatal(err)
	}

	sessionClient, err := clientredis.New(&clientredis.ClientOpts{
		URL: conf.Session.URL,
	})
	if err != nil {
		log.Fatal(err)
	}

	// repos
	usersRepo := tarantoolusers.New(&tarantoolusers.UsersOpts{
		Client: dbClient,
	})

	rolesRepo := tarantoolroles.New(&tarantoolroles.RolesOpts{
		Client: dbClient,
	})

	privilegesRepo := tarantoolprivileges.New(&tarantoolprivileges.PrivilegesOpts{
		Client: dbClient,
	})

	sessionsRepo := reposessions.New(&reposessions.SessionsOpts{
		Logger: logger.Named("session"),
		Client: sessionClient,
	})

	// services
	userService := usersvc.New(&usersvc.UserSvcOpts{
		Logger:  logger.Named("user"),
		Storage: usersRepo,
	})

	roleService := rolesvc.New(&rolesvc.RoleSvcOpts{
		Logger:  logger.Named("role"),
		Storage: rolesRepo,
	})

	userRoleService := userrolesvc.New(&userrolesvc.UserRoleSvcOpts{
		Logger:  logger.Named("user-role"),
		Storage: usersRepo,
		RoleSvc: roleService,
	})

	roleUserService := roleusersvc.New(&roleusersvc.RoleUserSvcOpts{
		Logger:  logger.Named("role-user"),
		Storage: rolesRepo,
		RoleSvc: roleService,
		UserSvc: userService,
	})

	privilegeService := privilegesvc.New(&privilegesvc.PrivilegeSvcOpts{
		Logger:  logger.Named("privilege"),
		Storage: privilegesRepo,
	})

	rolePrivilegeService := roleprivilegesvc.New(&roleprivilegesvc.RolePrivilegeSvcOpts{
		Logger:       logger.Named("role-privilege"),
		Storage:      rolesRepo,
		RoleSvc:      roleService,
		PrivilegeSvc: privilegeService,
	})

	userPrivilegeService := userprivilegesvc.New(&userprivilegesvc.UserPrivilegeSvcOpts{
		Logger:           logger.Named("user-privilege"),
		UserRoleSvc:      userRoleService,
		RolePrivilegeSvc: rolePrivilegeService,
	})

	sessionService := sessionsvc.New(&sessionsvc.SessionSvcOpts{
		Logger:           logger.Named("session"),
		Storage:          sessionsRepo,
		UserSvc:          userService,
		UserPrivilegeSvc: userPrivilegeService,
		SessionTTL:       conf.Session.SessionTTL,
		AccessTokenTTL:   conf.Session.AccessTokenTTL,
		RefreshTokenTTL:  conf.Session.RefreshTokenTTL,
		SigningKey:       conf.Session.SigningKey,
	})

	serverCtx, cancel := context.WithCancel(ctx)
	services := &services.SvcLayer{
		UserSvc:          userService,
		UserRoleSvc:      userRoleService,
		UserPrivilegeSvc: userPrivilegeService,
		RoleSvc:          roleService,
		RoleUserSvc:      roleUserService,
		RolePrivilegeSvc: rolePrivilegeService,
		PrivilegeSvc:     privilegeService,
		SessionSvc:       sessionService,
	}

	httpSrv := echo.New()
	httpSrv.HideBanner = true
	httpSrv.Use(
		httptransport.TracerMiddleware(),
		httptransport.LoggerMiddleware(logger),
		httptransport.AuthorizationMiddleware(
			sessionService,
			conf.Session.SigningKey,
		),
	)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(signalChannel) // Отмена подписки на системные события
	defer stopApp(
		ctx,
		logger,
		httpSrv,
		shutdownMeterProvider,
		shutdownTracerProvider,
	)

	go startApp(
		cancel,
		logger,
		httpSrv,
		httptransport.New(
			conf,
			services,
		),
		conf.Port,
	)

	for {
		select {
		case <-signalChannel:
			logger.Info("interrupted by a signal")

			return
		case <-serverCtx.Done():
			return
		}
	}
}

func stopApp(
	ctx context.Context,
	logger *zap.Logger,
	httpSrv *echo.Echo,
	meterShutdown otelmetrics.MeterShutdown,
	tracerShutdown oteltracing.TracerShutdown,
) {
	defer func(alogger *zap.Logger) {
		alogger.Debug("sync zap logs")

		_ = alogger.Sync()
	}(logger)

	err := httpSrv.Close()
	if err != nil {
		logger.Error("failed to close http server",
			zap.Error(err),
		)
	}

	g := errgroup.Group{}

	g.Go(func() error {
		err := tracerShutdown(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		err := meterShutdown(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	err = g.Wait()
	if err != nil {
		logger.Error("failed to shutdown Opentelemetry",
			zap.Error(err),
		)
	}
}

func startApp(
	cancel context.CancelFunc,
	logger *zap.Logger,
	httpSrv *echo.Echo,
	handler serverhttp.ServerInterface,
	port int,
) {
	defer cancel()

	serverhttp.RegisterHandlers(httpSrv, handler)

	if err := httpSrv.Start(fmt.Sprintf(":%d", port)); err != nil {
		logger.Error("error while serve http server",
			zap.Error(err),
		)
	}
}

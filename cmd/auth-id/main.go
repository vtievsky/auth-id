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
	tarantoolclient "github.com/vtievsky/auth-id/internal/repositories/tarantool/client"
	tarantoolprivileges "github.com/vtievsky/auth-id/internal/repositories/tarantool/privileges"
	tarantoolroles "github.com/vtievsky/auth-id/internal/repositories/tarantool/roles"
	tarantoolusers "github.com/vtievsky/auth-id/internal/repositories/tarantool/users"
	"github.com/vtievsky/auth-id/internal/services"
	privilegesvc "github.com/vtievsky/auth-id/internal/services/privileges"
	rolesvc "github.com/vtievsky/auth-id/internal/services/roles"
	usersvc "github.com/vtievsky/auth-id/internal/services/users"
	"github.com/vtievsky/golibs/runtime/logger"
	"go.uber.org/zap"
)

func main() {
	conf := conf.New()
	logger := logger.CreateZapLogger(conf.Debug, conf.Log.EnableStacktrace)
	httpSrv := echo.New()

	tarantoolClient, err := tarantoolclient.New(&tarantoolclient.ClientOpts{
		URL:       conf.DB.URL,
		RateLimit: 25, //nolint:mnd
	})
	if err != nil {
		log.Fatal(err)
	}

	// repos
	userRepo := tarantoolusers.New(&tarantoolusers.UsersOpts{
		Client: tarantoolClient,
	})

	roleRepo := tarantoolroles.New(&tarantoolroles.RolesOpts{
		Client: tarantoolClient,
	})

	rolePrivilegeRepo := tarantoolroles.New(&tarantoolroles.RolesOpts{
		Client: tarantoolClient,
	})

	roleUserRepo := tarantoolroles.New(&tarantoolroles.RolesOpts{
		Client: tarantoolClient,
	})

	privilegeRepo := tarantoolprivileges.New(&tarantoolprivileges.PrivilegesOpts{
		Client: tarantoolClient,
	})

	// services
	userService := usersvc.New(&usersvc.UserSvcOpts{
		Logger:  logger.Named("user"),
		Storage: userRepo,
	})

	privilegeService := privilegesvc.New(&privilegesvc.PrivilegeSvcOpts{
		Logger:  logger.Named("privilege"),
		Storage: privilegeRepo,
	})

	roleService := rolesvc.New(&rolesvc.RoleSvcOpts{
		Logger:         logger.Named("role"),
		Roles:          roleRepo,
		RolePrivileges: rolePrivilegeRepo,
		RoleUsers:      roleUserRepo,
		PrivilegeSvc:   privilegeService,
		UserSvc:        userService,
	})

	ctx := context.Background()
	serverCtx, cancel := context.WithCancel(ctx)
	services := &services.SvcLayer{
		UserSvc:          userService,
		RoleSvc:          roleService,
		RolePrivilegeSvc: roleService,
		RoleUserSvc:      roleService,
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(signalChannel) // Отмена подписки на системные события
	defer stopApp(logger, httpSrv)

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

func stopApp(logger *zap.Logger, httpSrv *echo.Echo) {
	defer func(alogger *zap.Logger) {
		alogger.Debug("sync zap logs")

		_ = alogger.Sync()
	}(logger)

	if err := httpSrv.Close(); err != nil {
		logger.Error("failed to close http server",
			zap.Error(err),
		)
	}
}

func startApp(
	cancel context.CancelFunc,
	logger *zap.Logger,
	httpSrv *echo.Echo,
	handlers serverhttp.StrictServerInterface,
	port int,
) {
	defer cancel()

	serverhttp.RegisterHandlers(httpSrv, serverhttp.NewStrictHandler(
		handlers,
		[]serverhttp.StrictMiddlewareFunc{},
	))

	address := fmt.Sprintf(":%d", port)

	if err := httpSrv.Start(address); err != nil {
		logger.Error("error while serve http server",
			zap.Error(err),
		)
	}
}

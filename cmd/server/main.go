package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	"github.com/casbin/casbin/v3"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goforj/wire"
	"github.com/google/uuid"
	_ "go.uber.org/automaxprocs"

	"github.com/tencat-dev/go-base/internal/authz"
	"github.com/tencat-dev/go-base/internal/conf"
	"github.com/tencat-dev/go-base/internal/server"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	hostname string
	id       uuid.UUID
)

// ProviderSetConfig is config providers.
var ProviderSetConfig = wire.NewSet(
	newServer,
	newGrpcServer,
	newHttpServer,
	newPprofServer,
	newAuth,
	newJwtConfig,
	newData,
	newDatabaseConfig,
	newAuthz,
)

func newServer(c *conf.Bootstrap) *conf.Server {
	return c.Server
}
func newGrpcServer(c *conf.Server) *conf.GRPCServer {
	return c.Grpc
}
func newHttpServer(c *conf.Server) *conf.HTTPServer {
	return c.Http
}
func newPprofServer(c *conf.Server) *conf.PprofServer {
	return c.Pprof
}
func newAuth(c *conf.Bootstrap) *conf.Auth {
	return c.Auth
}
func newJwtConfig(c *conf.Auth) *conf.JWT {
	return c.Jwt
}
func newData(c *conf.Bootstrap) *conf.Data {
	return c.Data
}
func newDatabaseConfig(c *conf.Data) *conf.DatabaseConfig {
	return c.Database
}
func newAuthz(c *conf.Bootstrap) *conf.Authz {
	return c.Authz
}

func init() {
	flag.StringVar(&flagconf, "conf", "./configs/config.yaml", "config path, eg: -conf config.yaml")

	// Initialize hostname and ID with proper error handling
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}

	id = uuid.Must(uuid.NewV7())
}

func newApp(
	ctx context.Context,
	confServer *conf.Server,
	confAuthz *conf.Authz,
	logger log.Logger,
	gs server.GrpcServer,
	hs server.HttpServer,
	ps server.PprofServer,
	enforcer casbin.IEnforcer,
	authzRegistry *authz.AuthzRegistry,
) (*kratos.App, func(), error) {
	var srvs []transport.Server

	if confServer.Http.Enable {
		srvs = append(srvs, hs)
	}

	if confServer.Grpc.Enable {
		srvs = append(srvs, gs)
	}

	if confServer.Pprof.Enable {
		srvs = append(srvs, ps)
	}

	if len(srvs) == 0 {
		return nil, nil, errors.New("no server configured")
	}

	if confAuthz.AutoSync {
		err := authz.SyncFromRegistry(enforcer, authzRegistry)
		if err != nil {
			return nil, nil, err
		}
	}

	app := kratos.New(
		kratos.Context(ctx),
		kratos.ID(id.String()),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(srvs...),
	)

	return app, func() {
		if err := app.Stop(); err != nil {
			logger.Log(log.LevelWarn, "app stop error", err)
		}
	}, nil
}

func main() {
	flag.Parse()

	// Create context with signal handling for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var kv = []any{
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service", Name,
		"version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	}
	logger := log.With(log.NewStdLogger(os.Stdout), kv...)
	logHelper := log.NewHelper(logger)

	// Load configuration
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer func() {
		if err := c.Close(); err != nil {
			logHelper.Error("Failed to close config", err)
		}
	}()

	if err := c.Load(); err != nil {
		logHelper.Fatalf("Failed to load config: %v", err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		logHelper.Fatalf("Failed to scan config: %v", err)
	}
	err := protovalidate.Validate(&bc)
	if err != nil {
		logHelper.Fatalf("Failed to validate config: %v", err)
	}

	// Initialize the application
	app, cleanup, err := wireApp(ctx, &bc, logger, logHelper)
	if err != nil {
		logHelper.Fatalf("Failed to initialize application: %v", err)
	}
	defer func() {
		if cleanup != nil {
			cleanup()
		}
	}()

	logHelper.Info(fmt.Sprintf("Starting %s service with version %s", Name, Version))
	logHelper.Info(fmt.Sprintf("Service ID: %s, Hostname: %s", id.String(), hostname))

	// Start the application and wait for stop signal
	if err := app.Run(); err != nil {
		logHelper.Fatalf("Application stopped with error: %v", err)
	}
}

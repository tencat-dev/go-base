package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/goforj/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob/drivers/pgx"

	"github.com/tencat-dev/go-base/internal/conf"
)

// ProviderSetData is data providers.
var ProviderSetData = wire.NewSet(
	NewData,
	NewCasbinEnforcer,
	NewCasbinAuthz,
	NewPermissionChecker,
	NewPermissionManager,
	NewUserRepo,
	NewAuthRepo,
)

// Data wraps database client.
type Data struct {
	db pgx.Pool
}

// NewData creates a new Data instance with PostgreSQL connection.
func NewData(ctx context.Context, c *conf.DatabaseConfig, logHelper *log.Helper) (*Data, func(), error) {
	if c == nil {
		return nil, nil, fmt.Errorf("database configuration is missing")
	}

	config, err := pgxpool.ParseConfig(c.Dsn)
	if err != nil {
		return nil, nil, err
	}
	config.ConnConfig.ConnectTimeout = 5 * time.Second
	config.PingTimeout = 5 * time.Second

	logHelper.Info("Connect to database")

	// Create connection pool
	pool, err := pgx.NewWithConfig(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %v", err)
	}

	logHelper.Info("Connect OKKK")

	d := &Data{
		db: pool,
	}

	cleanup := func() {
		log.Info("closing the database connection pool")
		pool.Close()
	}

	return d, cleanup, nil
}

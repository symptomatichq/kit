// Package dbutil provides utilities for managing connections to a SQL database.
package dbutil

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gocraft/dbr"
	migrate "github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres" // postgres driver
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/pkg/errors"

	"github.com/symptomatic/kit/pgutil"
)

type dbConfig struct {
	logger      log.Logger
	maxAttempts int
	maxConns    int
}

// WithLogger configures a logger Option.
func WithLogger(logger log.Logger) Option {
	return func(c *dbConfig) {
		c.logger = log.With(logger, "component", "database")
	}
}

// WithMaxAttempts configures the number of maximum attempts to make
func WithMaxAttempts(maxAttempts int) Option {
	return func(c *dbConfig) {
		c.maxAttempts = maxAttempts
	}
}

// WithMaxConnections configures the number of maximum number of connections to pool
func WithMaxConnections(maxConns int) Option {
	return func(c *dbConfig) {
		c.maxConns = maxConns
	}
}

// Option provides optional configuration for managing DB connections.
type Option func(*dbConfig)

// Open creates a dbr.Connection connection to the database driver.
// OpenDB uses a linear backoff timer when attempting to establish a connection,
// only returning after the connection is successful or the number of attempts exceeds
// the maxAttempts value (default 15).
func Open(driver, dsn string, opts ...Option) (*dbr.Connection, error) {
	config := &dbConfig{
		logger:      log.NewNopLogger(),
		maxAttempts: 15,
		maxConns:    10,
	}

	for _, opt := range opts {
		opt(config)
	}

	conn, err := dbr.Open(driver, dsn, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "opening %s connection, dsn=%s", driver, dsn)
	}

	conn.SetMaxOpenConns(config.maxConns)

	for attempt := 0; attempt < config.maxAttempts; attempt++ {
		err = conn.Ping()
		if err == nil {
			// we're connected!
			break
		}
		interval := time.Duration(attempt) * time.Second
		config.logger.Log(
			"message", "could not connect to db",
			"error", err.Error(),
		)
		time.Sleep(interval)
	}
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type assetGetter func(name string) ([]byte, error)

// Migrate upgrades the database to the newest version available
func Migrate(co *pgutil.ConnectionOptions, getter assetGetter, migrations []string) (err error) {
	s := bindata.Resource(migrations, func(name string) ([]byte, error) {
		return getter(name)
	})

	d, err := bindata.WithInstance(s)
	if err != nil {
		return
	}
	defer d.Close()

	m, err := migrate.NewWithSourceInstance("go-bindata", d, co.ConnectionString("postgres"))
	if err != nil {
		return
	}
	defer m.Close()

	_, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return
	}

	if dirty {
		err = errors.New("database is dirty unable to apply migrations")
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	} else if err != nil {
		// FIXME: consider force rolling back?
		// 	err = m.Force(int(currentVersion))
		// 	if err != nil {
		// 		return err
		// 	}
		return
	}

	return
}

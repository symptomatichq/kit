// Package pgutil provides utilities for Postgres
package pgutil

import (
	"flag"
	"fmt"

	"github.com/symptomatic/kit/env"
)

// ConfigFromEnv returns the connection configuration based on flags
func ConfigFromEnv() *ConnectionOptions {
	host := flag.String("database.host", env.String("DATABASE_HOST", "localhost"), "PostgreSQL server host")
	port := flag.Int("database.port", env.Int("DATABASE_PORT", 5432), "PostgreSQL server port")
	name := flag.String("database.name", env.String("DATABASE_NAME", "olympus"), "PostgreSQL database name")
	user := flag.String("database.user", env.String("DATABASE_USER", "olympus"), "PostgreSQL server user")
	password := flag.String("database.password", env.String("DATABASE_PASSWORD", "olympus"), "PostgreSQL server password")
	sslMode := flag.Bool("database.ssl", env.Bool("DATABASE_SSL", false), "PostgreSQL server ssl mode")

	return &ConnectionOptions{
		Host:     *host,
		Port:     *port,
		User:     *user,
		Password: *password,
		DBName:   *name,
		SSL:      *sslMode,
	}
}

// ConnectionOptions represents the configurable options of a connection to a
// Postgres database
type ConnectionOptions struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSL      bool
}

func (c *ConnectionOptions) sslMode() (mode string) {
	mode = "disable"
	if c.SSL == true {
		mode = "enable"
	}

	return mode
}

// String implements the Stringer interface so that a pgutil.ConnectionOptions
// can be converted into a value key/value connection string
func (c *ConnectionOptions) String() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.sslMode(),
	)
}

// ConnectionString returns the properly formatted connection string for a specified driver
func (c *ConnectionOptions) ConnectionString(driver string) string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		driver, c.User, c.Password, c.Host, c.Port, c.DBName, c.sslMode(),
	)
}

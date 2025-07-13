package main

import (
	"database/sql"
	"flag"
	"fmt"
	"runtime/debug"

	"github.com/go-gorp/gorp/v3"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

var dialects = map[string]gorp.Dialect{
	"sqlite3":  &gorp.SqliteDialect{},
	"postgres": &gorp.PostgresDialect{},
	"mysql":    &gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"},
}

var (
	ConfigFile        string
	ConfigEnvironment string
)

func ConfigFlags(f *flag.FlagSet) {
	f.StringVar(&ConfigFile, "config", "dbconfig.yml", "Configuration file to use.")
	f.StringVar(&ConfigEnvironment, "env", "development", "Environment to use.")
}

func GetEnvironment() (*migrate.Environment, error) {
	return migrate.GetEnvironment(ConfigFile, ConfigEnvironment)
}

func GetConnection(env *migrate.Environment) (*sql.DB, map[string]gorp.Dialect, error) {
	db, err := sql.Open(env.Dialect, env.DataSource)
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot connect to database: %w", err)
	}

	_, exists := dialects[env.Dialect]
	if !exists {
		return nil, nil, fmt.Errorf("Unsupported dialect: %s", env.Dialect)
	}

	return db, dialects, nil
}

func GetVersion() string {
	if buildInfo, ok := debug.ReadBuildInfo(); ok && buildInfo.Main.Version != "(devel)" {
		return buildInfo.Main.Version
	}
	return "dev"
}

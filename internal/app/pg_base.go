package app

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/tsatke/verbose-broccoli/internal/app/config"
)

//go:embed init.sql
var InitSQL string

type PostgresDatabaseProvider struct {
	Config config.Config
	DB     *sql.DB
}

func NewPostgresDatabaseProvider(log zerolog.Logger, cfg config.Config, ssl bool) (*PostgresDatabaseProvider, error) {
	endpoint := cfg.GetString(config.PGEndpoint)
	port := cfg.GetString(config.PGPort)
	user := cfg.GetString(config.PGUsername)
	pass := cfg.GetString(config.PGPassword)
	database := cfg.GetString(config.PGDatabase)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", endpoint, port, user, pass, database)
	if !ssl {
		dsn += " sslmode=disable"
	}

	log.
		Info().
		Str("host", endpoint).
		Str("port", port).
		Str("database", database).
		Msg("connect to database")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	p := &PostgresDatabaseProvider{
		Config: cfg,
		DB:     db,
	}

	start := time.Now()
	if err := p.initDB(); err != nil {
		return nil, fmt.Errorf("init DB: %w", err)
	}
	log.
		Info().
		Stringer("took", time.Since(start)).
		Msg("initialize database")

	return p, nil
}

func (i *PostgresDatabaseProvider) tx(fn func(tx *sql.Tx) error) error {
	return tx(i.DB, fn)
}

func (i *PostgresDatabaseProvider) initDB() error {
	return i.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(InitSQL)
		if err != nil {
			return fmt.Errorf("exec init: %w", err)
		}
		return nil
	})
}

func tx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // will be ignored if tx was already committed
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

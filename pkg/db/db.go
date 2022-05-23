package db

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hashicorp/go-hclog"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB(dbURL string) (*sqlx.DB, error) {
	dbConn, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		fmt.Println("failed to connect to DB", "error", err)
		return nil, err
	}

	return dbConn, err
}

func MigrateToDB(dbConn string,logger hclog.Logger) error {
	logger.Debug("preparing migration")
	m, err := migrate.New("file://schema/", dbConn)
	if err != nil {
		logger.Error("failed to prepare migration", "error", err)
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		logger.Debug("no change")
		return nil
	}
	if err != nil {
		logger.Error("failed to run migration", "error", err)
		return err
	}

	logger.Info("migration done")
	return nil
}

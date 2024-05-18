package db

import (
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

const DATABASE_PATH = "./internal/adapters/storage/db.sqlite"

var dbURL = os.Getenv("DATABASE_URL")

type DB struct {
	*sql.DB
	url string
}

//go:embed migrations/*.sql
var migrationFs embed.FS

func New() (*DB, error) {
	if dbURL == "" {
		dbURL = "sqlite://./internal/adapters/storage/db.sqlite"
	}
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		return nil, err
	}
	return &DB{db, dbURL}, nil
}

// migration possibility 1
func (db *DB) Migrate() error {
	// return the source driver from migrations files
	driver, err := iofs.New(migrationFs, "migrations")
	if err != nil {
		log.Println("Error creating migration driver:", err)
		return err
	}

	// migration instance from driver source and database url
	migrations, err := migrate.NewWithSourceInstance("iofs", driver, db.url)
	if err != nil {
		log.Println("Error creating migration instance:", err)
		return err
	}

	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("Error applying migrations:", err)
		return err
	}

	// if err := migrations.Force(13); err != nil && err != migrate.ErrNoChange {
	// 	log.Println("Error applying migrations:", err)
	// 	return err
	// }

	return nil
}

// migrations possibility 2
func (db *DB) MigrateBis() error {
	currenDir, _ := os.Getwd()
	migrationsDir := currenDir + "/internal/adapters/storage/migrations"
	// get migration instance from source and database url
	migrationInst, err := migrate.New("file://"+migrationsDir, "sqlite://"+dbURL)
	if err != nil {
		log.Println("Error creating migration instance:", err)
		return err
	}

	// migrate all the way up
	if err := migrationInst.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("Error applying migrations:", err)
		return err
	}
	return nil
}

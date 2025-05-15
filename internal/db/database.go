package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"plefi/internal/config"
	"plefi/internal/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB Database

type Database interface {
	Migrate(ctx context.Context) error
	SavePlexToken(ctx context.Context, tok models.PlexToken) error
	GetPlexToken(ctx context.Context, userID int) (*models.PlexToken, error)

	// Invite Code operations
	SaveInviteCode(ctx context.Context, inviteCode models.InviteCode) (int, error)
	GetInviteCode(ctx context.Context, id int) (*models.InviteCode, error)
	GetInviteCodeByCode(ctx context.Context, code string) (*models.InviteCode, error)
	UpdateInviteCodeUsage(ctx context.Context, codeID int) error
	ListActiveInviteCodes(ctx context.Context) ([]models.InviteCode, error)

	// Plex User operations
	SavePlexUser(ctx context.Context, user models.PlexUser) error
	GetPlexUser(ctx context.Context, userID int) (*models.PlexUser, error)
	GetPlexUserByEmail(ctx context.Context, email string) (*models.PlexUser, error)
	GetAllPlexUsers(ctx context.Context) ([]models.PlexUser, error)
	DeletePlexUser(ctx context.Context, userID int) error
	UpdateUserNotes(ctx context.Context, userID int, notes string) error

	// Plex User Invite operations
	AssociatePlexUserWithInviteCode(ctx context.Context, userID, inviteCodeID int) error
	GetPlexUserInvites(ctx context.Context, userID int) ([]models.PlexUserInvite, error)
	GetUsersWithActiveInviteCode(ctx context.Context, inviteCodeID int) ([]models.PlexUser, error)
	DisableInviteCode(ctx context.Context, codeID int) error
}

type sqlDB struct {
	conn   *sql.DB
	driver string
}

func Init(driver, dsn string) error {
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	DB = &sqlDB{conn: conn, driver: driver}
	return nil
}

func (db *sqlDB) Migrate(ctx context.Context) error {
	migrationsPath := config.C.Database.MigrationsPath
	// Setup the database driver
	var instance *migrate.Migrate
	var err error

	slog.Info("Running database migrations", "driver", db.driver, "path", migrationsPath)
	switch db.driver {
	case "postgres":
		driver, driverErr := postgres.WithInstance(db.conn, &postgres.Config{})
		if driverErr != nil {
			return fmt.Errorf("failed to create postgres driver: %w", driverErr)
		}
		instance, err = migrate.NewWithDatabaseInstance(
			"file://"+migrationsPath,
			"postgres", driver)
	case "sqlite3":
		driver, driverErr := sqlite3.WithInstance(db.conn, &sqlite3.Config{})
		if driverErr != nil {
			return fmt.Errorf("failed to create sqlite3 driver: %w", driverErr)
		}
		instance, err = migrate.NewWithDatabaseInstance(
			"file://"+migrationsPath,
			"sqlite3", driver)
	default:
		return fmt.Errorf("unsupported database driver: %s", db.driver)
	}

	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	if err := instance.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

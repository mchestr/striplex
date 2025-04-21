package db

import (
	"context"
	"database/sql"
	"plefi/internal/models"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB Database

type Database interface {
	Migrate(ctx context.Context) error
	SavePlexToken(ctx context.Context, tok models.PlexToken) error
	GetPlexToken(ctx context.Context, userID int) (models.PlexToken, error)

	// Invite Code operations
	SaveInviteCode(ctx context.Context, inviteCode models.InviteCode) (int, error)
	GetInviteCodeByCode(ctx context.Context, code string) (*models.InviteCode, error)
	UpdateInviteCodeUsage(ctx context.Context, codeID int) error
	ListActiveInviteCodes(ctx context.Context) ([]models.InviteCode, error)

	// Plex User operations
	SavePlexUser(ctx context.Context, user models.PlexUser) error
	GetPlexUser(ctx context.Context, userID int) (*models.PlexUser, error)
	GetPlexUserByEmail(ctx context.Context, email string) (*models.PlexUser, error)

	// Plex User Invite operations
	AssociatePlexUserWithInviteCode(ctx context.Context, userID, inviteCodeID int, expiresAt *time.Time) error
	GetPlexUserInvites(ctx context.Context, userID int) ([]models.PlexUserInvite, error)
	GetUsersWithActiveInviteCode(ctx context.Context, inviteCodeID int) ([]models.PlexUser, error)
	DisableInviteCode(ctx context.Context, codeID int) error
}

type sqlDB struct {
	conn *sql.DB
}

func Init(driver, dsn string) error {
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	DB = &sqlDB{conn: conn}
	return nil
}

func (db *sqlDB) Migrate(ctx context.Context) error {
	_, err := db.conn.ExecContext(ctx, `
    CREATE TABLE IF NOT EXISTS plex_tokens (
        user_id       INT PRIMARY KEY,
        access_token  TEXT NOT NULL,
        created_at	  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at 	  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`)
	if err != nil {
		return err
	}

	_, err = db.conn.ExecContext(ctx, `
    CREATE TABLE IF NOT EXISTS plex_users (
        id          INT PRIMARY KEY,
        uuid        TEXT NOT NULL UNIQUE,
        username    TEXT NOT NULL,
        email       TEXT NOT NULL UNIQUE,
        is_admin    BOOLEAN NOT NULL DEFAULT FALSE,
        created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`)
	if err != nil {
		return err
	}

	_, err = db.conn.ExecContext(ctx, `
    CREATE TABLE IF NOT EXISTS invite_codes (
        id                SERIAL PRIMARY KEY,
        code              TEXT NOT NULL UNIQUE,
        is_disabled       BOOLEAN NOT NULL DEFAULT FALSE,
        max_uses          INT NULL,
        used_count        INT NOT NULL DEFAULT 0,
        entitlement_name  TEXT NOT NULL,
        duration     	  TIMESTAMP NULL,
        expires_at        TIMESTAMP NULL,
        created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`)
	if err != nil {
		return err
	}

	_, err = db.conn.ExecContext(ctx, `
    CREATE TABLE IF NOT EXISTS plex_user_invites (
        id              SERIAL PRIMARY KEY,
        user_id         INT NOT NULL,
        invite_code_id  INT NOT NULL,
        used_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at      TIMESTAMP,
        CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES plex_users(id) ON DELETE CASCADE,
        CONSTRAINT fk_invite_code_id FOREIGN KEY (invite_code_id) REFERENCES invite_codes(id) ON DELETE CASCADE,
        CONSTRAINT unique_user_invite UNIQUE (user_id, invite_code_id)
    );`)
	return err
}

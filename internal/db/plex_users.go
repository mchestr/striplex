package db

import (
	"context"
	"database/sql"
	"plefi/internal/models"
)

func (db *sqlDB) SavePlexUser(ctx context.Context, user models.PlexUser) error {
	_, err := db.conn.ExecContext(ctx, `
    INSERT INTO plex_users(id, uuid, username, email, is_admin)
    VALUES($1, $2, $3, $4, $5)
    ON CONFLICT(id) DO UPDATE SET
        uuid = EXCLUDED.uuid,
        username = EXCLUDED.username,
        email = EXCLUDED.email,
        is_admin = EXCLUDED.is_admin,
        updated_at = CURRENT_TIMESTAMP;`,
		user.ID, user.UUID, user.Username, user.Email, user.IsAdmin,
	)
	return err
}

func (db *sqlDB) GetPlexUser(ctx context.Context, userID int) (*models.PlexUser, error) {
	user := &models.PlexUser{}
	err := db.conn.QueryRowContext(ctx, `
        SELECT id, uuid, username, email, is_admin, created_at, updated_at
        FROM plex_users
        WHERE id = $1`,
		userID).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email,
		&user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (db *sqlDB) GetPlexUserByEmail(ctx context.Context, email string) (*models.PlexUser, error) {
	user := &models.PlexUser{}
	err := db.conn.QueryRowContext(ctx, `
        SELECT id, uuid, username, email, is_admin, created_at, updated_at
        FROM plex_users
        WHERE email = $1`,
		email).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email,
		&user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

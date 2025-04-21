package db

import (
	"context"
	"plefi/internal/models"
)

func (db *sqlDB) SavePlexToken(ctx context.Context, tok models.PlexToken) error {
	_, err := db.conn.ExecContext(ctx, `
    INSERT INTO plex_tokens(user_id, access_token)
    VALUES($1,$2)
    ON CONFLICT(user_id) DO UPDATE SET
        access_token  = EXCLUDED.access_token,
        updated_at    = CURRENT_TIMESTAMP;`,
		tok.UserID, tok.AccessToken,
	)
	return err
}

func (db *sqlDB) GetPlexToken(ctx context.Context, userID int) (*models.PlexToken, error) {
	var tok models.PlexToken
	row := db.conn.QueryRowContext(ctx, `
        SELECT user_id, access_token, created_at, updated_at
          FROM plex_tokens
         WHERE user_id = $1`, userID)
	err := row.Scan(&tok.UserID, &tok.AccessToken, &tok.CreatedAt, &tok.UpdatedAt)
	return &tok, err
}

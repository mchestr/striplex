package db

import (
	"context"
	"database/sql"
	"plefi/internal/models"
)

func (db *sqlDB) SavePlexUser(ctx context.Context, user models.PlexUser) error {
	_, err := db.conn.ExecContext(ctx, `
    INSERT INTO plex_users(id, uuid, username, email, is_admin, notes)
    VALUES($1, $2, $3, $4, $5, $6)
    ON CONFLICT(id) DO UPDATE SET
        uuid = $2,
        username = EXCLUDED.username,
        email = EXCLUDED.email,
        is_admin = EXCLUDED.is_admin,
        notes = EXCLUDED.notes,
        updated_at = CURRENT_TIMESTAMP;`,
		user.ID, user.UUID, user.Username, user.Email, user.IsAdmin, user.Notes,
	)
	return err
}

func (db *sqlDB) GetPlexUser(ctx context.Context, userID int) (*models.PlexUser, error) {
	user := &models.PlexUser{}
	var notes sql.NullString
	err := db.conn.QueryRowContext(ctx, `
        SELECT id, uuid, username, email, is_admin, notes, created_at, updated_at
        FROM plex_users
        WHERE id = $1`,
		userID).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email,
		&user.IsAdmin, &notes, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Convert NullString to *string
	if notes.Valid {
		user.Notes = notes.String
	}
	return user, nil
}

func (db *sqlDB) GetPlexUserByEmail(ctx context.Context, email string) (*models.PlexUser, error) {
	user := &models.PlexUser{}
	var notes sql.NullString
	err := db.conn.QueryRowContext(ctx, `
        SELECT id, uuid, username, email, is_admin, notes, created_at, updated_at
        FROM plex_users
        WHERE email = $1`,
		email).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email,
		&user.IsAdmin, &notes, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Convert NullString to *string
	if notes.Valid {
		user.Notes = notes.String
	}
	return user, nil
}

func (db *sqlDB) GetAllPlexUsers(ctx context.Context) ([]models.PlexUser, error) {
	rows, err := db.conn.QueryContext(ctx, `
        SELECT id, uuid, username, email, is_admin, notes, created_at, updated_at
        FROM plex_users
        ORDER BY username ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.PlexUser
	for rows.Next() {
		var user models.PlexUser
		var notes sql.NullString
		err := rows.Scan(
			&user.ID, &user.UUID, &user.Username, &user.Email,
			&user.IsAdmin, &notes, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// Convert NullString to *string
		if notes.Valid {
			user.Notes = notes.String
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (db *sqlDB) DeletePlexUser(ctx context.Context, userID int) error {
	// First delete user-related records in dependent tables
	_, err := db.conn.ExecContext(ctx, `
		DELETE FROM plex_user_invites WHERE user_id = $1;`, userID)
	if err != nil {
		return err
	}

	_, err = db.conn.ExecContext(ctx, `
		DELETE FROM plex_tokens WHERE user_id = $1;`, userID)
	if err != nil {
		return err
	}

	// Finally delete the user record
	_, err = db.conn.ExecContext(ctx, `
		DELETE FROM plex_users WHERE id = $1;`, userID)
	return err
}

// UpdateUserNotes updates the notes for a user
func (db *sqlDB) UpdateUserNotes(ctx context.Context, userID int, notes string) error {
	_, err := db.conn.ExecContext(ctx, `
		UPDATE plex_users
		SET notes = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`,
		userID, notes)
	return err
}

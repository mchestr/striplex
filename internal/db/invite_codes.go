package db

import (
	"context"
	"database/sql"
	"plefi/internal/models"
)

// SaveInviteCode adds a new invite code to the database
func (db *sqlDB) SaveInviteCode(ctx context.Context, inviteCode models.InviteCode) (int, error) {
	var id int
	err := db.conn.QueryRowContext(ctx, `
		INSERT INTO invite_codes 
		(code, expires_at, max_uses, is_disabled, entitlement_name, duration)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, inviteCode.Code, inviteCode.ExpiresAt, inviteCode.MaxUses, inviteCode.IsDisabled,
		inviteCode.EntitlementName, inviteCode.Duration,
	).Scan(&id)

	return id, err
}

// GetInviteCodeByCode retrieves an invite code by its code value
func (db *sqlDB) GetInviteCode(ctx context.Context, id int) (*models.InviteCode, error) {
	inviteCode := &models.InviteCode{}

	err := db.conn.QueryRowContext(ctx, `
		SELECT id, code, created_at, updated_at, 
		       expires_at, max_uses, used_count, is_disabled, 
		       entitlement_name, duration
		FROM invite_codes
		WHERE id = $1
	`, id).Scan(
		&inviteCode.ID, &inviteCode.Code,
		&inviteCode.CreatedAt, &inviteCode.UpdatedAt, &inviteCode.ExpiresAt,
		&inviteCode.MaxUses, &inviteCode.UsedCount, &inviteCode.IsDisabled,
		&inviteCode.EntitlementName, &inviteCode.Duration,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return inviteCode, err
}

// GetInviteCodeByCode retrieves an invite code by its code value
func (db *sqlDB) GetInviteCodeByCode(ctx context.Context, code string) (*models.InviteCode, error) {
	inviteCode := &models.InviteCode{}

	err := db.conn.QueryRowContext(ctx, `
		SELECT id, code, created_at, updated_at, 
		       expires_at, max_uses, used_count, is_disabled, 
		       entitlement_name, duration
		FROM invite_codes
		WHERE code = $1
	`, code).Scan(
		&inviteCode.ID, &inviteCode.Code,
		&inviteCode.CreatedAt, &inviteCode.UpdatedAt, &inviteCode.ExpiresAt,
		&inviteCode.MaxUses, &inviteCode.UsedCount, &inviteCode.IsDisabled,
		&inviteCode.EntitlementName, &inviteCode.Duration,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return inviteCode, err
}

// UpdateInviteCodeUsage increments the usage count for an invite code
func (db *sqlDB) UpdateInviteCodeUsage(ctx context.Context, codeID int) error {
	_, err := db.conn.ExecContext(ctx, `
		UPDATE invite_codes
		SET used_count = used_count + 1, updated_at = NOW()
		WHERE id = $1
	`, codeID)

	return err
}

// ListActiveInviteCodes retrieves all active invite codes
func (db *sqlDB) ListActiveInviteCodes(ctx context.Context) ([]models.InviteCode, error) {
	rows, err := db.conn.QueryContext(ctx, `
		SELECT id, code, created_at, updated_at, 
		       expires_at, max_uses, used_count, is_disabled, 
		       entitlement_name, duration
		FROM invite_codes
		WHERE is_disabled = FALSE
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inviteCodes []models.InviteCode
	for rows.Next() {
		var code models.InviteCode
		if err := rows.Scan(
			&code.ID, &code.Code,
			&code.CreatedAt, &code.UpdatedAt, &code.ExpiresAt,
			&code.MaxUses, &code.UsedCount, &code.IsDisabled,
			&code.EntitlementName, &code.Duration,
		); err != nil {
			return nil, err
		}
		inviteCodes = append(inviteCodes, code)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return inviteCodes, nil
}

func (db *sqlDB) AssociatePlexUserWithInviteCode(ctx context.Context, userID, inviteCodeID int) error {
	_, err := db.conn.ExecContext(ctx, `
    INSERT INTO plex_user_invites(user_id, invite_code_id)
    VALUES($1, $2)
    ON CONFLICT(user_id, invite_code_id) DO UPDATE SET 
        used_at = CURRENT_TIMESTAMP;`,
		userID, inviteCodeID,
	)
	return err
}

func (db *sqlDB) GetPlexUserInvites(ctx context.Context, userID int) ([]models.PlexUserInvite, error) {
	rows, err := db.conn.QueryContext(ctx, `
        SELECT pui.id, pui.user_id, pui.invite_code_id, pui.used_at, ic.code, ic.entitlement_name
        FROM plex_user_invites pui
        JOIN invite_codes ic ON pui.invite_code_id = ic.id
        WHERE pui.user_id = $1
        ORDER BY pui.used_at DESC`,
		userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []models.PlexUserInvite
	for rows.Next() {
		var invite models.PlexUserInvite
		err := rows.Scan(
			&invite.ID, &invite.UserID, &invite.InviteCodeID,
			&invite.UsedAt, &invite.InviteCode, &invite.EntitlementName,
		)
		if err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}

	return invites, rows.Err()
}

func (db *sqlDB) GetUsersWithActiveInviteCode(ctx context.Context, inviteCodeID int) ([]models.PlexUser, error) {
	rows, err := db.conn.QueryContext(ctx, `
        SELECT u.id, u.uuid, u.username, u.email, u.is_admin, u.created_at, u.updated_at
        FROM plex_users u
        JOIN plex_user_invites pui ON u.id = pui.user_id
        WHERE pui.invite_code_id = $1
        ORDER BY u.username ASC`,
		inviteCodeID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.PlexUser
	for rows.Next() {
		var user models.PlexUser
		err := rows.Scan(
			&user.ID, &user.UUID, &user.Username, &user.Email,
			&user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// DisableInviteCode marks an invite code as disabled
func (db *sqlDB) DisableInviteCode(ctx context.Context, codeID int) error {
	_, err := db.conn.ExecContext(ctx, `
		UPDATE invite_codes
		SET is_disabled = TRUE, updated_at = NOW()
		WHERE id = $1
	`, codeID)

	return err
}

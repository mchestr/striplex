CREATE TABLE IF NOT EXISTS plex_tokens (
    user_id       INT PRIMARY KEY,
    access_token  TEXT NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS plex_users (
    id          INT PRIMARY KEY,
    uuid        TEXT NOT NULL UNIQUE,
    username    TEXT NOT NULL,
    email       TEXT NOT NULL UNIQUE,
    is_admin    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS invite_codes (
    id                SERIAL PRIMARY KEY,
    code              TEXT NOT NULL UNIQUE,
    is_disabled       BOOLEAN NOT NULL DEFAULT FALSE,
    max_uses          INT NULL,
    used_count        INT NOT NULL DEFAULT 0,
    entitlement_name  TEXT NOT NULL,
    duration          TIMESTAMP NULL,
    expires_at        TIMESTAMP NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS plex_user_invites (
    id              SERIAL PRIMARY KEY,
    user_id         INT NOT NULL,
    invite_code_id  INT NOT NULL,
    used_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES plex_users(id) ON DELETE CASCADE,
    CONSTRAINT fk_invite_code_id FOREIGN KEY (invite_code_id) REFERENCES invite_codes(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_invite UNIQUE (user_id, invite_code_id)
);

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
	id            BIGSERIAL PRIMARY KEY,
	created_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	name          TEXT                        NOT NULL,
	email         CITEXT UNIQUE               NOT NULL,
	password_hash BYTEA                       NOT NULL,
	activated     BOOL                        NOT NULL DEFAULT FALSE,
	version       INTEGER                     NOT NULL DEFAULT 1
);
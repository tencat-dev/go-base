CREATE TABLE users
(
    id            UUID PRIMARY KEY     DEFAULT uuidv7(),

    name          TEXT        NOT NULL,
    email         TEXT        NOT NULL UNIQUE,

    password      TEXT        NOT NULL,
    password_salt TEXT        NOT NULL,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

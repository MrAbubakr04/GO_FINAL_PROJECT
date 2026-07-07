CREATE TABLE users (
    id BIGSERIAL,
    login TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'employee',
    active_from TIMESTAMP NOT NULL DEFAULT NOW(),
    active_to TIMESTAMP,
    dt_created TIMESTAMP NOT NULL DEFAULT NOW(),
    dt_updated TIMESTAMP,

    CONSTRAINT users_id_pk PRIMARY KEY (id)
);

CREATE UNIQUE INDEX users_login_un
ON users (login)
WHERE active_to IS NULL;

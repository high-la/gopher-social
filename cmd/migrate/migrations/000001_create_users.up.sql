-- Run this command from terminal
-- sudo -u postgres psql -d gopher-social -c "CREATE EXTENSION IF NOT EXISTS citext;"

-- to run bellow commands run this on terminal
-- migrate -path ./cmd/migrate/migrations -database "$DSN" up

BEGIN;

CREATE TABLE IF NOT EXISTS users (

    id bigserial PRIMARY KEY,
    username varchar(20) UNIQUE NOT NULL,
    email citext UNIQUE NOT NULL,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);


-- If this fails, the table above isn't created and the version stays clean
COMMIT;
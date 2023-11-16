CREATE TABLE IF NOT EXISTS "users"
(
    "id"       bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email"    varchar NOT NULL,
    "password" varchar NOT NULL
);


CREATE TABLE IF NOT EXISTS "user_token"
(
    "id"            bigserial PRIMARY KEY,
    "token"         varchar    NOT NULL,
    "refresh_token" varchar    NOT NULL,
    "user_id"       int UNIQUE NOT NULL,
    "created_at"    timestamp  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"    timestamp  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "users"
(
    "id"           serial PRIMARY KEY,
    "username"     varchar        NOT NULL,
    "email"        varchar unique NOT NULL,
    "password"     varchar        NOT NULL,
    "role"         varchar        not null,
    "is_confirmed" boolean        not null
);

CREATE table if not exists "user_code"
(
    "id"        bigserial PRIMARY KEY,
    "user_id"   integer,
    "user_code" varchar not null,

    FOREIGN KEY ("user_id") REFERENCES users ("id")
);


CREATE TABLE IF NOT EXISTS "user_token"
(
    "id"            bigserial PRIMARY KEY,
    "token"         varchar    NOT NULL,
    "refresh_token" varchar    NOT NULL,
    "user_id"       int UNIQUE NOT NULL,
    "created_at"    timestamp  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"    timestamp  NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY ("user_id") REFERENCES users ("id")
);

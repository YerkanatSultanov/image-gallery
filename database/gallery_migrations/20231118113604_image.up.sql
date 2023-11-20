create table if not exists image
(
    "id"          bigserial PRIMARY KEY,
    "user_id"     integer not null ,
    "description" varchar,
    "image_link"  varchar   not null,
    "created_at"  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
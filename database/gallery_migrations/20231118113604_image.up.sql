create table if not exists image
(
    "id"          bigserial PRIMARY KEY,
    "user_id"     integer   not null,
    "description" varchar,
    "image_link"  varchar   not null,
    "created_at"  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create table if not exists tags
(
    "tag_id"   bigserial PRIMARY KEY,
    "tag_name" varchar not null
);

CREATE TABLE IF NOT EXISTS tag_images
(
    "image_id" INTEGER,
    "tag_id"   INTEGER,
    FOREIGN KEY ("image_id") REFERENCES image ("id"),
    FOREIGN KEY ("tag_id") REFERENCES tags ("tag_id")
);
create table if not exists image
(
    "id"          bigserial PRIMARY KEY,
    "user_id"     integer,
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

CREATE TABLE IF NOT EXISTS followers
(
    "follower_id"       INTEGER      NOT NULL,
    "followee_id"       INTEGER      NOT NULL,
    PRIMARY KEY ("follower_id", "followee_id")
);

CREATE TABLE if not exists likes (
    "user_id" INTEGER,
    "image_id" INTEGER,
    "created_at"  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("image_id") REFERENCES image ("id")
);
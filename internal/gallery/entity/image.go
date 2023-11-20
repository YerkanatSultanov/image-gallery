package entity

import "time"

type Photo struct {
	Id          int       `db:"id"`
	UserId      int       `db:"userId"`
	Description string    `db:"name"`
	ImageLink   string    `db:"image"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type ImageRequest struct {
	Description string `db:"name"`
	ImageLink   string `db:"image"`
}

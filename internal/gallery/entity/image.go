package entity

import "time"

type Image struct {
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

type PhotoResponse struct {
	Id          int       `db:"id"`
	UserId      int       `db:"userId"`
	Description string    `db:"name"`
	ImageLink   string    `db:"image"`
	CreatedAt   time.Time `db:"created_at"`
}

type Tags struct {
	TagId   int    `db:"tag_id"`
	TagName string `db:"tag_name"`
}

type TageRequest struct {
	TagName string `json:"TagName"`
	ImageId int    `json:"ImageId"`
}

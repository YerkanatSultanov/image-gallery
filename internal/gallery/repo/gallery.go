package repo

import (
	"context"
	"database/sql"
	"fmt"
	"image-gallery/internal/gallery/entity"
	"log"
	"time"
)

func (r *Repository) CreatePhoto(ph *entity.Photo) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var lastInsertId int
	query := `insert into image(user_id, description, image_link) values($1, $2, $3) returning id`
	err := r.db.QueryRowContext(c, query, ph.UserId, ph.Description, ph.ImageLink).Scan(&lastInsertId)

	if err != nil {
		return fmt.Errorf("query bake failed: %w", err)
	}

	ph.Id = (lastInsertId)
	return nil
}

func (r *Repository) GetAllPhotos() ([]*entity.Photo, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, user_id, description, image_link, created_at from image"
	rows, err := r.db.QueryContext(c, query)
	if err != nil {
		log.Fatalf("Error in database: %s", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var photos []*entity.Photo

	for rows.Next() {
		var photo entity.Photo
		if err := rows.Scan(&photo.Id, &photo.UserId, &photo.Description, &photo.ImageLink, &photo.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}

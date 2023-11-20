package repo

import (
	"context"
	"fmt"
	"image-gallery/internal/gallery/entity"
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

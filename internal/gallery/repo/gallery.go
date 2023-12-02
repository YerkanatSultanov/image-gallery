package repo

import (
	"context"
	"database/sql"
	"fmt"
	"image-gallery/internal/gallery/entity"
	"log"
	"time"
)

func (r *Repository) CreatePhoto(ph *entity.Image) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var lastInsertId int
	query := `insert into image(user_id, description, image_link) values($1, $2, $3) returning id`
	err := r.db.QueryRowContext(c, query, ph.UserId, ph.Description, ph.ImageLink).Scan(&lastInsertId)

	if err != nil {
		return fmt.Errorf("query bake failed: %w", err)
	}

	ph.Id = lastInsertId
	return nil
}

func (r *Repository) GetAllPhotos() ([]*entity.Image, error) {
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

	var photos []*entity.Image

	for rows.Next() {
		var photo entity.Image
		if err := rows.Scan(&photo.Id, &photo.UserId, &photo.Description, &photo.ImageLink, &photo.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("error in rows: %s", err)
		return nil, err
	}

	return photos, nil
}

func (r *Repository) GetGalleryById(id int) ([]*entity.PhotoResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "Select id, user_id, description, image_link, created_at from image where user_id = $1"
	rows, err := r.db.QueryContext(c, query, id)
	if err != nil {
		log.Fatalf("Error in database: %s", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var photos []*entity.PhotoResponse

	for rows.Next() {
		var photo entity.PhotoResponse
		if err := rows.Scan(&photo.Id, &photo.UserId, &photo.Description, &photo.ImageLink, &photo.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("error in rows: %s", err)
		return nil, err
	}

	return photos, nil
}

func (r *Repository) AddTagName(t *entity.Tags) (int, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var lastInsertId int
	query := "insert into tags(tag_name) values($1) returning tag_id"

	err := r.db.QueryRowContext(c, query, &t.TagName).Scan(&lastInsertId)

	if err != nil {
		return 0, fmt.Errorf("query bake failed: %v", err)
	}

	t.TagId = lastInsertId
	return t.TagId, nil
}

func (r *Repository) AddTagImage(tagId, imageId int) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO tag_images(image_id, tag_id) VALUES ($1, $2)"

	_, err := r.db.ExecContext(c, query, imageId, tagId)
	if err != nil {
		return fmt.Errorf("query bake failed: %v", err)
	}

	return nil
}

func (r *Repository) BeginTransaction() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *Repository) DeleteImage(imageId int) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "delete from image where id = $1"
	_, err := r.db.ExecContext(c, query, imageId)

	if err != nil {
		return fmt.Errorf("query bake failed: %v", err)
	}

	return nil
}

func (r *Repository) Follow(followerId int, followeeId int) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "insert into followers(follower_id, followee_id) values ($1, $2)"

	_, err := r.db.ExecContext(c, query, followerId, followeeId)

	if err != nil {
		return fmt.Errorf("error in exec query: %s", err)
	}

	return nil

}

func (r *Repository) UserLikedPhoto(userId, imageID int) (bool, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	query := "SELECT COUNT(*) FROM likes WHERE user_id = $1 and image_id = $2"
	var count int
	err := r.db.QueryRowContext(c, query, userId, imageID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Error checking if user liked photo: %s\n", err)
	}
	return count > 0, nil
}

func (r *Repository) LikePhoto(like *entity.Likes) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO likes (user_id, image_id, created_at) VALUES ($1, $2, NOW())"
	_, err := r.db.ExecContext(c, query, like.UserId, like.ImageId)
	if err != nil {
		return fmt.Errorf("Error liking photo: %s\n", err)
	}
	return nil
}

func (r *Repository) UserHasImage(imageID int) (bool, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM image WHERE id = $1"
	var count int
	err := r.db.QueryRowContext(c, query, imageID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Error checking if user has image: %s\n", err)
	}
	return count > 0, nil
}

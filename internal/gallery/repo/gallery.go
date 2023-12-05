package repo

import (
	"context"
	"database/sql"
	"fmt"
	"image-gallery/internal/gallery/entity"
	"image-gallery/pkg/util"
	"log"
	"time"
)

func (r *Repository) CreatePhoto(ph entity.Image) error {
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
			return
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
			return
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

	err := r.DeleteTagFromImage(imageId)
	if err != nil {
		return fmt.Errorf("can not delete tag from image: %s", err)
	}

	err = r.DeleteLikeFromImage(imageId)
	if err != nil {
		return fmt.Errorf("can not delete like from image: %s", err)
	}

	query := "DELETE FROM image WHERE id = $1"
	_, err = r.db.ExecContext(c, query, imageId)

	if err != nil {
		return fmt.Errorf("error in exec query: %s", err)
	}

	return nil
}

func (r *Repository) DeleteTagFromImage(imageId int) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "DELETE FROM tag_images WHERE image_id = $1"
	_, err := r.db.ExecContext(c, query, imageId)

	if err != nil {
		return fmt.Errorf("error in exec query: %s", err)
	}

	return nil
}

func (r *Repository) DeleteLikeFromImage(imageId int) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "DELETE FROM likes WHERE image_id = $1"
	_, err := r.db.ExecContext(c, query, imageId)

	if err != nil {
		return fmt.Errorf("error in exec query: %s", err)
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

func (r *Repository) FindImagesByTag(tag string) ([]entity.Image, error) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT i.id, i.user_id, i.description, i.image_link, i.created_at, i.updated_at
        FROM image i
        JOIN tag_images ti ON i.id = ti.image_id
        JOIN tags t ON ti.tag_id = t.tag_id
        WHERE t.tag_name LIKE $1
    `

	rows, err := r.db.QueryContext(c, query, tag+"%")
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var images []entity.Image
	for rows.Next() {
		var img entity.Image
		if err := rows.Scan(&img.Id, &img.UserId, &img.Description, &img.ImageLink, &img.CreatedAt, &img.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		images = append(images, img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return images, nil
}

func (r *Repository) GetImages(sortKey, sortBy string) ([]*entity.Image, error) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	validSortKeys := []string{"created_at", "updated_at", "id", "user_id"}
	if !util.Contains(validSortKeys, sortKey) {
		return nil, fmt.Errorf("invalid sortKey parameter")
	}

	validSortByValues := []string{"asc", "desc"}
	if !util.Contains(validSortByValues, sortBy) {
		return nil, fmt.Errorf("invalid sortBy parameter")
	}

	tableName := " image"

	query := fmt.Sprintf(`
		SELECT id, user_id, description, image_link, created_at, updated_at
		FROM%s 
		ORDER BY %s %s;
    `, tableName, sortKey, sortBy)

	rows, err := r.db.QueryContext(c, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var images []*entity.Image
	for rows.Next() {
		var img entity.Image
		if err := rows.Scan(&img.Id, &img.UserId, &img.ImageLink, &img.Description, &img.CreatedAt, &img.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		images = append(images, &img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return images, nil
}

func (r *Repository) GetImagesForFollower(followerId int, followeeId int) ([]*entity.Image, error) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT i.id, i.user_id, i.description, i.image_link, i.created_at, i.updated_at
		FROM image i
		JOIN followers f ON i.user_id = f.followee_id
		WHERE f.follower_id = $1 AND f.followee_id = $2
		ORDER BY i.created_at DESC
	`

	fmt.Printf("Executing query: %s with followerID: %d, followeeID: %d\n", query, followerId, followeeId)

	rows, err := r.db.QueryContext(c, query, followerId, followeeId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var images []*entity.Image
	for rows.Next() {
		var img entity.Image
		if err := rows.Scan(&img.Id, &img.UserId, &img.Description, &img.ImageLink, &img.CreatedAt, &img.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		images = append(images, &img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %v", err)
	}

	return images, nil
}

func (r *Repository) GetLikedImages(userId int) ([]*entity.Image, error) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		SELECT i.id, i.user_id, i.description, i.image_link, i.created_at, i.updated_at
		FROM image i
		JOIN likes l ON i.id = l.image_id
		WHERE l.user_id = $1
	`

	rows, err := r.db.QueryContext(c, query, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var images []*entity.Image
	for rows.Next() {
		var img entity.Image
		if err := rows.Scan(&img.Id, &img.UserId, &img.Description, &img.ImageLink, &img.CreatedAt, &img.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		images = append(images, &img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %v", err)
	}

	return images, nil
}

func (r *Repository) UpdateImage(imageId, userId int, description string) (*entity.Image, error) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE image SET description = $1 WHERE user_id = $2 AND id = $3 
                                  RETURNING id, user_id, description, image_link, created_at, updated_at`

	var updatedImage entity.Image
	err := r.db.QueryRowContext(c, query, description, userId, imageId).
		Scan(&updatedImage.Id, &updatedImage.UserId, &updatedImage.Description, &updatedImage.ImageLink, &updatedImage.CreatedAt, &updatedImage.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update image: %v", err)
	}

	return &updatedImage, nil

}

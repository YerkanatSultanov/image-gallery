package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"image-gallery/internal/gallery/config"
)

type DataBase struct {
	db *sql.DB
}

type Config config.Database

func (c Config) dsn() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SslMode,
	)
}

func NewDataBase(cfg config.Database) (*DataBase, error) {
	conf := Config(cfg)

	db, err := sql.Open(conf.Driver, conf.dsn())
	if err != nil {
		return nil, err
	}

	return &DataBase{db: db}, nil
}

func (d *DataBase) Close() {
	err := d.db.Close()
	if err != nil {
		return
	}
}

func (d *DataBase) GetDB() *sql.DB {
	return d.db
}

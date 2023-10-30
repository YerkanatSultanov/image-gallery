package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type DataBase struct {
	db *sql.DB
}

func NewDataBase() (*DataBase, error) {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/authDatabase?sslmode=disable")
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

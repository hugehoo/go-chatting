package repository

import (
	"chat_controller/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	cfg *config.Config
	db  *sql.DB
}

func NewRepository(cfg *config.Config) (*Repository, error) {

	r := &Repository{
		cfg: cfg,
	}
	var err error

	if r.db, err = sql.Open(cfg.DB.Database, cfg.DB.Url); err != nil {
		return nil, err
	}
	return r, nil
}

package repository

import (
	"chat_backend/config"
	"chat_backend/repository/kafka"
	"database/sql"
)

type Repository struct {
	cfg   *config.Config
	db    *sql.DB
	Kafka *kafka.Kafka
}

const (
	room       = "chatting.room"
	chat       = "chatting.room"
	serverInfo = "chatting.server_info"
)

func (r *Repository) ServerSet(ip string, available bool) error {
	// upsert
	_, err := r.db.Exec("INSERT server_info(`ìp`, `available`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `available` = VALUES(`available`)",
		ip, available)
	return err
}

func NewRepository(c *config.Config) (*Repository, error) {
	r := &Repository{cfg: c}
	var err error

	if r.db, err = sql.Open(c.DB.Database, c.DB.URL); err != nil {
		return nil, err
	} else if r.Kafka, err = kafka.NewKafka(c); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

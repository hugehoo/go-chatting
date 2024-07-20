package repository

import (
	"database/sql"
	"websocket/config"
)

type Repository struct {
	cfg *config.Config
	db  *sql.DB
}

const (
	room       = "chatting.room"
	chat       = "chatting.room"
	serverInfo = "chatting.server_info"
)

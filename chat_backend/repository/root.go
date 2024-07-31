package repository

import (
	"chat_backend/config"
	messageBroker "chat_backend/repository/kafka"
	"chat_backend/types/schema"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"strings"
)

type Repository struct {
	cfg   *config.Config
	db    *sql.DB
	Kafka *messageBroker.Kafka
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
	} else if r.Kafka, err = messageBroker.NewKafka(c); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (r *Repository) InsertChatting(user, message, roomName string) error {

	// 2. kafka 로 던진 후 컨슈머 서버에서 db insert
	ch := make(chan kafka.Event)
	v, _ := json.Marshal(message)
	r.Kafka.PublishEvent("chat-message", v, ch)

	// [original] 1. chat server 에서 직접 db insert
	_, err := r.db.Exec("INSERT INTO chatting.chat(room, name, message) VALUES (?, ?, ?)",
		roomName, user, message)
	return err
}

func (r *Repository) GetChatList(roomName string) ([]*schema.Chat, error) {
	qs := query([]string{"SELECT * FROM", chat, "WHERE room = ? ORDER BY `when` DESC LIMIT 10"})
	if cursor, err := r.db.Query(qs, roomName); err != nil {
		return nil, err
	} else {
		defer cursor.Close()
		var result []*schema.Chat
		for cursor.Next() {
			c := new(schema.Chat)
			if err = cursor.Scan(&c.ID, &c.Room, &c.Name, &c.Message, &c.When); err != nil {
				return nil, err
			} else {
				result = append(result, c)
			}
		}
		return result, err // 아래랑 같은 것 아님..?
		//if len(result) == 0 {
		//	return []*schema.Chat{}, err
		//} else {
		//	return result, err
		//}
	}
}

func (r *Repository) RoomList() ([]*schema.Room, error) {
	qs := query([]string{"SELECT * FROM", room})
	if cursor, err := r.db.Query(qs); err != nil {
		return nil, err
	} else {
		defer cursor.Close()
		var result []*schema.Room

		for cursor.Next() {
			d := new(schema.Room)
			if err = cursor.Scan(&d.ID, &d.Name, &d.CreateAt, &d.UpdatedAt); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}
		return result, err
	}
}

func (r *Repository) MakeRoom(name string) error {
	_, err := r.db.Exec("INSERT INTO chatting.room(name) VALUES (?)", name)
	return err
}

func (r *Repository) Room(name string) (*schema.Room, error) {
	d := new(schema.Room)
	qs := query([]string{"SELECT * FROM", room, "WHERE name = ?"})

	err := r.db.QueryRow(qs, name).Scan(&d.ID, &d.Name, &d.CreateAt, &d.UpdatedAt)
	if err = noResult(err); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}

func noResult(err error) error {
	if errors.Is(sql.ErrNoRows, err) {
		return nil
	}
	return err
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}

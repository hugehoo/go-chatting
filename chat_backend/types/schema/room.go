package schema

import "time"

type Room struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreateAt  time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

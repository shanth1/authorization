package domain

import "time"

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
}

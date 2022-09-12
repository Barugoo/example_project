package models

import "time"

type Order struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	ItemIDs   []int64   `json:"item_ids"`
	CreatedAt time.Time `json:"created_at"`
}

type Item struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

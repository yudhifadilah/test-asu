package models

import (
	"time"
)

// Article struct
type Article struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Excerpt string    `json:"excerpt"`
	Content string    `json:"content"`
	Image   string    `json:"image"`
	RegDate time.Time `json:"reg_date"`
}

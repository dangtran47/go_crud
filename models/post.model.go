package models

import "time"

type Post struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Title     string    `gorm:"uniqueIndex;not null" json:"title,omitempty" binding:"required"`
	Content   string    `gorm:"not null" json:"content,omitempty"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	ID        uint      `gorm:"primary_key"`
	AuthorID  uint      `gorm:"not null"`
}

type CreatePost struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
}

type UpdatePost struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

type PostResponse struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ID        uint      `json:"id"`
	AuthorID  uint      `json:"author_id"`
}

func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

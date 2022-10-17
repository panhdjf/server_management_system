package models

import (
	"time"
	// "github.com/google/uuid"
)

type Server struct {
	// ID          string    `gorm:"uniqueIndex;default:uuid_generate_v4();primary_key" json:"server_id,omitempty"`
	ID          string    `gorm:"uniqueIndex;primary_key" json:"server_id,omitempty"`
	Name        string    `gorm:"uniqueIndex;not null" json:"server_name,omitempty"`
	Status      string    `gorm:"not null" json:"status,omitempty"`
	Ipv4        string    `gorm:"not null" json:"ipv4,omitempty"`
	User        int       `gorm:"not null" json:"user,omitempty"`
	CreatedTime time.Time `gorm:"not null" json:"created_time,omitempty"`
	LastUpdated time.Time `gorm:"not null" json:"last_updated,omitempty"`
}

type CreateServerRequest struct {
	Name        string    `json:"server_name" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	Ipv4        string    `json:"ipv4" binding:"required"`
	User        string    `json:"user,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
	LastUpdated time.Time `json:"last_updated,omitempty"`
}

type UpdateServer struct {
	Name        string    `json:"server_name,omitempty"`
	Status      string    `json:"status,omitempty"`
	Ipv4        string    `json:"ipv4,omitempty"`
	User        string    `json:"user,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
	LastUpdated time.Time `json:"last_updated,omitempty"`
}

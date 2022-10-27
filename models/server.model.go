package models

import (
	"time"
	// "github.com/google/uuid"
)

type Server struct {
	// ID          string    `gorm:"uniqueIndex;default:uuid_generate_v4();primary_key" json:"server_id,omitempty"`
	ID            string    `gorm:"uniqueIndex;primary_key" json:"server_id,omitempty"`
	Name          string    `gorm:"uniqueIndex;not null" json:"server_name,omitempty"`
	Status        string    `gorm:"not null" json:"status,omitempty"`
	Uptime        float64   `gorm:"not null" json:"uptime,omitempty"`
	Ipv4          string    `gorm:"not null" json:"ipv4,omitempty"`
	IdUserManager int       `gorm:"not null" json:"idUserManager,omitempty"`
	CreatedTime   time.Time `gorm:"not null" json:"created_time,omitempty"`
	LastUpdated   time.Time `gorm:"not null" json:"last_updated,omitempty"`
}

type CreateServerRequest struct {
	Name          string    `json:"server_name" binding:"required"`
	Status        string    `json:"status" binding:"required"`
	Uptime        float64   `json:"uptime" binding:"required"`
	Ipv4          string    `json:"ipv4" binding:"required"`
	IdUserManager string    `json:"idUserManager,omitempty"`
	CreatedTime   time.Time `json:"created_time,omitempty"`
	LastUpdated   time.Time `json:"last_updated,omitempty"`
}

type UpdateServer struct {
	Name          string    `json:"server_name,omitempty"`
	Status        string    `json:"status,omitempty"`
	Uptime        float64   `json:"uptime,omitempty"`
	Ipv4          string    `json:"ipv4,omitempty"`
	IdUserManager string    `json:"idUserManager,omitempty"`
	CreatedTime   time.Time `json:"created_time,omitempty"`
	LastUpdated   time.Time `json:"last_updated,omitempty"`
}

type ImportExcel struct {
	ID   string ` json:"server_id,omitempty"`
	Name string `json:"server_name,omitempty"`
}

type ServerStatus struct {
	UpdateTime string `json:"uptime"`
	Status     string `json:"status"`
}

//them update time -

//lay yeu cau, phan tích thiết kế, cài đặt, kiểm thử

// tiêu đề
// slide: required - -
// model: server (inf - struct ) vẽ draw.io, chup ảnh db trong psql
// kiến trúc ứng dụng:
// các chức năng (API gắn vói chức năng nào)
// tính năng phi chức năng: token
// cảm ơn

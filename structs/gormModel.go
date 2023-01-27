package structs

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        *uint          `json:"id" query:"id" form:"id" param:"id" gorm:"primarykey;autoIncrement;not null"`
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime" query:"created_at" form:"created_at" param:"created_at"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime" query:"updated_at" form:"updated_at" param:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"default:null;index" query:"deleted_at" form:"deleted_at" param:"deleted_at"`
}

type ModelNoID struct {
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime" query:"created_at" form:"created_at" param:"created_at"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime" query:"updated_at" form:"updated_at" param:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"default:null;index" query:"deleted_at" form:"deleted_at" param:"deleted_at"`
}

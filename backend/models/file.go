package models

import (
	//"time"
	//"gorm.io/gorm"
)

type File struct {
	ID        uint           `gorm:"primaryKey" json:"id"`

	FileName  string `json:"file_name"`   
	FilePath  string `json:"file_path"`   
	FileSize  int64  `json:"file_size"`   
	FileType  string `json:"file_type"`   
	
	OwnerID   uint   `json:"owner_id"`
}
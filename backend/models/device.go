package models

import (
	"time"
	//"gorm.io/gorm"
)

type Device struct {
    ID           uint      `gorm:"primaryKey"`
    DeviceID     string    `gorm:"uniqueIndex"`
    DeviceName   string    `json:"device_name"`  
    LastSync     time.Time `json:"last_sync"`
    UserAgent    string    `json:"user_agent"` 
    IPAddress    string    `json:"ip_address"`   
}
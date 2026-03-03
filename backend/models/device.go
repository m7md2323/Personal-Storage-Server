package models

import (
	"time"
	//"gorm.io/gorm"
)

type Device struct {
    ID           uint      `gorm:"primaryKey"`
    DeviceID     string    `gorm:"uniqueIndex"` // The unique UUID string
    DeviceName   string    `json:"device_name"`  // e.g., "M7md's iPhone" or "Home PC"
    LastSync     time.Time `json:"last_sync"`
    UserAgent    string    `json:"user_agent"`   // Tells you if it's Chrome, Safari, Linux, etc.
    IPAddress    string    `json:"ip_address"`   // The current local IP
}
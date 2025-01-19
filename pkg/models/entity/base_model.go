package entity

import "time"

/**
  @author: victor2022
  @since: 2025/1/13
*/
type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GmtCreate time.Time `json:"gmt_create"`
	GmtModify time.Time `json:"gmt_modify"`
}

func GetNewBaseModel() *BaseModel {
	return &BaseModel{
		GmtCreate: time.Now(),
		GmtModify: time.Now(),
	}
}

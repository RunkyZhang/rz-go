package models

import "time"

type CallbackBasePo struct {
	Finished             bool      `gorm:"column:finished"`
	FinishedTime         time.Time `gorm:"column:finishedTime"`
	FinishedCallbackUrls string    `gorm:"column:finishedCallbackUrls"`
	States               string    `gorm:"column:states"`
	ErrorMessages        string    `gorm:"column:errorMessages"`
	ExpireTime           time.Time `gorm:"column:expireTime"`
}

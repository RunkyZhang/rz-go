package models

import "time"

type CallbackBasePo struct {
	Disable              bool      `gorm:"column:disable"`
	Finished             bool      `gorm:"column:finished"`
	FinishedTime         time.Time `gorm:"column:finishedTime"`
	FinishedCallbackUrls string    `gorm:"column:finishedCallbackUrls"`
	ExpireCallbackUrls   string    `gorm:"column:expireCallbackUrls"`
	States               string    `gorm:"column:states"`
	ErrorMessages        string    `gorm:"column:errorMessages"`
	ExpireTime           time.Time `gorm:"column:expireTime"`
	ProviderId           string    `gorm:"column:providerId"`
}

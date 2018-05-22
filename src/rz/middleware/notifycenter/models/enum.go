package models

type SendChannel int

const (
	Email    SendChannel = iota
	SmS
	QYWeixin
	Weixin
	JPush
	Voice
)

type ToType int

const (
	Auto   ToType = iota
	Phone
	UserId
	Mail
)

type QYWeixinMessageType int

const (
	Text QYWeixinMessageType = iota
	TextCard
)

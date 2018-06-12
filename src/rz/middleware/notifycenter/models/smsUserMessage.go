package models

type SmsUserCallbackMessageDto struct {
	Id             string `json:"id"`
	Content        string `json:"content"`
	Sign           string `json:"sign"`
	Time           int64  `json:"time"`
	CreatedTime    int64  `json:"createdTime"`
	Finished       bool   `json:"finished"`
	FinishedTime   int64  `json:"finishedTime"`
	ErrorMessage   string `json:"errorMessage"`
	UserCallbackId string `json:"userCallbackId"`
}

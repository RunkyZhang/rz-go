package external

type SmsUserCallbackMessageRequestExternalDto struct {
	Extend     string `json:"extend"`
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
	Sign       string `json:"sign"`
	Text       string `json:"text"`
	Time       int64  `json:"time"`
}

type SmsUserCallbackMessageResponseExternalDto struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
}
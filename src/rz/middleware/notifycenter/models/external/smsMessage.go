package external

type PhoneNumberExternalDto struct {
	Nationcode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
}
type SmsMessageRequestExternalDto struct {
	Tel    []PhoneNumberExternalDto `json:"tel"`
	Type   string                   `json:"type"`
	Msg    string                   `json:"msg"`
	Sig    string                   `json:"sig"`
	Time   int64                    `json:"time"`
	Extend string                   `json:"extend"`
	Ext    string                   `json:"ext"`
}

type SmsMessageResponseExternalDto struct {
	Result int                           `json:"result"`
	Errmsg string                        `json:"errmsg"`
	Ext    string                        `json:"ext"`
	Detail []SmsMessageDetailExternalDto `json:"detail"`
}

type SmsMessageDetailExternalDto struct {
	Result     int    `json:"result"`
	Errmsg     string `json:"errmsg"`
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
	Sid        string `json:"sid"`
	Fee        int    `json:"fee"`
}

package external

type PhoneNumberPackExternalDto struct {
	Nationcode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
}
type SmsMessageRequestExternalDto struct {
	Tel    []PhoneNumberPackExternalDto `json:"tel"`
	Msg    string                       `json:"msg,omitempty"`
	Params []string                     `json:"params"`
	Sig    string                       `json:"sig"`
	Time   int64                        `json:"time"`
	Extend string                       `json:"extend"`
	Ext    string                       `json:"ext"`
	TplId  int                          `json:"tpl_id"`
	Sign   string                       `json:"sign"`
}

type SmsMessageResponseExternalDto struct {
	Result int                                 `json:"result"`
	Errmsg string                              `json:"errmsg"`
	Ext    string                              `json:"ext"`
	Detail []SmsMessageResultDetailExternalDto `json:"detail"`
}

type SmsMessageResultDetailExternalDto struct {
	Result     int    `json:"result"`
	Errmsg     string `json:"errmsg"`
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
	Sid        string `json:"sid"`
	Fee        int    `json:"fee"`
}

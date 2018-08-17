package external

type TencentPhoneNumberPackDto struct {
	Nationcode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
}
type TencentSmsMessageRequestDto struct {
	Tel    []TencentPhoneNumberPackDto `json:"tel"`
	Msg    string                      `json:"msg,omitempty"`
	Params []string                    `json:"params"`
	Sig    string                      `json:"sig"`
	Time   int64                       `json:"time"`
	Extend string                      `json:"extend"`
	Ext    string                      `json:"ext"`
	TplId  int                         `json:"tpl_id"`
	Sign   string                      `json:"sign"`
}

type TencentSmsMessageResponseDto struct {
	Result int                                `json:"result"`
	Errmsg string                             `json:"errmsg"`
	Ext    string                             `json:"ext"`
	Detail []TencentSmsMessageResultDetailDto `json:"detail"`

	ActionStatus string `json:"ActionStatus"`
	ErrorCode    int    `json:"ErrorCode"`
	ErrorInfo    string `json:"ErrorInfo"`
}

type TencentSmsMessageResultDetailDto struct {
	Result     int    `json:"result"`
	Errmsg     string `json:"errmsg"`
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
	Sid        string `json:"sid"`
	Fee        int    `json:"fee"`
}

type TencentSmsUserCallbackRequestDto struct {
	Extend     string `json:"extend"`
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
	Sign       string `json:"sign"`
	Text       string `json:"text"`
	Time       int64  `json:"time"`
}

type TencentSmsUserCallbackResponseDto struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
}

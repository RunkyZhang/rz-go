package external

type DahanSmsMessageRequestDto struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	MsgId    string `json:"msgId"`
	Phones   string `json:"phones"`
	Content  string `json:"content"`
	Sign     string `json:"sign"`
	SubCode  string `json:"subcode"`
	SendTime int64  `json:"sendtime,string"`
}

//0	    提交成功
//1	    账号无效
//2	    密码错误
//3	    msgid太长，不得超过64位
//4	    错误号码/限制运营商号码
//5	    手机号码个数超过最大限制
//6	    短信内容超过最大限制
//7	    扩展子号码无效
//8	    定时时间格式错误
//14	手机号码为空
//19	用户被禁发或禁用
//20	ip鉴权失败
//21	短信内容为空
//24	无可用号码
//25	批量提交短信数超过最大限制
//98	系统正忙
//99	消息格式错误
type DahanSmsMessageResponseDto struct {
	MsgId      string `json:"msgId"`
	Result     string `json:"result"`
	Desc       string `json:"desc"`
	FailPhones string `json:"failPhones"`
}

type DahanSmsUserCallbackRequestDto struct {
	Result   string                                  `json:"result"`
	Desc     string                                  `json:"desc"`
	Delivers []*DahanSmsUserCallbackDeliverRequestDto `json:"delivers"`
}

type DahanSmsUserCallbackDeliverRequestDto struct {
	Phone       string `json:"phone"`
	Content     string `json:"content"`
	SubCode     string `json:"subcode"`
	DeliverTime string `json:"delivertime"`
}

type DahanSmsUserCallbackResponseDto struct {
	Status string `json:"status"`
}

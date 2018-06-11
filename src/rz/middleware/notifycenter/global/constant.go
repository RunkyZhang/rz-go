package global

const (
	ArgumentNameConfig              = "config"
	RedisKeyMessage                 = "middleware_notifyCenter_"
	RedisKeyMessageIds             = RedisKeyMessage + "Ids_"
	//RedisKeyMessageValues           = RedisKeyMessage + "values_"
	RedisKeySmsUserCallbcaks        = RedisKeyMessage + "sms_user_callbacks"
	RedisKeySmsTemplates            = RedisKeyMessage + "sms_templates"
	RedisKeySmsUserCallbackMessages = RedisKeyMessage + "sms_user_callback_messages"
)

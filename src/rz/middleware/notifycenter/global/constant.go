package global

const (
	ArgumentNameConfig              = "config"
	RedisKeyMessage                 = "middleware_notifyCenter_"
	RedisKeyMessageKeys             = RedisKeyMessage + "keys_"
	RedisKeyMessageValues           = RedisKeyMessage + "values_"
	RedisKeySmsUserCallbcaks        = RedisKeyMessage + "sms_user_callbacks"
	RedisKeySmsTemplates            = RedisKeyMessage + "sms_templates"
	RedisKeySmsUserCallbackMessages = RedisKeyMessage + "sms_user_callback_messages"
)

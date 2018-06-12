package models

type MailMessageDto struct {
	MessageBaseDto

	Subject string `json:"subject"`
}

type MailMessagePo struct {
	MessageBasePo

	Subject string `gorm:"column:subject"`
}

func MailMessageDtoToPo(mailMessageDto *MailMessageDto) (*MailMessagePo) {
	mailMessagePo := &MailMessagePo{}
	mailMessagePo.MessageBasePo = *MessageBaseDtoToPo(&mailMessageDto.MessageBaseDto)
	mailMessagePo.Subject = mailMessageDto.Subject

	return mailMessagePo
}

func MailMessagePoToDto(mailMessagePo *MailMessagePo) (*MailMessageDto) {
	mailMessageDto := &MailMessageDto{}
	mailMessageDto.MessageBaseDto = *MessageBasePoToDto(&mailMessagePo.MessageBasePo)
	mailMessageDto.Subject = mailMessagePo.Subject

	return mailMessageDto
}

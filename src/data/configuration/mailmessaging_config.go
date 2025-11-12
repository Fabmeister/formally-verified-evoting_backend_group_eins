package configuration

type MessagingConfig struct {
	DomainInMessage     string
	UseCliInsteadOfSmtp bool
	SendMails           bool // togglebar z.B. für Unittests
	SmtpServer          string
	SmtpPort            int
	SmtpSenderAddress   string
}

func (this *MessagingConfig) SetStandard() {
	this.DomainInMessage = "domain.example"
	this.UseCliInsteadOfSmtp = false
	this.SendMails = false
	this.SmtpServer = "localhost"
	this.SmtpPort = 25
	this.SmtpSenderAddress = "noreply@domain.example"
}

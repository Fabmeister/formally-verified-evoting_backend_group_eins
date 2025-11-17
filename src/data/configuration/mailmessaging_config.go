// Formally verified E-Voting using Dafny
// Copyright (C) 2025 Authors Gruppe EinS
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

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

package email

import (
	"e-voting-service/data/configuration"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"os/exec"
)

func SendMail(msg Mail_Message) {
	if !configuration.GlobalConfig.Messaging.SendMails {
		return
	}

	if configuration.GlobalConfig.Messaging.UseCliInsteadOfSmtp {
		sendMailCLI(msg)
	} else {
		sendMailSmtpModule(msg)
	}
}

func sendMailSmtpModule(msg Mail_Message) {
	// mit smtp modul
	// https://pkg.go.dev/net/smtp

	err := smtp.SendMail(
		fmt.Sprintf(
			"%s:%d",
			configuration.GlobalConfig.Messaging.SmtpServer,
			configuration.GlobalConfig.Messaging.SmtpPort),
		nil, // Future TODO: Auth
		configuration.GlobalConfig.Messaging.SmtpSenderAddress,
		msg.GetMailAddresses(),
		[]byte(fmt.Sprintf(
			"Subject: %s\r\n\r\n%s\r\n",
			msg.GetSubject(),
			msg.GetBody(),
		)),
	)

	if err != nil {
		log.Printf("Fehler beim Mailversand: %v", err)
	}
}

func sendMailCLI(msg Mail_Message) {
	// mit command line tool
	args := []string{"-s", msg.GetSubject()}
	args = append(args, msg.GetMailAddresses()...)
	cmd := exec.Command("mail", args...)

	reader, writer := io.Pipe()
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout

	go func() {
		defer writer.Close()

		writer.Write([]byte("\n"))
		writer.Write([]byte(msg.GetBody()))
		writer.Write([]byte("\n"))
	}()

	err := cmd.Run()

	if err != nil {
		log.Printf("Fehler beim Mailversand: %v", err)
	}
}

// Erweiterungsmöglichkeit: Sprachauswahl und/oder Möglichkeit zum customizen der Nachrichten -> laden aus Ressourcendatei
// Potenziel auch variablen für Kontext (z.B. Wahlname, Ablaufdatum, etc.) für custom Mails erlauben
type Mail_Message interface {
	GetSubject() string
	GetBody() string
	SetMailAddresses(to []string)
	GetMailAddresses() []string
}

type StandardMailAddressHolder struct {
	To []string
}

func (mail_holder *StandardMailAddressHolder) SetMailAddresses(to []string) {
	mail_holder.To = to
}

func (mail_holder StandardMailAddressHolder) GetMailAddresses() []string {
	return mail_holder.To
}

type VoteInviteMail struct {
	StandardMailAddressHolder
	Token string
}

func (VoteInviteMail) GetSubject() string {
	return "Einladung zur Wahl"
}

func (mail VoteInviteMail) GetBody() string {
	return fmt.Sprintf(
		"Sehr geehrter User,\n"+
			"Sie wurden zur Wahl eingeladen: \n"+
			"%s/?VoterToken=%s \n mit freundlichen Grüßen,\n"+
			"Ihr Wahlteam",
		configuration.GlobalConfig.Messaging.DomainInMessage,
		mail.Token,
	)
}

package mailconsumer

import (
	"fmt"
	"html/template"
	"net/smtp"
)

import (
	"github.com/shal/pigeon/pkg/mailer"
	"github.com/shal/pigeon/pkg/utils"
)

type MailAccountTpl struct {
	mailer.MailTpl

	Record AccountRecord `json:"record"`
}

func SendEmail(record AccountRecord, cli *smtp.Client) error {
	sender := utils.GetEnv("SENDER_EMAIL", "example@domain.com")

	fmt.Println("BR 1")

	if err := cli.Mail(sender); err != nil {
		return err
	}

	fmt.Println("BR 2")

	if err := cli.Rcpt(record.Email); err != nil {
		return err
	}

	fmt.Println("BR 3")

	wc, err := cli.Data()
	if err != nil {
		return err
	}

	fmt.Println("BR 4")

	mailTpl := MailAccountTpl{
		mailer.MailTpl{
			To:   sender,
			From: record.Email,
		},
		record,
	}

	fmt.Println("BR 5")

	tpl, err := template.ParseFiles("templates/sign_up.tpl")
	if err != nil {
		return err
	}

	fmt.Println("BR 6")

	fmt.Fprintf(wc, "Message-Id: <%s>\r\n", mailTpl.MessageID())

	tpl.Execute(wc, mailTpl)
	if err != nil {
		return err
	}

	fmt.Println("BR 7")

	err = wc.Close()
	if err != nil {
		return err
	}

	fmt.Println("BR 8")

	// Send the QUIT command and close the connection.
	err = cli.Quit()
	if err != nil {
		return err
	}

	fmt.Println("BR 9")

	return nil
}

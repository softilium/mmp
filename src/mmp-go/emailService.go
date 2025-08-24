package main

import (
	"fmt"

	"github.com/resend/resend-go/v2"
	"github.com/wneessen/go-mail"
)

func SendEmail(To string, Subject string, Body string) error {
	switch Cfg.EmailMode {
	case EMAIL_MODE_DIRECT:
		return sendEmailDirect(To, Subject, Body)
	case EMAIL_MODE_RESEND:
		return sendEmailResend(To, Subject, Body)
	}
	return nil
}

func sendEmailResend(To string, Subject string, Body string) error {

	apiKey := Cfg.Email_Resend_ApiKey

	client := resend.NewClient(apiKey)

	keys, _ := client.ApiKeys.List()

	fmt.Printf("%+v'n", keys)

	fmt.Printf("key=%+v\n", client.ApiKey)

	params := &resend.SendEmailRequest{
		From:    Cfg.EmailFrom,
		To:      []string{To},
		Subject: Subject,
		Html:    Body,
		Text:    Body,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		rerr := fmt.Errorf("failed to send email via Resend: %w", err)
		fmt.Printf("%v\n", rerr)
		return rerr
	} else {
		fmt.Printf("Email sent via Resend with ID: %s\n", sent.Id)
		return nil
	}
}

func sendEmailDirect(To string, Subject string, Body string) error {
	message := mail.NewMsg()
	if err := message.From(Cfg.EmailFrom); err != nil {
		rerr := fmt.Errorf("failed to set From address: %w", err)
		fmt.Printf("%v\n", rerr)
		return rerr
	}
	if err := message.To(To); err != nil {
		rerr := fmt.Errorf("failed to set To address: %s", err)
		fmt.Printf("%v\n", rerr)
		return rerr
	}
	message.Subject(Subject)
	message.SetBodyString(mail.TypeTextPlain, Body)
	client, err := mail.NewClient(Cfg.EmailSMTPHost,
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithPort(Cfg.EmailSMTPPort),
		mail.WithUsername(Cfg.EmailSMTPUserName),
		mail.WithPassword(Cfg.EmailSMTPPassword))
	if err != nil {
		rerr := fmt.Errorf("failed to create mail client: %w", err)
		fmt.Printf("%v\n", rerr)
		return rerr
	}
	if err := client.DialAndSend(message); err != nil {
		rerr := fmt.Errorf("failed to send mail: %w", err)
		fmt.Printf("%v\n", rerr)
		return rerr
	}
	return nil
}

package main

import (
	"fmt"
	"golang.org/x/net/context"
	"net/smtp"
	pb "server/rpcdefsql"
	"strings"
)

const (
	SMTP_SERVER = "smtp.gmail.com"
	SMTP_PORT   = "587"
)

type Sender struct {
	User     string
	Password string
}

func NewSender(Username, Password string) Sender {

	return Sender{Username, Password}
}

var subject string = "Welcome to Soulfit!"
var content string = "Thanks for joining Soulfit! Please check www.soulfit.com for details."

func (sender Sender) SendMail(Dest []string, Subject, bodyMessage string) error {

	msg := "From: " + sender.User + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(SMTP_SERVER+":"+SMTP_PORT,
		smtp.PlainAuth("", sender.User, sender.Password, SMTP_SERVER),
		sender.User, Dest, []byte(msg))

	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return err
	}

	fmt.Println("Mail sent successfully!")
	return nil
}

func SendWelcomeEmail(toUser string) error {
	sender := NewSender("libera.labs@gmail.com", "Lib3r@99")
	return sender.SendMail([]string{toUser, "libera.labs@gmail.com"},
		subject,
		content,
	)
}

func (s *server) SendWelcomeEmailToUser(ctx context.Context, in *pb.SendWelcomeEmailReq) (*pb.SendWelcomeEmailReply,
	error) {

	var resp pb.SendWelcomeEmailReply

	err := SendWelcomeEmail(in.GetToEmail())

	return &resp, err
}

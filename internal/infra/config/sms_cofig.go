package config

import (
	"os"
)

type SendOtpSms struct {
	AccountSid string
	AuthToken  string
}

func NewSendOtpSms() *SendOtpSms {
	return &SendOtpSms{
		AccountSid: os.Getenv("ACCOUNT_SID"),
		AuthToken:  os.Getenv("AUTH_TOKEN"),
	}
}

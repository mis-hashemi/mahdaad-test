package sms_client

import (
	"fmt"
	"github.com/mis-hashemi/mahdaad-test/adapter/sms_client/test_client"
)

var storageClientNotExist = func(name string) error { return fmt.Errorf("invalid sms client : %s", name) }

const (
	Test = "test"
)

func NewSMSClient(name string, config any) (SMSClient, error) {
	switch name {
	case Test:
		return test_client.New()
	default:
		return nil, storageClientNotExist(name)
	}
}

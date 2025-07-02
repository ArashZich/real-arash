package sms

import (
	"fmt"

	"gitag.ir/armogroup/armo/services/reality/notification/notif"
	"github.com/kavenegar/kavenegar-go"
)

type kvnegar struct{ api *kavenegar.Kavenegar }

func NewKavenegar(key string) notif.Driver {
	api := kavenegar.New(key)

	return kvnegar{api}
}

func (k kvnegar) Send(to string, subject, message string) error {
	panic("implement me. later decouple from sendWithTemplate")
}

func (k kvnegar) SendWithTemplate(to string, subject, message string, tmpl string) error {
	receptor := to
	template := tmpl
	token := message
	params := &kavenegar.VerifyLookupParam{}
	if res, err := k.api.Verify.Lookup(receptor, template, token, params); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			fmt.Println(err.Error())
		case *kavenegar.HTTPError:
			fmt.Println(err.Error())
		default:
			fmt.Println(err.Error())
		}
		return err
	} else {
		fmt.Println("MessageID 	= ", res.MessageID)
		fmt.Println("Status    	= ", res.Status)
		return nil
	}
}

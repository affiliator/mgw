package mailgun_processor

import (
	"context"
	"fmt"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/mailgun/mailgun-go"
	"log"
	"mailgun-mgw/config"
	"time"
)

var apiCredentials ApiCredentials

type ApiCredentials struct {
	Domain string `json:"mailgun_domain"`
	ApiKey string `json:"mailgun_api_key"`
}

// The MyFoo decorator [enter what it does]
var MailgunProcessor = func() backends.Decorator {

	initializer := backends.InitializeWith(func(c backends.BackendConfig) error {
		r := config.GetInstance().Paths.Credentials.ReadTo(&apiCredentials)
		if r != nil {
			return r
		}

		return nil
	})

	backends.Svc.AddInitializer(initializer)

	return func(p backends.Processor) backends.Processor {
		return backends.ProcessWith(
			func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {
				if task == backends.TaskValidateRcpt {

					// if you want your processor to validate recipents,
					// validate recipient by checking
					// the last item added to `e.RcptTo` slice
					// if error, then return something like this:
					/* return backends.NewResult(
					   response.Canned.FailNoSenderDataCmd),
					   backends.NoSuchUser
					*/
					// if no error:
					return p.Process(e, task)
				} else if task == backends.TaskSaveMail {

					fmt.Println("Sending using mailgun..")
					mg := mailgun.NewMailgun(apiCredentials.Domain, apiCredentials.ApiKey)
					mg.SetAPIBase("")

					// The message object allows you to add attachments and Bcc recipients
					message := mg.NewMessage(e.MailFrom.String(), e.Subject, "Das ist ein Test!", e.RcptTo[0].String())

					fmt.Println("Sending from: " + e.MailFrom.String() + " / to: " + e.RcptTo[0].String())

					var ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()

					// Send the message	with a 10 second timeout
					resp, id, err := mg.Send(ctx, message)

					if err != nil {
						log.Fatal(err)
					}

					fmt.Printf("ID: %s Resp: %s\n", id, resp)

					// if you want your processor to do some processing after
					// receiving the email, continue here.
					// if want to stop processing, return
					// errors.New("Something went wrong")
					// return backends.NewBackendResult(fmt.Sprintf("554 Error: %s", err)), err
					// call the next processor in the chain
					return p.Process(e, task)
				}
				return p.Process(e, task)
			},
		)
	}
}

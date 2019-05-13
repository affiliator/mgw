package mailgun_processor

import (
	"context"
	"fmt"
	"github.com/affiliator/mgw/config"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/mailgun/mailgun-go"
	"log"
	"time"
)

var (
	conf      = config.Ptr()
	overrides Overrides
)

type Overrides struct {
	Credentials config.Credentials `json:"mailgun,omitempty"`
}

func createProviderClient() *mailgun.MailgunImpl {
	c := config.Ptr().Providers.Mailgun
	mg := mailgun.NewMailgun(c.Credentials.Domain, c.Credentials.ApiKey)

	if c.ApiBase != "" {
		mg.SetAPIBase(c.ApiBase)
	}

	return mg
}

var MailgunProcessor = func() backends.Decorator {
	initializer := backends.InitializeWith(func(c backends.BackendConfig) error {
		if result, _ := conf.Paths.Credentials.Exists(); result == false {
			return nil
		}

		e := conf.Paths.Credentials.ReadTo(&overrides)
		if e != nil {
			return e
		}

		cred := &conf.Providers.Mailgun.Credentials
		if overrides.Credentials.ApiKey != "" {
			cred.ApiKey = overrides.Credentials.ApiKey
		}

		if overrides.Credentials.Domain != "" {
			cred.Domain = overrides.Credentials.Domain
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
					mg := createProviderClient()

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

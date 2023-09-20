package mail

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/alayton/papers"
)

func NewSendgridMailer(apiKey string) *SendgridMailer {
	return &SendgridMailer{
		ApiKey:       apiKey,
		Endpoint:     "/v3/mail/send",
		Host:         "https://api.sendgrid.com",
		MessageTypes: map[papers.MessageType]string{},
	}
}

type SendgridMailer struct {
	ApiKey       string
	Endpoint     string
	Host         string
	MessageTypes map[papers.MessageType]string
}

func (s SendgridMailer) AddMessageType(msgType papers.MessageType, templateID string) {
	s.MessageTypes[msgType] = templateID
}

// Send an e-mail
func (s SendgridMailer) Send(ctx context.Context, p *papers.Papers, msg papers.Message) error {
	templateID, ok := s.MessageTypes[msg.Type]
	if !ok {
		return fmt.Errorf("%w: %s", papers.ErrNoMessageTemplate, msg.Type)
	}

	m := mail.NewV3Mail()
	m.SetFrom(mail.NewEmail(p.Config.Mailer.From.Name, p.Config.Mailer.From.Address))

	p13n := mail.NewPersonalization()
	tos := []*mail.Email{}
	for _, to := range msg.To {
		tos = append(tos, mail.NewEmail(to.Name, to.Address))
	}
	p13n.AddTos(tos...)

	ccs := []*mail.Email{}
	for _, cc := range msg.Cc {
		ccs = append(ccs, mail.NewEmail(cc.Name, cc.Address))
	}
	p13n.AddCCs(ccs...)

	bccs := []*mail.Email{}
	for _, bcc := range msg.Bcc {
		bccs = append(bccs, mail.NewEmail(bcc.Name, bcc.Address))
	}
	p13n.AddBCCs(bccs...)

	for k, v := range msg.Data {
		p13n.SetDynamicTemplateData(k, v)
	}

	m.SetTemplateID(templateID)
	m.AddPersonalizations(p13n)

	request := sendgrid.GetRequest(s.ApiKey, s.Endpoint, s.Host)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		return fmt.Errorf("%w: %v", papers.ErrMessageFailed, err)
	} else {
		if response.StatusCode >= 400 {
			p.Logger.Print("Unexpected status code from SendGrid API:", response.StatusCode, response.Body)
		}
	}

	return err
}

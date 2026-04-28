package actions

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/mudler/LocalAGI/core/types"
	"github.com/mudler/LocalAGI/pkg/config"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func NewSendMail(config map[string]string) *SendMailAction {
	s := &SendMailAction{
		username: config["username"],
		password: config["password"],
		email:    config["email"],
		smtpHost: config["smtpHost"],
		smtpPort: config["smtpPort"],
	}
	s.fromEmail = config["fromEmail"]
	if s.fromEmail == "" {
		s.fromEmail = config["email"]
	}
	return s
}

type SendMailAction struct {
	username  string
	password  string
	email     string
	fromEmail string
	smtpHost  string
	smtpPort  string
}

func (a *SendMailAction) Run(ctx context.Context, sharedState *types.AgentSharedState, params types.ActionParams) (types.ActionResult, error) {
	result := struct {
		Message string `json:"message"`
		To      string `json:"to"`
		Subject string `json:"subject"`
	}{}
	err := params.Unmarshal(&result)
	if err != nil {
		fmt.Printf("error: %v", err)

		return types.ActionResult{}, err
	}

	// Build RFC 2822 message with proper headers
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: ION <%s>\r\n", a.fromEmail))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", result.To))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", result.Subject))
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(result.Message)

	// Authentication.
	auth := smtp.PlainAuth("", a.email, a.password, a.smtpHost)

	// Sending email.
	err = smtp.SendMail(
		fmt.Sprintf("%s:%s", a.smtpHost, a.smtpPort),
		auth, a.email, []string{
			result.To,
		}, buf.Bytes())
	if err != nil {
		return types.ActionResult{}, err
	}
	return types.ActionResult{Result: fmt.Sprintf("Email sent to %s", result.To)}, nil
}

func (a *SendMailAction) Definition() types.ActionDefinition {
	return types.ActionDefinition{
		Name:        "send_email",
		Description: "Send an email.",
		Properties: map[string]jsonschema.Definition{
			"to": {
				Type:        jsonschema.String,
				Description: "The email address to send the email to.",
			},
			"subject": {
				Type:        jsonschema.String,
				Description: "The subject of the email.",
			},
			"message": {
				Type:        jsonschema.String,
				Description: "The message to send.",
			},
		},
		Required: []string{"to", "subject", "message"},
	}
}

func (a *SendMailAction) Plannable() bool {
	return true
}

// SendMailConfigMeta returns the metadata for SendMail action configuration fields
func SendMailConfigMeta() []config.Field {
	return []config.Field{
		{
			Name:     "smtpHost",
			Label:    "SMTP Host",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "SMTP server host (e.g., smtp.gmail.com)",
		},
		{
			Name:         "smtpPort",
			Label:        "SMTP Port",
			Type:         config.FieldTypeText,
			Required:     true,
			DefaultValue: "587",
			HelpText:     "SMTP server port (e.g., 587)",
		},
		{
			Name:     "username",
			Label:    "SMTP Username",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "SMTP username/email address",
		},
		{
			Name:     "password",
			Label:    "SMTP Password",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "SMTP password or app password",
		},
		{
			Name:     "email",
			Label:    "Auth Email (SMTP)",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "SMTP authentication email address",
		},
		{
			Name:     "fromEmail",
			Label:    "From Email",
			Type:     config.FieldTypeText,
			Required: false,
			HelpText: "Display From address (defaults to Auth Email if empty)",
		},
	}
}

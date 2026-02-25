package fnMailer

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"time"
	"tower/pkg/fnEnv"

	"github.com/mailgun/mailgun-go/v5"
)

type Mailer struct {
	mg     *mailgun.Client
	domain string
	sender string
}

var defaultMailer *Mailer

func Init() {
	domain := fnEnv.GetString("MAILGUN_DOMAIN", "")
	apiKey := fnEnv.GetString("MAILGUN_API_KEY", "")
	sender := fnEnv.GetString("MAIL_SENDER_ADDRESS", "")

	if domain == "" || apiKey == "" {
		log.Println("[fnMailer] Mailgun 환경변수가 설정되지 않아 메일러를 초기화할 수 없습니다.")
		return
	}

	mg := mailgun.NewMailgun(apiKey)
	defaultMailer = &Mailer{
		mg:     mg,
		domain: domain,
		sender: sender,
	}
}

func CompileTemplate(tmplStr string, data map[string]interface{}) (string, error) {
	t, err := template.New("mail_tmpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Send(to, subject, htmlContent string) error {
	if defaultMailer == nil {
		log.Println("[fnMailer] 메일러가 초기화되지 않았습니다. 발송 스킵:", to)
		return nil
	}

	message := mailgun.NewMessage(defaultMailer.domain, defaultMailer.sender, subject, "", to)
	message.SetHTML(htmlContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := defaultMailer.mg.Send(ctx, message)
	if err != nil {
		log.Printf("[fnMailer] 이메일 발송 실패 (%s): %v\n", to, err)
		return err
	}

	log.Printf("[fnMailer] 이메일 발송 성공 (ID: %s, To: %s)\n", resp.ID, to)
	return nil
}

package fnMailer

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"sync"
	"time"
	"tower/pkg/fnEnv"

	"github.com/mailgun/mailgun-go/v5"
)

type Mailer struct {
	mg     *mailgun.Client
	domain string
	sender string
}

var (
	instance *Mailer
	once     sync.Once
)

func getMailer() *Mailer {
	once.Do(func() {
		domain := fnEnv.App.MailgunDomain
		apiKey := fnEnv.App.MailgunAPIKey
		sender := fnEnv.App.MailSender

		if domain == "" || apiKey == "" {
			log.Println("[fnMailer] Mailgun 환경변수가 설정되지 않아 메일러를 초기화할 수 없습니다.")
			return
		}

		mg := mailgun.NewMailgun(apiKey)
		instance = &Mailer{
			mg:     mg,
			domain: domain,
			sender: sender,
		}
		log.Println("[fnMailer] Mailgun 객체 지연 초기화 완료 📧")
	})
	return instance
}

func CompileTemplate(tmplStr string, data any) (string, error) {
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
	mailer := getMailer()
	if mailer == nil {
		log.Println("[fnMailer] 메일러가 초기화되지 않았습니다. 발송 스킵:", to)
		return nil
	}

	message := mailgun.NewMessage(mailer.domain, mailer.sender, subject, "", to)
	message.SetHTML(htmlContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := mailer.mg.Send(ctx, message)
	if err != nil {
		log.Printf("[fnMailer] 이메일 발송 실패 (%s): %v\n", to, err)
		return err
	}

	log.Printf("[fnMailer] 이메일 발송 성공 (ID: %s, To: %s)\n", resp.ID, to)
	return nil
}

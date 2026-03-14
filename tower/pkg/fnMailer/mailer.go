package fnMailer

import (
	"bytes"
	"context"
	"log"
	"text/template"
	"time"
	"tower/model/maindb"

	"github.com/mailgun/mailgun-go/v5"
	"gorm.io/gorm"
)

type Config struct {
	Domain string
	APIKey string
	Sender string
}

type Mailer struct {
	mg     *mailgun.Client
	domain string
	sender string
}

var (
	instance *Mailer
	logDB    *gorm.DB
)

func Init(cfg Config, db *gorm.DB) {
	if cfg.Domain == "" || cfg.APIKey == "" {
		log.Println("[fnMailer] ⚠️ Mailgun 환경변수가 누락되어 메일러를 초기화할 수 없습니다.")
		return
	}

	mg := mailgun.NewMailgun(cfg.APIKey)
	instance = &Mailer{
		mg:     mg,
		domain: cfg.Domain,
		sender: cfg.Sender,
	}

	logDB = db
	log.Println("[fnMailer] Mailgun 클라이언트 초기화 및 DB 주입 완료 📧")
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

func Send(templateCode, to, subject, htmlContent string) error {
	message := mailgun.NewMessage(instance.domain, instance.sender, subject, "", to)
	message.SetHTML(htmlContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := instance.mg.Send(ctx, message)

	status := maindb.EmailLogStatusSuccess
	var providerID, errMsg string
	if err != nil {
		status = maindb.EmailLogStatusFailed
		errMsg = err.Error()
		log.Printf("발송 실패 (%s): %v\n", to, err)
	} else {
		providerID = resp.ID
		log.Printf("발송 성공 (ID: %s, To: %s)\n", resp.ID, to)
	}

	if logDB != nil {
		go func(logData maindb.EmailLog) {
			if insertErr := logDB.Create(&logData).Error; insertErr != nil {
				log.Printf("[fnMailer] 🚨 이메일 로그 DB 저장 실패: %v\n", insertErr)
			}
		}(maindb.EmailLog{
			TemplateCode: templateCode,
			Sender:       instance.sender,
			Recipient:    to,
			Subject:      subject,
			HTMLBody:     htmlContent,
			Status:       status,
			ProviderID:   providerID,
			ErrorMessage: errMsg,
		})
	}
	return nil
}

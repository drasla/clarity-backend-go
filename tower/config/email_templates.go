package config

import "tower/pkg/fnMailer"

type TemplateCode string

const (
	TplInquiryReply TemplateCode = "INQUIRY_REPLY"
	TplWelcomeUser  TemplateCode = "WELCOME_USER"
	TplBilling      TemplateCode = "BILLING"
)

var EssentialEmailTemplates = []fnMailer.SystemEmailTemplate{
	{
		Code:        string(TplInquiryReply),
		Subject:     "[골든넷] 문의하신 내용에 대한 답변이 등록되었습니다.",
		Variables:   `["Inquiry.Title", "Inquiry.Category", "Answer"]`,
		Description: "문의 답변 등록 시 발송되는 이메일 템플릿입니다.",
	},
	{
		Code:        string(TplWelcomeUser),
		Subject:     "[골든넷] 환영합니다! 회원가입이 완료되었습니다.",
		Variables:   `["User.Name", "User.Email"]`,
		Description: "회원가입 직후 자동으로 발송되는 환영 메일입니다.",
	},
	{
		Code:        string(TplBilling),
		Subject:     "[골든넷] 청구서가 발행되었습니다.",
		Variables:   `["User.Name"]`,
		Description: "회원의 청구서 발행 시 자동으로 발송되는 이메일 템플릿입니다.",
	},
}

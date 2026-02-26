package config

import (
	"log"
	"tower/model/maindb"

	"gorm.io/gorm"
)

func seedData(db *gorm.DB) {
	essentialTemplates := []maindb.EmailTemplate{
		{
			TemplateCode: "INQUIRY_REPLY",
			Subject:      "[안내] 문의하신 내용에 대한 답변이 등록되었습니다.",
			HTMLBody: `
<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h2>문의하신 내용에 대한 답변입니다.</h2>
	<hr />
	<p><strong>문의 제목:</strong> {{.Inquiry.Title}}</p>
	<div style="background-color: #f9f9f9; padding: 15px; margin-top: 20px; border-radius: 5px;">
		<p>{{.Answer}}</p>
	</div>
</div>`,
			Variables:   `["Inquiry.Title", "Inquiry.Category", "Answer"]`,
			Description: "1:1 문의 답변 시 발송되는 필수 기본 템플릿입니다.",
		},
	}

	for _, t := range essentialTemplates {
		var existing maindb.EmailTemplate
		if err := db.Where(maindb.EmailTemplate{TemplateCode: t.TemplateCode}).FirstOrCreate(&existing, t).Error; err != nil {
			log.Printf("[Registry] 🚨 필수 이메일 템플릿(%s) 시딩 실패: %v\n", t.TemplateCode, err)
		}
	}
	log.Println("[Registry] 🌱 필수 초기 데이터(Seed) 검증 및 주입 완료")
}

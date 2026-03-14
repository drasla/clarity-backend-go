package config

import (
	"log"
	"tower/model/maindb"

	"gorm.io/gorm"
)

func seedData(db *gorm.DB) {
	for _, sysTmpl := range EssentialEmailTemplates {
		seedModel := maindb.EmailTemplate{
			TemplateCode: string(sysTmpl.Code),
			Subject:      sysTmpl.Subject,
			HTML:         "",
			Design:       "{}",
			Variables:    sysTmpl.Variables,
			Description:  sysTmpl.Description,
		}

		var existing maindb.EmailTemplate
		if err := db.Where(maindb.EmailTemplate{TemplateCode: string(sysTmpl.Code)}).
			FirstOrCreate(&existing, seedModel).Error; err != nil {
			log.Printf("[Registry] 🚨 필수 이메일 템플릿(%s) 시딩 실패: %v\n", sysTmpl.Code, err)
		}
	}

	log.Println("[Registry] 🌱 필수 초기 데이터(Seed) 검증 및 주입 완료")
}

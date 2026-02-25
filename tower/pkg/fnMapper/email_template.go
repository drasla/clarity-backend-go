package fnMapper

import (
	"tower/graph/model"
	"tower/model/maindb"
)

func EmailTemplateToGraphQL(emailTemplate *maindb.EmailTemplate) *model.EmailTemplate {
	if emailTemplate == nil {
		return nil
	}

	return &model.EmailTemplate{
		ID:           int(emailTemplate.ID),
		CreatedAt:    emailTemplate.CreatedAt,
		UpdatedAt:    emailTemplate.UpdatedAt,
		TemplateCode: emailTemplate.TemplateCode,
		Subject:      emailTemplate.Subject,
		HTMLBody:     emailTemplate.HTMLBody,
		Variables:    &emailTemplate.Variables,
		Description:  &emailTemplate.Description,
	}
}

func EmailTemplatesToGraphQL(emailTemplates []maindb.EmailTemplate) []*model.EmailTemplate {
	var list []*model.EmailTemplate
	for i := range emailTemplates {
		list = append(list, EmailTemplateToGraphQL(&emailTemplates[i]))
	}
	return list
}

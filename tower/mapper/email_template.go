package mapper

import (
	"tower/graph/model"
	"tower/model/maindb"
)

func EmailTemplateToGraphQL(emailTemplate *maindb.EmailTemplate) *model.EmailTemplate {
	if emailTemplate == nil {
		return nil
	}

	var variables, description *string
	if emailTemplate.Variables != "" {
		variables = new(emailTemplate.Variables)
	}
	if emailTemplate.Description != "" {
		description = new(emailTemplate.Description)
	}

	return &model.EmailTemplate{
		ID:           int(emailTemplate.ID),
		CreatedAt:    emailTemplate.CreatedAt,
		UpdatedAt:    emailTemplate.UpdatedAt,
		TemplateCode: emailTemplate.TemplateCode,
		Subject:      emailTemplate.Subject,
		HTML:         emailTemplate.HTML,
		Design:       emailTemplate.Design,
		Variables:    variables,
		Description:  description,
	}
}

func EmailTemplatesToGraphQL(emailTemplates []maindb.EmailTemplate) []*model.EmailTemplate {
	var list []*model.EmailTemplate
	for i := range emailTemplates {
		list = append(list, EmailTemplateToGraphQL(&emailTemplates[i]))
	}
	return list
}

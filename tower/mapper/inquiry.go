package mapper

import (
	"tower/graph/model"
	"tower/model/maindb"
)

func InquiryToGraphQL(inquiry *maindb.Inquiry) *model.Inquiry {
	if inquiry == nil {
		return nil
	}

	var files []*model.File
	for _, f := range inquiry.Attachments {
		files = append(files, &model.File{
			ID:           int(f.ID),
			OriginalName: f.OriginalName,
			StoredName:   f.StoredName,
			URL:          f.URL,
			Size:         f.Size,
			Extension:    f.Extension,
			CreatedAt:    f.CreatedAt,
		})
	}

	var userId *int
	if inquiry.UserID != nil {
		userId = new(int(*inquiry.UserID))
	}

	return &model.Inquiry{
		ID:          int(inquiry.ID),
		CreatedAt:   inquiry.CreatedAt,
		UpdatedAt:   inquiry.UpdatedAt,
		UserID:      userId,
		Category:    model.InquiryCategory(inquiry.Category),
		Domain:      inquiry.Domain,
		Title:       inquiry.Title,
		Content:     inquiry.Content,
		Email:       inquiry.Email,
		PhoneNumber: inquiry.PhoneNumber,
		Status:      model.InquiryStatus(inquiry.Status),
		Answer:      inquiry.Answer,
		AnsweredAt:  inquiry.AnsweredAt,
		Attachments: files,
	}
}

func InquiriesToGraphQL(inquiries []maindb.Inquiry) []*model.Inquiry {
	var list []*model.Inquiry
	for i := range inquiries {
		list = append(list, InquiryToGraphQL(&inquiries[i]))
	}
	return list
}

func InquiryToPublicGraphQL(inquiry *maindb.Inquiry) *model.Inquiry {
	if inquiry == nil {
		return nil
	}

	phone := inquiry.PhoneNumber
	maskedPhone := "***-****-****"
	if len(phone) >= 4 {
		maskedPhone = "***-****-" + phone[len(phone)-4:]
	}

	return &model.Inquiry{
		ID:          int(inquiry.ID),
		CreatedAt:   inquiry.CreatedAt,
		UpdatedAt:   inquiry.UpdatedAt,
		UserID:      nil,
		Category:    model.InquiryCategory(inquiry.Category),
		Domain:      inquiry.Domain,
		Status:      model.InquiryStatus(inquiry.Status),
		Title:       inquiry.Title,
		Content:     "🔒 비회원으로 작성한 1:1문의 글입니다.",
		Email:       "***@***",
		PhoneNumber: maskedPhone,
		Answer:      nil,
		Attachments: []*model.File{},
	}
}

func InquiriesToPublicGraphQL(inquiries []maindb.Inquiry) []*model.Inquiry {
	var list []*model.Inquiry
	for i := range inquiries {
		list = append(list, InquiryToPublicGraphQL(&inquiries[i]))
	}
	return list
}

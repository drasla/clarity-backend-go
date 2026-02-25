package graph

import service "tower/services"

type Resolver struct {
	AuthService          service.AuthService
	VerificationService  service.VerificationService
	UserService          service.UserService
	InquiryService       service.InquiryService
	EmailTemplateService service.EmailTemplateService
	FileService          service.FileService
}

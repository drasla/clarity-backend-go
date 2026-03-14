package config

import (
	"log"
	"net/http"
	"tower/graph"
	"tower/pkg/fnEcho"
	"tower/pkg/fnError"
	"tower/pkg/fnMailer"
	"tower/repository"
	service "tower/services"

	"github.com/99designs/gqlgen/graphql"
)

type ServiceContainer struct {
	AuthService          service.AuthService
	VerificationService  service.VerificationService
	UserService          service.UserService
	InquiryService       service.InquiryService
	EmailTemplateService service.EmailTemplateService
	FileService          service.FileService
}

func newContainer(db *ProjectDB) *ServiceContainer {
	fnMailer.Init(fnMailer.Config{
		Domain: App.Mailgun.Domain,
		APIKey: App.Mailgun.APIKey,
		Sender: App.Mailgun.Sender,
	}, db.MainDB)
	var systemEmailTemplateCodes []string
	for _, sysTmpl := range EssentialEmailTemplates {
		systemEmailTemplateCodes = append(systemEmailTemplateCodes, sysTmpl.Code)
	}

	userRepo := repository.NewUserRepository(db.MainDB)
	sessionRepo := repository.NewSessionRepository(db.MainDB)
	verificationRepo := repository.NewVerificationRepository(db.MainDB)
	inquiryRepo := repository.NewInquiryRepository(db.MainDB)
	emailTemplateRepo := repository.NewEmailTemplateRepository(db.MainDB)

	verificationService := service.NewVerificationService(verificationRepo)

	authOpts := service.AuthOptions{
		JwtSecret:       App.Echo.JwtSecret,
		WelcomeTmplCode: string(TplWelcomeUser), // config/email_templates.go 에 정의된 상수 사용
	}
	authService := service.NewAuthService(userRepo, sessionRepo, verificationService, emailTemplateRepo, authOpts)
	userService := service.NewUserService(userRepo)

	inquiryService := service.NewInquiryService(inquiryRepo, emailTemplateRepo, service.TelegramOptions{
		TelegramBotToken: App.Telegram.BotToken,
		TelegramChatID:   App.Telegram.ChatID,
		ReplyTmplCode:    string(TplInquiryReply),
	})
	emailTemplateService := service.NewEmailTemplateService(emailTemplateRepo, systemEmailTemplateCodes)

	fileService, err := service.NewS3Service(service.S3Options{
		Domain:     App.S3.Domain,
		Endpoint:   App.S3.Endpoint,
		AccessKey:  App.S3.AccessKey,
		SecretKey:  App.S3.SecretKey,
		BucketName: App.S3.BucketName,
	})
	if err != nil {
		log.Fatalf("❌ Failed to initialize FileService: %v", err)
	}

	return &ServiceContainer{
		AuthService:          authService,
		VerificationService:  verificationService,
		UserService:          userService,
		InquiryService:       inquiryService,
		EmailTemplateService: emailTemplateService,
		FileService:          fileService,
	}
}

func NewExecutableSchema(db *ProjectDB) graphql.ExecutableSchema {
	services := newContainer(db)

	config := graph.Config{
		Resolvers: &graph.Resolver{
			AuthService:          services.AuthService,
			VerificationService:  services.VerificationService,
			UserService:          services.UserService,
			InquiryService:       services.InquiryService,
			EmailTemplateService: services.EmailTemplateService,
			FileService:          services.FileService,
		},
	}

	config.Directives.Auth = graph.AuthDirective
	config.Directives.Admin = graph.AdminDirective

	return graph.NewExecutableSchema(config)
}

func StartWebServer(errHandler *fnError.ErrorHandler, execSchema graphql.ExecutableSchema) *http.Server {
	echoCfg := fnEcho.Config{
		Port:      App.Echo.Port,
		JwtSecret: App.Echo.JwtSecret,
	}

	return fnEcho.StartEchoServer(echoCfg, errHandler, execSchema)
}

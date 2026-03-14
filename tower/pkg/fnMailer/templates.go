package fnMailer

type SystemEmailTemplate struct {
	Code        string
	Subject     string
	Variables   string
	Description string
}

const AllowedCustomVariables = `["User.Name", "User.Email", "User.PhoneNumber", "User.Type"]`

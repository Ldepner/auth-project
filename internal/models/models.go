package models

type TemplateData struct {
	Error     string
	Success   string
	StringMap map[string]string
}

type UserRecord struct {
	Email    string `bson:"email,omitempty"`
	Password string `bson:"password,omitempty"`
}

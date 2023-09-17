package responses

type ServiceAccount struct {
	AccountName string `json:"account_name"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Endpoint string `json:"endpoint"`
	Email string `'json':"email"`
}
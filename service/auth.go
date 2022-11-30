package service

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}

type AuthService interface {
	EncodeToken(id string) string
	DecodeToken(token string) (*User, error)
}

type authServices struct {
}

func NewAuthService(secretKey string) AuthService {
	return &authServices{}
}

func (service *authServices) EncodeToken(uuidStr string) string {

	return ""
}

func (service *authServices) DecodeToken(encodedToken string) (*User, error) {

	return &User{ID: -1}, nil
}

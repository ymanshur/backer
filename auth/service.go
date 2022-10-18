package auth

import "github.com/dgrijalva/jwt-go"

type Service interface {
	GenerateToken(userID int) (string, error)
}

type service struct {
}

func NewService() *service {
	return &service{}
}

var SECRET_KEY = []byte("BWABACKERSTARTUP_53cr3t_k3y")

func (s *service) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{
		"user_id": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

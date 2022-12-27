package authService

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
	"tosinjs/cloud-backup/internal/entity/errorEntity"
	"tosinjs/cloud-backup/internal/entity/tokenEntity"
)

type authService struct {
	key string
}

type Claims struct {
	Username string `json:"username"`
	Status   string `json:"status"`
	jwt.StandardClaims
}

type AuthService interface {
	GenerateJWT(xpTime int, payload tokenEntity.JWTPayload) (string, *errorEntity.ServiceError)
	ValidateJWT(jwt string) (*tokenEntity.JWTPayload, bool)
}

func New(key string) AuthService {
	return authService{
		key: key,
	}
}

func (a authService) GenerateJWT(xpTime int, payload tokenEntity.JWTPayload) (string, *errorEntity.ServiceError) {
	tokenLife := time.Now().Add(time.Duration(xpTime) * time.Minute)

	claims := &Claims{
		Username: payload.Username,
		Status:   payload.Status,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenLife.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.key))

	if err != nil {
		return "", errorEntity.InternalServerError(err)
	}
	return tokenString, nil
}

func (a authService) ValidateJWT(jwtString string) (*tokenEntity.JWTPayload, bool) {
	token, err := jwt.ParseWithClaims(jwtString, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(a.key), nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &tokenEntity.JWTPayload{
			Username: claims.Username,
			Status:   claims.Status,
		}, token.Valid
	}

	return nil, token.Valid
}

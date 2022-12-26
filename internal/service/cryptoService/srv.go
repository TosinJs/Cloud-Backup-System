package cryptoService

import (
	"golang.org/x/crypto/bcrypt"
	"tosinjs/cloud-backup/internal/entity/errorEntity"
)

type cryptoService struct{}

type CryptoService interface {
	Hash(input string) (string, *errorEntity.ServiceError)
	Compare(input, hashed string) bool
}

func New() CryptoService {
	return cryptoService{}
}

func (c cryptoService) Hash(input string) (string, *errorEntity.ServiceError) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), 10)

	if err != nil {
		return "", errorEntity.InternalServerError(err)
	}
	return string(bytes), nil

}

func (c cryptoService) Compare(input, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))

	return err == nil
}

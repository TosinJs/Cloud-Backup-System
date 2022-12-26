package validationService

import (
	"github.com/go-playground/validator/v10"
	"tosinjs/cloud-backup/internal/entity/errorEntity"
)

type validationService struct{}

type ValidationService interface {
	Validate(value any) *errorEntity.ServiceError
}

func New() ValidationService {
	return validationService{}
}

func (v validationService) Validate(value any) *errorEntity.ServiceError {
	validate := validator.New()
	if err := validate.Struct(value); err != nil {
		return errorEntity.BadRequestError(err.Error(), err)
	}
	return nil
}

package utils

import (
	"regexp"

	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/entity"
)

func CEPValidation(cep string) error {
	re := regexp.MustCompile(`^\d{8}$`)
	if !re.MatchString(cep) {
		return entity.ErrZipcodeNotValid
	}

	return nil
}

package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/entity"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/usecase"
)

type TemperatureHandler struct {
	useCase *usecase.GetTemperatureUseCase
}

func NewTemperatureHandler(useCase *usecase.GetTemperatureUseCase) *TemperatureHandler {
	return &TemperatureHandler{
		useCase: useCase,
	}
}

func (h *TemperatureHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if err := h.CEPValidation(cep); err != nil {
		h.handlerError(w, err)
		return
	}

	weather, err := h.useCase.Execute(r.Context(), cep)
	if err != nil {
		h.handlerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weather); err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TemperatureHandler) CEPValidation(cep string) error {
	re := regexp.MustCompile(`^\d{5}-\d{3}$`)
	if !re.MatchString(cep) {
		return entity.ErrZipcodeNotValid
	}

	return nil
}

func (h *TemperatureHandler) handlerError(w http.ResponseWriter, err error) {
	fmt.Println("error:", err)

	switch {
	case errors.Is(err, entity.ErrZipcodeNotValid):
		Error(w, entity.ErrZipcodeNotValid.Error(), http.StatusUnprocessableEntity)
		return
	case errors.Is(err, entity.ErrZipcodeNotFound):
		Error(w, entity.ErrZipcodeNotFound.Error(), http.StatusNotFound)
		return
	case err != nil:
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

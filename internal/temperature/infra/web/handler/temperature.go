package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/entity"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/usecase"
	"github.com/israelalvesmelo/desafio-otel/internal/utils"
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
	if err := utils.CEPValidation(cep); err != nil {
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
		utils.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TemperatureHandler) handlerError(w http.ResponseWriter, err error) {
	fmt.Println("error:", err)

	switch {
	case errors.Is(err, entity.ErrZipcodeNotValid):
		utils.Error(w, entity.ErrZipcodeNotValid.Error(), http.StatusUnprocessableEntity)
		return
	case errors.Is(err, entity.ErrZipcodeNotFound):
		utils.Error(w, entity.ErrZipcodeNotFound.Error(), http.StatusNotFound)
		return
	case err != nil:
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

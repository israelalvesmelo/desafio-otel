package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/israelalvesmelo/desafio-otel/internal/inputhandle/domain/dto"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/entity"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/usecase"
	"github.com/israelalvesmelo/desafio-otel/internal/utils"
)

type InputHandler struct {
}

func NewInputHandler(useCase *usecase.GetTemperatureUseCase) *InputHandler {
	return &InputHandler{}
}

func (h *InputHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bodyBytes, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		utils.Error(w, readErr.Error(), http.StatusBadRequest)
		return
	}
	var location dto.LocationInput
	if unmErr := json.Unmarshal(bodyBytes, &location); unmErr != nil {
		utils.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if err := utils.CEPValidation(location.Cep); err != nil {
		utils.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	reqCtx, reqCtxErr := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://service_b:50055/temperature?cep="+location.Cep,
		nil,
	)
	if reqCtxErr != nil {
		utils.Error(w, reqCtxErr.Error(), http.StatusInternalServerError)
		return
	}

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	clientDo, doErr := httpClient.Do(reqCtx)
	if doErr != nil {
		utils.Error(w, doErr.Error(), http.StatusInternalServerError)
		return
	}
	defer clientDo.Body.Close()
	h.handlerResponse(w, clientDo)
}

func (h *InputHandler) handlerResponse(w http.ResponseWriter, resp *http.Response) {
	switch resp.StatusCode {
	case http.StatusUnprocessableEntity:
		utils.Error(w, entity.ErrZipcodeNotValid.Error(), http.StatusUnprocessableEntity)
		return
	case http.StatusNotFound:
		utils.Error(w, entity.ErrZipcodeNotFound.Error(), http.StatusNotFound)
		return
	}

	locRespBody, bReadErr := io.ReadAll(resp.Body)
	if bReadErr != nil {
		utils.Error(w, bReadErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var locTempResp dto.TemperatureOutput
	if err := json.Unmarshal(locRespBody, &locTempResp); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, marshErr := json.Marshal(dto.TemperatureOutput{
		City:  locTempResp.City,
		TempC: locTempResp.TempC,
		TempF: locTempResp.TempF,
		TempK: locTempResp.TempK,
	})
	if marshErr != nil {
		utils.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	if _, wErr := w.Write(jsonData); wErr != nil {
		utils.Error(w, "Error writing", http.StatusInternalServerError)
		return
	}
}

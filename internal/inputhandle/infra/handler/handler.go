package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/israelalvesmelo/desafio-otel/cmd/config"
	"github.com/israelalvesmelo/desafio-otel/infra/tracer"
	"github.com/israelalvesmelo/desafio-otel/internal/inputhandle/domain/dto"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/entity"
	"github.com/israelalvesmelo/desafio-otel/internal/utils"
)

type InputHandler struct {
	config       *config.Config
	tracerHelper *tracer.TracerHelper
}

func NewInputHandler(config *config.Config, tracerHelper *tracer.TracerHelper) *InputHandler {
	return &InputHandler{
		config:       config,
		tracerHelper: tracerHelper,
	}
}

func (h *InputHandler) PostWeather(w http.ResponseWriter, r *http.Request) {
	ctx := h.tracerHelper.ExtractContext(r)
	ctx, span := h.tracerHelper.StartSpan(ctx)
	defer span.End()

	location, err := h.readAndParseLocation(r)
	if err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.CEPValidation(location.Cep); err != nil {
		utils.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	req, err := h.buildTemperatureRequest(ctx, location.Cep)
	if err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := h.doTemperatureRequest(req)
	if err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	h.handleTemperatureResponse(w, resp)
}

func (h *InputHandler) readAndParseLocation(r *http.Request) (*dto.LocationInput, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var location = &dto.LocationInput{}
	if err := json.Unmarshal(bodyBytes, &location); err != nil {
		return nil, errors.New("Error parsing JSON")
	}

	return location, nil
}

func (h *InputHandler) buildTemperatureRequest(ctx context.Context, cep string) (*http.Request, error) {
	url := h.config.ServiceB.Host + "/temperature?cep=" + cep

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	h.tracerHelper.InjectContext(ctx, req)
	return req, nil
}

func (h *InputHandler) doTemperatureRequest(req *http.Request) (*http.Response, error) {
	client := http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func (h *InputHandler) handleTemperatureResponse(w http.ResponseWriter, resp *http.Response) {
	switch resp.StatusCode {
	case http.StatusUnprocessableEntity:
		utils.Error(w, entity.ErrZipcodeNotValid.Error(), http.StatusUnprocessableEntity)
		return
	case http.StatusNotFound:
		utils.Error(w, entity.ErrZipcodeNotFound.Error(), http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var temp dto.TemperatureOutput
	if err := json.Unmarshal(body, &temp); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, temp)
}

func (h *InputHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		utils.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsonData); err != nil {
		utils.Error(w, "Error writing", http.StatusInternalServerError)
	}
}

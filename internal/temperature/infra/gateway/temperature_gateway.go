package gateway

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/israelalvesmelo/desafio-otel/cmd/config"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/dto"
	gatewaydomain "github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/gateway"
)

type TemperatureGatewayImpl struct {
	config *config.Temperature
}

func NewTemperatureGateway(config *config.Temperature) gatewaydomain.TemperatureGateway {
	return TemperatureGatewayImpl{
		config: config,
	}
}

var createWeatherEndpoint = func(baseUrl string) string {
	return strings.Join([]string{baseUrl, "v1", "current.json"}, "/")
}

func (g TemperatureGatewayImpl) GetTempCelsius(ctx context.Context, location string) (*float64, error) {
	u, urlErr := url.Parse(createWeatherEndpoint(g.config.URL))
	if urlErr != nil {
		fmt.Printf("Error parsing URL: %s\n", urlErr)
		return nil, urlErr
	}
	apiKey := g.config.ApiKey
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}

	q := u.Query()
	q.Set("key", g.config.ApiKey)
	q.Set("q", location)
	q.Set("aqi", "no")
	u.RawQuery = q.Encode()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if reqErr != nil {
		fmt.Printf("Error creating request: %s\n", reqErr)
		return nil, urlErr
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, doErr := client.Do(req)
	if doErr != nil {
		fmt.Printf("Error making GET request: %s\n", doErr)
		return nil, doErr
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Printf("Error reading response body: %s\n", readErr)
		return nil, readErr
	}

	var weatherData dto.TemperatureResponseOut
	if unmErr := json.Unmarshal(bodyBytes, &weatherData); unmErr != nil {
		fmt.Printf("Error parsing JSON: %s\n", unmErr)
		return nil, unmErr
	}

	return &weatherData.Current.TempC, nil
}

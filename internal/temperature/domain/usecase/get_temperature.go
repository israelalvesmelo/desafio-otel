package usecase

import (
	"context"
	"fmt"

	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/dto"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/entity"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/gateway"
)

type GetTemperatureUseCase struct {
	cepGateway         gateway.LocationGateway
	temperatureGateway gateway.TemperatureGateway
}

func NewGetTemperatureUseCase(
	cepGateway gateway.LocationGateway,
	temperatureGateway gateway.TemperatureGateway,
) *GetTemperatureUseCase {
	return &GetTemperatureUseCase{
		cepGateway:         cepGateway,
		temperatureGateway: temperatureGateway,
	}
}

func (uc *GetTemperatureUseCase) Execute(ctx context.Context, cep string) (*dto.TemperatureOutput, error) {
	location, err := uc.cepGateway.GetLocation(ctx, cep)
	if err != nil {
		return nil, err
	}

	fmt.Println("location: ", *location)

	temperatureCelsius, err := uc.temperatureGateway.GetTempCelsius(ctx, location.Localidade)
	if err != nil {
		return nil, err
	}
	fmt.Println("temperature celsius: ", *temperatureCelsius)
	temperature := entity.NewTemperature(*temperatureCelsius)

	return &dto.TemperatureOutput{
		City:  location.Localidade,
		TempC: temperature.Celsius(),
		TempF: temperature.Fahrenheit(),
		TempK: temperature.Kelvin(),
	}, nil

}

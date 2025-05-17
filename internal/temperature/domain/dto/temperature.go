package dto

type TemperatureResponseOut struct {
	Current TemperatureCurrent `json:"current"`
}

type TemperatureCurrent struct {
	TempC float64 `json:"temp_c"`
}

type TemperatureOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

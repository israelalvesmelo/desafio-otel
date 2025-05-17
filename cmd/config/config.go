package config

type Config struct {
	ServiceA    ServiceA
	ServiceB    ServiceB
	Temperature Temperature
	CEP         CEP
}

type ServiceB struct {
	Port string
	Host string
}

type ServiceA struct {
	Port string
}

type Temperature struct {
	ApiKey string
	URL    string
}

type CEP struct {
	URL string
}

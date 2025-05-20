package config

type Config struct {
	ServiceA    ServiceA
	ServiceB    ServiceB
	Temperature Temperature
	CEP         CEP
	Zipkin      Zipkin
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

type Zipkin struct {
	Endpoint string
}

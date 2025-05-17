package dto

type LocationInput struct {
	Cep string `json:"cep"`
}

type LocationOutput struct {
	Localidade string `json:"localidade"`

	Erro string `json:"erro"`
}

func (l LocationOutput) HasError() bool {
	return l.Erro != ""
}
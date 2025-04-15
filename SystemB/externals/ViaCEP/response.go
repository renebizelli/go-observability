package viacep

type APIResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Uf          string `json:"uf"`
	Localidade  string `json:"localidade"`
	Erro        string `json:"erro"`
}

package weatherAPI

type APIResponse struct {
	Current Current `json:"current"`
}

type Current struct {
	Celsius    float32 `json:"temp_c"`
	Fahrenheit float32 `json:"temp_f"`
}

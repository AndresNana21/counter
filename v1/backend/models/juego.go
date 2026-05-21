package models

// ContadorResponse es el molde para enviar datos a Astro
type ContadorResponse struct {
	Valor int `json:"valor"`
}

// ActualizarRequest es el molde para recibir datos desde Astro
type ActualizarRequest struct {
	Valor int `json:"valor"`
}

package models

// TransmisionResponse representa la estructura de una respuesta de transmisión
type TransmisionResponse struct {
	TransmittedFiles []TransmittedFile `json:"transmittedFiles"`
}

// TransmittedFile representa un archivo transmitido y su resultado.
type TransmittedFile struct {
	FileName           string             `json:"fileName"`
	TransmissionResult TransmissionResult `json:"transmissionResult"`
}

// TransmissionResult contiene el resultado de la transmisión.
type TransmissionResult struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// Response estructura las respuestas del handler.
type Response struct {
	Message    string `json:"message"`
	ErrorCount int    `json:"error_count,omitempty"`
	TotalFiles int    `json:"total_files"`
	Success    bool   `json:"success"`
}

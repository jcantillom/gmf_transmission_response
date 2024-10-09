package models

// TransmisionResponse representa la estructura de una respuesta de transmisión
type TransmisionResponse struct {
	TransmittedFiles []TransmittedFile `json:"transmittedFiles"` // Array de archivos transmitidos
}

// TransmittedFile representa un archivo transmitido y su resultado.
type TransmittedFile struct {
	FileName           string             `json:"fileName"`           // Nombre del archivo transmitido
	TransmissionResult TransmissionResult `json:"transmissionResult"` // Contenedor de resultado de la transmisión
}

// TransmissionResult contiene el resultado de la transmisión.
type TransmissionResult struct {
	Status string `json:"status"` // Estado de la transmisión: ERROR o SUCCESSFUL
	Code   string `json:"code"`   // Código de error de longitud 4
	Detail string `json:"detail"` // Detalle del resultado de la transmisión
}

// Response estructura las respuestas del handler.
type Response struct {
	Message    string `json:"message"`
	ErrorCount int    `json:"error_count,omitempty"`
	TotalFiles int    `json:"total_files"`
	Success    bool   `json:"success"`
}

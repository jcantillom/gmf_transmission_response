package handler

import (
	"encoding/json"
	"fmt"
	"gmf_transmission_response/internal/logs"
	"gmf_transmission_response/internal/models"
	"gmf_transmission_response/internal/service"
	"net/http"
)

// ArchivoHandlerInterface define la interfaz para manejar transmisiones
type ArchivoHandlerInterface interface {
	HandleTransmisionResponses(w http.ResponseWriter, r *http.Request)
}

// ArchivoHandler maneja las solicitudes relacionadas con archivos.
type ArchivoHandler struct {
	ArchivoService service.ArchivoServiceInterface
}

// NewArchivoHandler crea una nueva instancia de ArchivoHandler.
func NewArchivoHandler(archivoService service.ArchivoServiceInterface) *ArchivoHandler { // Cambia esto a la interfaz
	return &ArchivoHandler{
		ArchivoService: archivoService,
	}
}

// HandleTransmisionResponses es el controlador para procesar el array de respuestas de transmisión.
// Recibe un array de transmisiones a través de API Gateway y procesa cada una de ellas.
func (h *ArchivoHandler) HandleTransmisionResponses(w http.ResponseWriter, r *http.Request) {
	var transmisionResponse models.TransmisionResponse

	// Decodificar el cuerpo de la solicitud JSON en la estructura de TransmisionResponse
	if err := json.NewDecoder(r.Body).Decode(&transmisionResponse); err != nil {
		logs.LogError(nil, "Error al decodificar el cuerpo de la solicitud: %v", err)
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	// Inicializar contadores
	var errorCount, successCount int

	// Procesar cada archivo transmitido en la solicitud
	for _, transmittedFile := range transmisionResponse.TransmittedFiles {
		// Procesar cada archivo utilizando el servicio
		if err := h.ArchivoService.ProcesarTransmision(transmittedFile); err != nil {
			// Si hay un error, registrar el error y continuar con el siguiente archivo
			logs.LogError(nil, "Error al procesar archivo transmitido: %v", err)
			errorCount++ // Incrementar el contador de errores
			continue     // Continuar con el siguiente archivo
		}
		successCount++ // Incrementar el contador de archivos procesados exitosamente
	}

	// Loggear el resultado del procesamiento
	logs.LogInfo(nil, "Archivos procesados correctamente: %d ✅", successCount)
	logs.LogWarn(nil, "Archivos con errores: %d ❌", errorCount)

	// Crear una respuesta en formato JSON con el resumen
	response := models.Response{
		TotalFiles: len(transmisionResponse.TransmittedFiles),
		ErrorCount: errorCount,
		Success:    errorCount == 0,
	}

	if errorCount > 0 {
		response.Message = fmt.Sprintf(
			"Se procesaron con errores %d de %d archivos", errorCount, len(transmisionResponse.TransmittedFiles))
		w.WriteHeader(http.StatusPartialContent)
	} else {
		response.Message = "Todos los archivos fueron procesados correctamente ✅"
		w.WriteHeader(http.StatusOK)
	}

	// Enviar la respuesta JSON
	json.NewEncoder(w).Encode(response)
}

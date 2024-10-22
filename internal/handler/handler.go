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

	if err := json.NewDecoder(r.Body).Decode(&transmisionResponse); err != nil {
		logs.Logger.LogError("Error al decodificar el cuerpo de la solicitud", err, "N/A")
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	var errorCount, successCount int

	for _, transmittedFile := range transmisionResponse.TransmittedFiles {
		fileName := transmittedFile.FileName
		if err := h.ArchivoService.ProcesarTransmision(transmittedFile); err != nil {
			logs.Logger.LogError("Error al procesar archivo transmitido", err, fileName)
			errorCount++
			continue
		}
		logs.Logger.LogInfo("Archivo procesado exitosamente", fileName)
		successCount++
	}

	logs.Logger.LogInfo(fmt.Sprintf("Archivos procesados correctamente: %d", successCount), "N/A")
	logs.Logger.LogWarn(fmt.Sprintf("Archivos con errores: %d", errorCount), "N/A")

	response := models.Response{
		TotalFiles: len(transmisionResponse.TransmittedFiles),
		ErrorCount: errorCount,
		Success:    errorCount == 0,
	}

	if errorCount > 0 {
		response.Message = fmt.Sprintf("Se procesaron con errores %d de %d archivos", errorCount, len(transmisionResponse.TransmittedFiles))
		w.WriteHeader(http.StatusPartialContent)
	} else {
		response.Message = "Todos los archivos fueron procesados correctamente"
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

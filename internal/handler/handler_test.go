package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gmf_transmission_response/internal/handler"
	"gmf_transmission_response/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockArchivoService es un mock que implementa la interfaz ArchivoServiceInterface
type MockArchivoService struct {
	ProcessError         bool  // Define si todas las llamadas fallan
	ErrorInSpecificCalls []int // Define en qué llamadas específicas debe fallar
	CallCount            int   // Cuenta cuántas veces se ha llamado al servicio
}

// ProcesarTransmision simula el procesamiento de transmisión y falla según lo indicado
func (m *MockArchivoService) ProcesarTransmision(transmittedFile models.TransmittedFile) error {
	m.CallCount++

	// Verificar si la llamada actual debe fallar
	for _, call := range m.ErrorInSpecificCalls {
		if m.CallCount == call {
			return fmt.Errorf("mock error")
		}
	}

	// Si ProcessError está activado, siempre falla
	if m.ProcessError {
		return fmt.Errorf("mock error")
	}

	return nil
}

// Otros métodos necesarios para cumplir con la interfaz
func (m *MockArchivoService) RemoveExtension(fileName string) string { return fileName }
func (m *MockArchivoService) IsAnulacion(fileName string) bool       { return false }
func (m *MockArchivoService) ValidateIDLength(id string) error       { return nil }

func TestHandleTransmisionResponses_Success(t *testing.T) {
	// Caso de prueba exitoso, sin errores
	mockService := &MockArchivoService{
		ProcessError: false, // No habrá errores
	}

	h := handler.NewArchivoHandler(mockService)

	// Crear el cuerpo de la solicitud con un archivo exitoso
	body := models.TransmisionResponse{
		TransmittedFiles: []models.TransmittedFile{
			{
				FileName: "TUTGMF000100012024031-0001.txt",
				TransmissionResult: models.TransmissionResult{
					Status: "SUCCESSFUL",
					Code:   "0000",
					Detail: "Transmisión exitosa",
				},
			},
		},
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/transmision", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	// Ejecutar el handler
	h.HandleTransmisionResponses(w, req)

	// Validar el resultado
	assert.Equal(t, http.StatusOK, w.Code) // El código de estado debe ser 200
	var resp models.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 1, resp.TotalFiles)
	assert.Equal(t, 0, resp.ErrorCount)
	assert.Equal(t, "Todos los archivos fueron procesados correctamente", resp.Message)
}

func TestHandleTransmisionResponses_PartialFailure(t *testing.T) {
	// Caso con error en la primera llamada
	mockService := &MockArchivoService{
		ErrorInSpecificCalls: []int{1}, // Solo falla en la primera llamada
	}

	h := handler.NewArchivoHandler(mockService)

	// Crear el cuerpo de la solicitud con dos archivos (uno fallido y uno exitoso)
	body := models.TransmisionResponse{
		TransmittedFiles: []models.TransmittedFile{
			{
				FileName: "TUTGMF000100012024031-0001.txt",
				TransmissionResult: models.TransmissionResult{
					Status: "ERROR",
					Code:   "0001",
					Detail: "Error en la transmisión",
				},
			},
			{
				FileName: "TUTGMF0001000120240312-0002-A.txt",
				TransmissionResult: models.TransmissionResult{
					Status: "SUCCESSFUL",
					Code:   "0000",
					Detail: "Transmisión exitosa",
				},
			},
		},
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/transmision", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	// Ejecutar el handler
	h.HandleTransmisionResponses(w, req)

	// Validar el resultado
	assert.Equal(t, http.StatusPartialContent, w.Code)
	var resp models.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, 2, resp.TotalFiles)                                        // Deben ser 2 archivos
	assert.Equal(t, 1, resp.ErrorCount)                                        // Debe haber 1 error
	assert.Equal(t, "Se procesaron con errores 1 de 2 archivos", resp.Message) // Mensaje esperado
}

func TestHandleTransmisionResponses_InvalidRequest(t *testing.T) {
	mockService := &MockArchivoService{
		ProcessError: false, // No habrá errores en el mock
	}

	h := handler.NewArchivoHandler(mockService)

	// Crear un request con un body inválido
	req := httptest.NewRequest(http.MethodPost, "/transmision", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	// Ejecutar el handler
	h.HandleTransmisionResponses(w, req)

	// Validar el resultado
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Solicitud inválida\n", w.Body.String())
}

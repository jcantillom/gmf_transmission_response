package service_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gmf_transmission_response/internal/models"
	"gmf_transmission_response/internal/service"
	"testing"
)

// MockRepository es un mock del repositorio que implementa RepositoryInterface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetArchivoByNombreArchivo(nombreArchivo string) (*models.CGDArchivo, error) {
	args := m.Called(nombreArchivo)
	if archivo, ok := args.Get(0).(*models.CGDArchivo); ok {
		return archivo, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) UpdateArchivo(archivo *models.CGDArchivo) error {
	return m.Called(archivo).Error(0)
}

func (m *MockRepository) InsertEstadoArchivo(estado *models.CGDArchivoEstado) error {
	return m.Called(estado).Error(0)
}

// Test de procesamiento de transmisión exitosa
func TestProcesarTransmision_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	archivoService := service.NewArchivoService(mockRepo)

	transmittedFile := models.TransmittedFile{
		FileName: "TUTGMF0001000120240312-0001",
		TransmissionResult: models.TransmissionResult{
			Status: "SUCCESSFUL",
			Code:   "0000",
			Detail: "Transmisión exitosa",
		},
	}

	archivo := &models.CGDArchivo{
		IDArchivo:         10001202403120001,
		GAWRtaTransEstado: "PENDING",
	}

	// Simular respuestas del mock
	mockRepo.On("GetArchivoByNombreArchivo", "TUTGMF0001000120240312-0001").Return(archivo, nil)
	mockRepo.On("UpdateArchivo", archivo).Return(nil)
	mockRepo.On("InsertEstadoArchivo", mock.Anything).Return(nil)

	err := archivoService.ProcesarTransmision(transmittedFile)

	// Validaciones
	assert.NoError(t, err)
	assert.Equal(t, "ENVIADO", archivo.Estado)
	mockRepo.AssertExpectations(t)
}

// Test de procesamiento de transmisión fallida por error en transmisión
func TestProcesarTransmision_Failure(t *testing.T) {
	mockRepo := new(MockRepository)
	archivoService := service.NewArchivoService(mockRepo)

	transmittedFile := models.TransmittedFile{
		FileName: "TUTGMF0001000120240312-0001",
		TransmissionResult: models.TransmissionResult{
			Status: "ERROR",
			Code:   "0001",
			Detail: "Error en la transmisión",
		},
	}

	archivo := &models.CGDArchivo{
		IDArchivo:         10001202403120001,
		GAWRtaTransEstado: "PENDING",
	}

	// Simular respuestas del mock
	mockRepo.On("GetArchivoByNombreArchivo", "TUTGMF0001000120240312-0001").Return(archivo, nil)
	mockRepo.On("UpdateArchivo", archivo).Return(nil)
	mockRepo.On("InsertEstadoArchivo", mock.Anything).Return(nil)

	err := archivoService.ProcesarTransmision(transmittedFile)

	// Validaciones
	assert.NoError(t, err)
	assert.Equal(t, "ENVIO_FALLIDO", archivo.Estado)
	mockRepo.AssertExpectations(t)
}

// Test de fallo por no encontrar el archivo en la base de datos
func TestProcesarTransmision_ArchivoNoEncontrado(t *testing.T) {
	mockRepo := new(MockRepository)
	archivoService := service.NewArchivoService(mockRepo)

	transmittedFile := models.TransmittedFile{
		FileName: "TUTGMF0001000120240312-0001",
		TransmissionResult: models.TransmissionResult{
			Status: "ERROR",
			Code:   "0001",
			Detail: "Error en la transmisión",
		},
	}

	// Simular que no se encuentra el archivo
	mockRepo.On("GetArchivoByNombreArchivo", "TUTGMF0001000120240312-0001").Return(nil, errors.New("archivo no encontrado"))

	err := archivoService.ProcesarTransmision(transmittedFile)

	// Validaciones
	assert.Error(t, err)
	assert.EqualError(t, err, "archivo no encontrado")
	mockRepo.AssertExpectations(t)
}

// Test de validación de longitud de ID inválido
func TestProcesarTransmision_ValidacionID(t *testing.T) {
	mockRepo := new(MockRepository)
	archivoService := service.NewArchivoService(mockRepo)

	archivo := &models.CGDArchivo{
		IDArchivo:         10001202403120001, // Un ID válido para `int64`
		GAWRtaTransEstado: "PENDING",
	}

	// Simulamos que el ID tiene longitud incorrecta (más de 16 dígitos) utilizando `strconv.FormatInt`.
	mockRepo.On("GetArchivoByNombreArchivo", "TUTGMF0001000120240312-0001").Return(archivo, nil)
	mockRepo.On("UpdateArchivo", archivo).Return(nil)
	mockRepo.On("InsertEstadoArchivo", mock.Anything).Return(nil)

	// Forzar que el archivo tenga un ID con más de 16 caracteres para la validación.
	err := archivoService.ValidateIDLength("1000120240312000123") // Cadena de longitud incorrecta

	// Validar que el error de longitud sea detectado
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "el ID del archivo debe tener una longitud de 16 caracteres")
	mockRepo.AssertNotCalled(t, "UpdateArchivo", mock.Anything)
	mockRepo.AssertNotCalled(t, "InsertEstadoArchivo", mock.Anything)
}
func TestProcesarTransmision_Anulacion(t *testing.T) {
	mockRepo := new(MockRepository)
	archivoService := service.NewArchivoService(mockRepo)

	archivo := &models.CGDArchivo{
		IDArchivo:         10001202403120001,
		GAWRtaTransEstado: "PENDING",
	}

	transmittedFile := models.TransmittedFile{
		FileName: "TUTGMF0001000120240312-0002-A", // Este nombre debe hacer que IsAnulacion devuelva true
		TransmissionResult: models.TransmissionResult{
			Status: "SUCCESSFUL",
			Code:   "0000",
			Detail: "Transmisión exitosa",
		},
	}

	mockRepo.On("GetArchivoByNombreArchivo", "TUTGMF0001000120240312-0002-A").Return(archivo, nil)
	mockRepo.On("UpdateArchivo", archivo).Return(nil)
	mockRepo.On("InsertEstadoArchivo", mock.Anything).Return(nil)

	err := archivoService.ProcesarTransmision(transmittedFile)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "GetArchivoByNombreArchivo", "TUTGMF0001000120240312-0002-A")
	mockRepo.AssertCalled(t, "UpdateArchivo", archivo)
	mockRepo.AssertCalled(t, "InsertEstadoArchivo", mock.Anything)
}

func TestProcesarTransmision_ErrorAlInsertarEstado(t *testing.T) {
	mockRepo := new(MockRepository)
	archivoService := service.NewArchivoService(mockRepo)

	archivo := &models.CGDArchivo{
		IDArchivo:         10001202403120001,
		GAWRtaTransEstado: "PENDING",
	}

	transmittedFile := models.TransmittedFile{
		FileName: "TUTGMF0001000120240312-0001",
		TransmissionResult: models.TransmissionResult{
			Status: "SUCCESSFUL",
			Code:   "0000",
			Detail: "Transmisión exitosa",
		},
	}

	mockRepo.On("GetArchivoByNombreArchivo", "TUTGMF0001000120240312-0001").Return(archivo, nil)
	mockRepo.On("UpdateArchivo", archivo).Return(nil)
	mockRepo.On("InsertEstadoArchivo", mock.Anything).Return(fmt.Errorf("mock error"))

	err := archivoService.ProcesarTransmision(transmittedFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock error")
	mockRepo.AssertCalled(t, "InsertEstadoArchivo", mock.Anything)
}

func TestProcesarTransmision_AnulacionFallida(t *testing.T) {
	mockRepo := new(MockRepository)
	archivoService := service.NewArchivoService(mockRepo)

	archivo := &models.CGDArchivo{
		IDArchivo:         10001202403120001,
		GAWRtaTransEstado: "PENDING",
	}

	transmittedFile := models.TransmittedFile{
		FileName: "TUTGMF0001000120240312-0002-A", // Anulación por el sufijo "-A"
		TransmissionResult: models.TransmissionResult{
			Status: "ERROR",
			Code:   "0001",
			Detail: "Error en la transmisión",
		},
	}

	mockRepo.On("GetArchivoByNombreArchivo", "TUTGMF0001000120240312-0002-A").Return(archivo, nil)
	mockRepo.On("UpdateArchivo", archivo).Return(nil)
	mockRepo.On("InsertEstadoArchivo", mock.Anything).Return(nil)

	err := archivoService.ProcesarTransmision(transmittedFile)

	assert.NoError(t, err)
	assert.Equal(t, "ANULACION_FALLIDA", archivo.Estado) // Comprobamos el estado
	mockRepo.AssertCalled(t, "UpdateArchivo", archivo)
}

func TestProcesarTransmision_ErrorAlActualizarArchivo(t *testing.T) {
	mockRepo := new(MockRepository)
	archivoService := service.NewArchivoService(mockRepo)

	archivo := &models.CGDArchivo{
		IDArchivo:         10001202403120001,
		GAWRtaTransEstado: "PENDING",
	}

	transmittedFile := models.TransmittedFile{
		FileName: "TUTGMF0001000120240312-0001",
		TransmissionResult: models.TransmissionResult{
			Status: "SUCCESSFUL",
			Code:   "0000",
			Detail: "Transmisión exitosa",
		},
	}

	mockRepo.On("GetArchivoByNombreArchivo", "TUTGMF0001000120240312-0001").Return(archivo, nil)
	mockRepo.On("UpdateArchivo", archivo).Return(fmt.Errorf("mock error"))
	mockRepo.On("InsertEstadoArchivo", mock.Anything).Return(nil)

	err := archivoService.ProcesarTransmision(transmittedFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock error")
	mockRepo.AssertCalled(t, "UpdateArchivo", archivo)
}

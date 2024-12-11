package models_test

import (
	"testing"
	"time"

	"gmf_transmission_response/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestCGDArchivoTableName(t *testing.T) {
	archivo := models.CGDArchivos{}
	assert.Equal(
		t, "cgd_archivos", archivo.TableName(), "El nombre de la tabla debe ser CGD_ARCHIVO")
}

func TestCGDArchivoEstadoTableName(t *testing.T) {
	archivoEstado := models.CGDArchivoEstados{}
	assert.Equal(t, "cgd_archivo_estados", archivoEstado.TableName(), "El nombre de la tabla debe ser CGD_ARCHIVO_ESTADO")
}

func TestCGDArchivoFields(t *testing.T) {
	fecha := time.Now()

	archivo := models.CGDArchivos{
		IDArchivo:                   1234567890123456,
		NombreArchivo:               "archivo.txt",
		PlataformaOrigen:            "AB",
		TipoArchivo:                 "TX",
		ConsecutivoPlataformaOrigen: 1,
		FechaNombreArchivo:          "20221012",
		NroTotalRegistros:           1000,
		NroRegistrosError:           5,
		NroRegistrosValidos:         995,
		Estado:                      "PROCESADO",
		FechaRecepcion:              fecha,
		FechaCiclo:                  fecha,
		ACGFechaGeneracion:          fecha,
		ACGConsecutivo:              100,
		ACGNombreArchivo:            "TUTGMF000100012024031-0001.txt",
		GAWRtaTransEstado:           "SUCCESS",
		GAWRtaTransCodigo:           "0000",
		GAWRtaTransDetalle:          "Transmisión exitosa",
	}

	// Verificación de campos clave
	assert.Equal(t, int64(1234567890123456), archivo.IDArchivo)
	assert.Equal(t, "archivo.txt", archivo.NombreArchivo)
	assert.Equal(t, "AB", archivo.PlataformaOrigen)
	assert.Equal(t, "TX", archivo.TipoArchivo)
	assert.Equal(t, "SUCCESS", archivo.GAWRtaTransEstado)
	assert.Equal(t, "0000", archivo.GAWRtaTransCodigo)
	assert.Equal(t, "Transmisión exitosa", archivo.GAWRtaTransDetalle)
}

func TestCGDArchivoEstadoFields(t *testing.T) {
	fecha := time.Now()

	archivoEstado := models.CGDArchivoEstados{
		IDArchivo:         1234567890123456,
		EstadoInicial:     "PENDIENTE",
		EstadoFinal:       "PROCESADO",
		FechaCambioEstado: fecha,
	}

	// Verificación de campos clave
	assert.Equal(t, int64(1234567890123456), archivoEstado.IDArchivo)
	assert.Equal(t, "PENDIENTE", archivoEstado.EstadoInicial)
	assert.Equal(t, "PROCESADO", archivoEstado.EstadoFinal)
	assert.Equal(t, fecha, archivoEstado.FechaCambioEstado)
}

func TestTransmisionResponseFields(t *testing.T) {
	// Crear instancias de ejemplo
	transmittedFiles := []models.TransmittedFile{
		{
			FileName: "TUTGMF000100012024031-0001.txt",
			TransmissionResult: models.TransmissionResult{
				Status: "SUCCESSFUL",
				Code:   "0000",
				Detail: "Transmisión exitosa",
			},
		},
	}

	transmisionResponse := models.TransmisionResponse{
		TransmittedFiles: transmittedFiles,
	}

	// Validar que los archivos transmitidos sean correctos
	assert.Equal(t, 1, len(transmisionResponse.TransmittedFiles))
	assert.Equal(t, "TUTGMF000100012024031-0001.txt", transmisionResponse.TransmittedFiles[0].FileName)
	assert.Equal(t, "SUCCESSFUL", transmisionResponse.TransmittedFiles[0].TransmissionResult.Status)
	assert.Equal(t, "0000", transmisionResponse.TransmittedFiles[0].TransmissionResult.Code)
	assert.Equal(t, "Transmisión exitosa", transmisionResponse.TransmittedFiles[0].TransmissionResult.Detail)
}

func TestResponseFields(t *testing.T) {
	// Crear instancia de ejemplo
	response := models.Response{
		Message:    "Todos los archivos fueron procesados correctamente ✅",
		ErrorCount: 0,
		TotalFiles: 1,
		Success:    true,
	}

	// Validar los valores de los campos
	assert.Equal(t, "Todos los archivos fueron procesados correctamente ✅", response.Message)
	assert.Equal(t, 0, response.ErrorCount)
	assert.Equal(t, 1, response.TotalFiles)
	assert.True(t, response.Success)
}

func TestTransmittedFileFields(t *testing.T) {
	// Crear instancia de ejemplo
	transmittedFile := models.TransmittedFile{
		FileName: "TUTGMF000100012024031-0002.txt",
		TransmissionResult: models.TransmissionResult{
			Status: "ERROR",
			Code:   "0001",
			Detail: "Error en la transmisión",
		},
	}

	// Validar los valores de los campos
	assert.Equal(t, "TUTGMF000100012024031-0002.txt", transmittedFile.FileName)
	assert.Equal(t, "ERROR", transmittedFile.TransmissionResult.Status)
	assert.Equal(t, "0001", transmittedFile.TransmissionResult.Code)
	assert.Equal(t, "Error en la transmisión", transmittedFile.TransmissionResult.Detail)
}

func TestTransmissionResultFields(t *testing.T) {
	// Crear instancia de ejemplo
	transmissionResult := models.TransmissionResult{
		Status: "SUCCESSFUL",
		Code:   "0000",
		Detail: "Transmisión exitosa",
	}

	// Validar los valores de los campos
	assert.Equal(t, "SUCCESSFUL", transmissionResult.Status)
	assert.Equal(t, "0000", transmissionResult.Code)
	assert.Equal(t, "Transmisión exitosa", transmissionResult.Detail)
}

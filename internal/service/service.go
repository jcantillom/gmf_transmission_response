package service

import (
	"fmt"
	"gmf_transmission_response/internal/logs"
	"gmf_transmission_response/internal/models"
	"gmf_transmission_response/internal/repository"
	"path/filepath"
	"strings"
	"time"
)

// ArchivoServiceInterface define los métodos que el servicio de archivos debe implementar.
type ArchivoServiceInterface interface {
	ProcesarTransmision(transmittedFile models.TransmittedFile) error
	RemoveExtension(fileName string) string
	IsAnulacion(fileName string) bool
	ValidateIDLength(id string) error
}

// ArchivoService implementa el servicio de archivos.
type ArchivoService struct {
	repo repository.RepositoryInterface
}

// NewArchivoService crea una nueva instancia de ArchivoService.
func NewArchivoService(repo repository.RepositoryInterface) *ArchivoService {
	return &ArchivoService{
		repo: repo,
	}
}

// ProcesarTransmision procesa una respuesta de transmisión (movimiento o anulación).
func (s *ArchivoService) ProcesarTransmision(transmittedFile models.TransmittedFile) error {
	fileName := transmittedFile.FileName
	isAnulacion := s.IsAnulacion(s.RemoveExtension(fileName))

	if isAnulacion {
		logs.Logger.LogInfo("La transmisión es una anulación", fileName)
	} else {
		logs.Logger.LogInfo("La transmisión es un movimiento", fileName)
	}

	archivo, err := s.repo.GetArchivoByNombreArchivo(fileName)
	if err != nil {
		logs.Logger.LogError("Error al obtener archivo de la base de datos", err, fileName)
		return err
	}

	if err := s.actualizarEstadoArchivo(archivo, transmittedFile, isAnulacion); err != nil {
		return err
	}

	estadoArchivo := &models.CGDArchivoEstado{
		IDArchivo:         archivo.IDArchivo,
		EstadoInicial:     archivo.GAWRtaTransEstado,
		EstadoFinal:       transmittedFile.TransmissionResult.Status,
		FechaCambioEstado: time.Now(),
	}

	if err := s.repo.InsertEstadoArchivo(estadoArchivo); err != nil {
		logs.Logger.LogError("Error al insertar estado del archivo", err, fileName)
		return err
	}

	logs.Logger.LogInfo("Estado insertado correctamente en la tabla CGD_ARCHIVO_ESTADO", fileName)
	return nil
}

// actualizarEstadoArchivo actualiza el estado del archivo dependiendo si es anulación o movimiento.
func (s *ArchivoService) actualizarEstadoArchivo(
	archivo *models.CGDArchivo, transmittedFile models.TransmittedFile, isAnulacion bool) error {
	// Actualizar el estado en función del resultado de la transmisión
	archivo.GAWRtaTransEstado = transmittedFile.TransmissionResult.Status
	archivo.GAWRtaTransCodigo = transmittedFile.TransmissionResult.Code
	archivo.GAWRtaTransDetalle = transmittedFile.TransmissionResult.Detail

	var filename = archivo.NombreArchivo

	if transmittedFile.TransmissionResult.Status == "ERROR" {
		if isAnulacion {
			archivo.Estado = "ANULACION_FALLIDA"
			logs.LogInfo("Se ha marcado el archivo en estado ANULACION_FALLIDA.", filename)
		} else {
			archivo.Estado = "ENVIO_FALLIDO"
			logs.LogInfo("Se ha marcado el archivo en estado ENVIO_FALLIDO.", filename)
		}
	} else {
		// Si la transmisión fue exitosa
		if isAnulacion {
			archivo.Estado = "ANULACION_ENVIADA"
			logs.LogInfo("Se ha marcado el archivo en estado ANULACION_ENVIADA.", filename)
		} else {
			archivo.Estado = "ENVIADO"
			logs.LogInfo("Se ha marcado el archivo en estado ENVIADO.", filename)
		}
	}

	// Actualizar el archivo en la base de datos
	if err := s.repo.UpdateArchivo(archivo); err != nil {
		logs.LogError("Error al actualizar el archivo en la base de datos: %v", err, filename)
		return err
	}

	return nil
}

// removeExtension elimina la extensión de un nombre de archivo.
func (s *ArchivoService) RemoveExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// isAnulacion valida si el archivo es de tipo "anulación" verificando si tiene "-A" en el nombre.
func (s *ArchivoService) IsAnulacion(fileName string) bool {
	return strings.HasSuffix(fileName, "-A")
}

// validateIDLength valida que el ID del archivo tenga una longitud de 16 caracteres.
func (s *ArchivoService) ValidateIDLength(id string) error {
	if len(id) != 16 {
		return fmt.Errorf("el ID del archivo debe tener una longitud de 16 caracteres")
	}
	return nil
}

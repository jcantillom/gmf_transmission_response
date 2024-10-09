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
	// Eliminar la extensión del nombre del archivo
	fileNameWithoutExtension := s.RemoveExtension(transmittedFile.FileName)
	// Validar si la respuesta es una anulación basándose en el nombre del archivo.
	isAnulacion := s.IsAnulacion(fileNameWithoutExtension)
	if isAnulacion {
		logs.LogInfo(
			nil, "La transmisión es una anulación 🅰️ para el archivo: %s", transmittedFile.FileName)
	} else {
		logs.LogInfo(
			nil, "La transmisión es un movimiento Ⓜ️ para el archivo: %s", transmittedFile.FileName)
	}

	// Buscar el archivo en la base de datos por el nombre del archivo (ACGNombreArchivo)
	archivo, err := s.repo.GetArchivoByNombreArchivo(fileNameWithoutExtension)
	if err != nil {
		logs.LogError(nil, "Error al obtener el archivo desde la base de datos: %v", err)
		return err
	}

	//// Validar la longitud del ID del archivo
	//if err := s.ValidateIDLength(strconv.FormatInt(archivo.IDArchivo, 10)); err != nil {
	//	logs.LogError(nil, "Error en la longitud del ID del archivo: %v", err)
	//}

	// Actualizar el estado del archivo en función de la transmisión
	if err := s.actualizarEstadoArchivo(archivo, transmittedFile, isAnulacion); err != nil {
		return err
	}

	// Insertar un nuevo estado en la tabla CGD_ARCHIVO_ESTADO
	estadoArchivo := &models.CGDArchivoEstado{
		IDArchivo:         archivo.IDArchivo,
		EstadoInicial:     archivo.GAWRtaTransEstado,
		EstadoFinal:       transmittedFile.TransmissionResult.Status,
		FechaCambioEstado: time.Now(),
	}
	if err := s.repo.InsertEstadoArchivo(estadoArchivo); err != nil {
		logs.LogError(
			nil, "Error al insertar el estado del archivo en la tabla CGD_ARCHIVO_ESTADO: %v", err)
		return err
	}
	logs.LogInfo(
		nil, "Se ha insertado un nuevo estado para el archivo {ID: %d} en la tabla CGD_ARCHIVO_ESTADO. 📝",
		archivo.IDArchivo)
	return nil
}

// actualizarEstadoArchivo actualiza el estado del archivo dependiendo si es anulación o movimiento.
func (s *ArchivoService) actualizarEstadoArchivo(
	archivo *models.CGDArchivo, transmittedFile models.TransmittedFile, isAnulacion bool) error {
	// Actualizar el estado en función del resultado de la transmisión
	archivo.GAWRtaTransEstado = transmittedFile.TransmissionResult.Status
	archivo.GAWRtaTransCodigo = transmittedFile.TransmissionResult.Code
	archivo.GAWRtaTransDetalle = transmittedFile.TransmissionResult.Detail

	if transmittedFile.TransmissionResult.Status == "ERROR" {
		if isAnulacion {
			archivo.Estado = "ANULACION_FALLIDA"
			logs.LogInfo(nil, "Se ha marcado el archivo en estado ANULACION_FALLIDA.")
		} else {
			archivo.Estado = "ENVIO_FALLIDO"
			logs.LogInfo(nil, "Se ha marcado el archivo en estado ENVIO_FALLIDO.")
		}
	} else {
		// Si la transmisión fue exitosa
		if isAnulacion {
			archivo.Estado = "ANULACION_ENVIADA"
			logs.LogInfo(nil, "Se ha marcado el archivo en estado ANULACION_ENVIADA.")
		} else {
			archivo.Estado = "ENVIADO"
			logs.LogInfo(nil, "Se ha marcado el archivo en estado ENVIADO.")
		}
	}

	// Actualizar el archivo en la base de datos
	if err := s.repo.UpdateArchivo(archivo); err != nil {
		logs.LogError(nil, "Error al actualizar el archivo en la base de datos: %v", err)
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

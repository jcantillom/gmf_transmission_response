package repository

import (
	"gmf_transmission_response/internal/models"
	"gorm.io/gorm"
)

// RepositoryInterface methods de la interfaz GormArchivoRepository
type RepositoryInterface interface {
	GetArchivoByNombreArchivo(nombreArchivo string) (*models.CGDArchivo, error)
	UpdateArchivo(archivo *models.CGDArchivo) error
	InsertEstadoArchivo(estado *models.CGDArchivoEstado) error
}

// GormArchivoRepository implementa el repositorio de Archivo utilizando GORM.
type GormArchivoRepository struct {
	DB *gorm.DB
}

// NewArchivoRepository crea una nueva instancia de GormArchivoRepository
func NewArchivoRepository(db *gorm.DB) *GormArchivoRepository {
	return &GormArchivoRepository{
		DB: db,
	}
}

// GetArchivoByNombreArchivo obtiene un archivo por su nombre de archivo (ACGNombreArchivo).
func (r *GormArchivoRepository) GetArchivoByNombreArchivo(nombreArchivo string) (*models.CGDArchivo, error) {
	var archivo models.CGDArchivo
	if err := r.DB.Where(
		"acg_nombre_archivo = ?", nombreArchivo).First(&archivo).Error; err != nil {
		return nil, err
	}
	return &archivo, nil
}

// UpdateArchivo actualiza el archivo en la base de datos con el nuevo estado de la transmisión.
func (r *GormArchivoRepository) UpdateArchivo(archivo *models.CGDArchivo) error {
	return r.DB.Model(&archivo).Updates(map[string]interface{}{
		"gaw_rta_trans_estado":  archivo.GAWRtaTransEstado,
		"gaw_rta_trans_codigo":  archivo.GAWRtaTransCodigo,
		"gaw_rta_trans_detalle": archivo.GAWRtaTransDetalle,
		"estado":                archivo.Estado,
	}).Error
}

// InsertEstadoArchivo inserta un nuevo estado en la tabla CGD_ARCHIVO_ESTADO.
func (r *GormArchivoRepository) InsertEstadoArchivo(estado *models.CGDArchivoEstado) error {
	return r.DB.Create(estado).Error
}

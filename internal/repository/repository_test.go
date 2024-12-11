package repository_test

import (
	"gmf_transmission_response/internal/models"
	"gmf_transmission_response/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Configuración inicial del mock de la base de datos
func SetupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestGetArchivoByNombreArchivo(t *testing.T) {
	// Configurar la base de datos de prueba y el mock
	gormDB, mock := SetupTestDB(t)
	repo := repository.NewArchivoRepository(gormDB)

	// Definir el archivo que será devuelto por la consulta
	nombreArchivo := "TUTGMF000100012024031-0001.txt"
	archivo := models.CGDArchivos{
		IDArchivo:         1,
		NombreArchivo:     nombreArchivo,
		GAWRtaTransEstado: "PENDING",
	}

	// Configurar el mock para la consulta SQL
	mock.ExpectQuery(
		`SELECT \* FROM "cgd_archivos" WHERE acg_nombre_archivo = \$1 ORDER BY "cgd_archivos"."id_archivo" LIMIT \$2`).
		WithArgs(nombreArchivo, 1). // Debemos pasar también el argumento para el límite (GORM lo agrega automáticamente)
		WillReturnRows(sqlmock.NewRows([]string{"id_archivo", "nombre_archivo", "gaw_rta_trans_estado"}).
			AddRow(archivo.IDArchivo, archivo.NombreArchivo, archivo.GAWRtaTransEstado))

	// Ejecutar el método
	result, err := repo.GetArchivoByNombreArchivo(nombreArchivo)

	// Verificar los resultados
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, archivo.NombreArchivo, result.NombreArchivo)
	assert.NoError(t, mock.ExpectationsWereMet()) // Verificar que todas las expectativas fueron cumplidas
}

func TestGetArchivoByNombreArchivo_NotFound(t *testing.T) {
	// Configurar la base de datos de prueba y el mock
	gormDB, mock := SetupTestDB(t)
	repo := repository.NewArchivoRepository(gormDB)

	// Definir el nombre de archivo que será buscado
	nombreArchivo := "TUTGMF000100012024031-0001.txt"

	// Configurar el mock para la consulta SQL
	mock.ExpectQuery(
		`SELECT \* FROM "cgd_archivos" WHERE acg_nombre_archivo = \$1 ORDER BY "cgd_archivos"."id_archivo" LIMIT \$2`).
		WithArgs(nombreArchivo, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id_archivo", "nombre_archivo", "gaw_rta_trans_estado"}))

	// Ejecutar el método
	result, err := repo.GetArchivoByNombreArchivo(nombreArchivo)

	// Verificar los resultados
	assert.Error(t, err)
	assert.Nil(t, result)

}

func TestUpdateArchivo(t *testing.T) {
	// Configurar la base de datos de prueba y el mock
	gormDB, mock := SetupTestDB(t)
	repo := repository.NewArchivoRepository(gormDB)

	// Datos de prueba
	archivo := &models.CGDArchivos{
		IDArchivo:          1,
		GAWRtaTransEstado:  "SUCCESS",
		GAWRtaTransCodigo:  "0000",
		GAWRtaTransDetalle: "Transmisión exitosa",
		Estado:             "ENVIADO",
	}

	// Agregar expectativas para el inicio de la transacción
	mock.ExpectBegin()

	// Configurar el mock para la consulta SQL de actualización
	mock.ExpectExec(`UPDATE "cgd_archivos"`).
		WithArgs(
			archivo.Estado,
			archivo.GAWRtaTransCodigo,
			archivo.GAWRtaTransDetalle,
			archivo.GAWRtaTransEstado,
			archivo.IDArchivo,
		).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Simular éxito en la actualización

	// Agregar expectativas para el commit de la transacción
	mock.ExpectCommit()

	// Ejecutar el método
	err := repo.UpdateArchivo(archivo)

	// Verificar que no hubo error
	assert.NoError(t, err)

	// Verificar que las expectativas del mock se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertEstadoArchivo(t *testing.T) {
	// Configurar la base de datos de prueba y el mock
	gormDB, mock := SetupTestDB(t)
	repo := repository.NewArchivoRepository(gormDB)

	// Definir el estado de archivo que será insertado
	estadoArchivo := models.CGDArchivoEstados{
		IDArchivo:         1,
		EstadoInicial:     "PENDING",
		EstadoFinal:       "SUCCESS",
		FechaCambioEstado: time.Now(),
	}

	// Configurar el mock para la consulta SQL de inserción
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "cgd_archivo_estados"`).
		WithArgs(
			estadoArchivo.IDArchivo,
			estadoArchivo.EstadoInicial,
			estadoArchivo.EstadoFinal,
			sqlmock.AnyArg()). // sqlmock.AnyArg para la fecha
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Ejecutar el método
	err := repo.InsertEstadoArchivo(&estadoArchivo)

	// Verificar los resultados
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet()) // Verificar que todas las expectativas fueron cumplidas
}

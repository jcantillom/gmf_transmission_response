package config

import (
	"context"
	"log"

	"gmf_transmission_response/connection"
	"gmf_transmission_response/internal/handler"
	"gmf_transmission_response/internal/logs"
	"gmf_transmission_response/internal/repository"
	"gmf_transmission_response/internal/service"
)

var ctx = context.TODO()

// InitApplication inicializa todos los componentes necesarios para la aplicación.
func InitApplication() (*handler.ArchivoHandler, *connection.DBManager) {
	// Inicializar el ConfigManager y cargar la configuración
	configManager := NewConfigManager()
	configManager.InitConfig()

	// Inicializar el DBManager y abrir la conexión a la base de datos
	dbManager := connection.NewDBManager()
	if err := dbManager.InitDB(); err != nil {
		log.Fatalf("Error inicializando la base de datos: %v", err)
	}

	// Inicializar el repositorio GORM con la conexión a la base de datos
	repo := repository.NewArchivoRepository(dbManager.GetDB())

	// Inicializar el servicio de archivos con el repositorio
	archivoService := service.NewArchivoService(repo)

	// Inicializar el handler de archivos
	archivoHandler := handler.NewArchivoHandler(archivoService)

	logs.LogInfo(ctx, "Aplicación inicializada correctamente 🚀")

	return archivoHandler, dbManager
}

// CleanupApplication maneja la limpieza de recursos, como cerrar conexiones a la base de datos.
func CleanupApplication(dbManager connection.DBManagerInterface) {
	// Cerrar la conexión a la base de datos
	dbManager.CloseDB()
	logs.LogInfo(ctx, "Recursos limpiados correctamente 🧹")
}

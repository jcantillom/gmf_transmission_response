package config

import (
	"log"

	"gmf_transmission_response/connection"
	"gmf_transmission_response/internal/handler"
	"gmf_transmission_response/internal/logs"
	"gmf_transmission_response/internal/repository"
	"gmf_transmission_response/internal/service"
)

// InitApplication inicializa todos los componentes necesarios para la aplicación.
func InitApplication() (*handler.ArchivoHandler, *connection.DBManager) {
	// Inicializar el ConfigManager y cargar la configuración
	configManager := NewConfigManager()
	configManager.InitConfig()

	// Inicializar el DBManager y abrir la conexión a la base de datos
	dbManager := connection.NewDBManager()
	if err := dbManager.InitDB(); err != nil {
		logs.Logger.LogError("Error inicializando la base de datos", err, "APP_INIT")
		log.Fatalf("Error inicializando la base de datos: %v", err)
	}

	// Inicializar el repositorio con la conexión a la base de datos
	repo := repository.NewArchivoRepository(dbManager.GetDB())

	// Inicializar el servicio de archivos con el repositorio
	archivoService := service.NewArchivoService(repo)

	// Inicializar el handler de archivos
	archivoHandler := handler.NewArchivoHandler(archivoService)

	logs.Logger.LogInfo("Aplicación inicializada correctamente ✅ ", "APP_INIT")

	return archivoHandler, dbManager
}

// CleanupApplication maneja la limpieza de recursos, como cerrar conexiones a la base de datos.
func CleanupApplication(dbManager connection.DBManagerInterface) {
	dbManager.CloseDB()
	logs.Logger.LogInfo("Recursos limpiados correctamente 🧹", "APP_CLEANUP")
}

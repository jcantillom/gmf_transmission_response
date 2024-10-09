package main

import (
	"gmf_transmission_response/config"
	"gmf_transmission_response/connection"
	"gmf_transmission_response/internal/handler"
	"gmf_transmission_response/internal/routes"
	"os"
)

var (
	archivoHandler *handler.ArchivoHandler
	dbManager      *connection.DBManager
)

func init() {
	// Inicializar la aplicación con todos los componentes
	archivoHandler, dbManager = config.InitApplication()
}

func main() {
	// configurar las rutas de la aplicación
	routes.SetupRoutes(archivoHandler)

	// Iniciar el servidor HTTP
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	connection.StartServer(host, port)

	// Limpiar los recursos de la aplicación
	defer config.CleanupApplication(dbManager)
}

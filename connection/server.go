package connection

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gmf_transmission_response/internal/logs"
)

// StartServer inicia el servidor HTTP en el host y puerto proporcionados.
func StartServer(host string, portStr string) {
	// Convertir el puerto a int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logs.Logger.LogError("Error al convertir el puerto", err, "SERVER_START")
		log.Fatalf("Error al convertir el puerto: %v", err)
	}

	address := fmt.Sprintf("%s:%d", host, port)

	// Mostrar el mensaje de que el servidor está iniciando
	logs.Logger.LogInfo(fmt.Sprintf(
		"Servidor iniciado exitosamente en http://%s 🚀", address), "SERVER_START")

	// Iniciar el servidor HTTP
	if err := http.ListenAndServe(address, nil); err != nil {
		logs.Logger.LogError("Error al iniciar el servidor", err, "SERVER_START")
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

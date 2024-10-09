package connection

import (
	"fmt"
	"gmf_transmission_response/internal/logs"
	"log"
	"net/http"
	"strconv"
)

// StartServer inicia el servidor HTTP en el host y puerto proporcionados.
func StartServer(host string, portStr string) {
	// Convertir el puerto a int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error al convertir el puerto: %v", err)
	}
	address := fmt.Sprintf("%s:%d", host, port)

	// Mostrar el mensaje de que el servidor está iniciando
	logs.LogInfo(nil, fmt.Sprintf(
		"Servidor iniciado exitosamente en http://%s 🚀", address))

	// Iniciar el servidor HTTP
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

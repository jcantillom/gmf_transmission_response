package routes

import (
	"gmf_transmission_response/internal/handler"
	"net/http"
)

func SetupRoutes(archivoHandle handler.ArchivoHandlerInterface) {
	http.HandleFunc("/transmission", func(w http.ResponseWriter, r *http.Request) {
		archivoHandle.HandleTransmisionResponses(w, r)
	})
}

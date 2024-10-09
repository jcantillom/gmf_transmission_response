package routes_test

import (
	"gmf_transmission_response/internal/routes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockArchivoHandler simula el comportamiento de ArchivoHandler
type MockArchivoHandler struct{}

func (m *MockArchivoHandler) HandleTransmisionResponses(w http.ResponseWriter, r *http.Request) {
	// Simulación de una respuesta básica para el mock
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Mock response"}`))
}

func TestSetupRoutes(t *testing.T) {
	// Crear un mock de ArchivoHandler
	mockHandler := &MockArchivoHandler{}

	// Configurar las rutas usando la función SetupRoutes
	routes.SetupRoutes(mockHandler)

	// Crear una solicitud HTTP POST (ya que la ruta maneja POST)
	req, err := http.NewRequest("POST", "/transmission", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Grabar la respuesta
	rr := httptest.NewRecorder()
	handler := http.DefaultServeMux

	// Ejecutar la solicitud
	handler.ServeHTTP(rr, req)

	// Verificar el código de estado de la respuesta
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verificar el contenido de la respuesta
	expected := `{"message": "Mock response"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

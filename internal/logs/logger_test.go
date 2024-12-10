package logs

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	// Captura la salida estándar
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Ejecuta la función
	f()

	// Restaura la salida estándar
	w.Close()
	os.Stdout = old

	// Lee el resultado
	var buf strings.Builder
	io.Copy(&buf, r)
	return buf.String()
}

func TestLogInfo(t *testing.T) {
	// Establece la variable de entorno LOG_FORMAT para activar el formato JSON
	t.Setenv("LOG_FORMAT", "JSON")

	// Ejecuta la función LogInfo
	message := "Este es un mensaje de información"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		LogInfo(message, fileName)
	})

	// Verifica que el nivel en el JSON sea "INFO"
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("Se esperaba el nivel \"level\":\"INFO\" en la salida, pero no se encontró. Salida: %s", output)
	}

	// Verifica que el mensaje esté presente
	if !strings.Contains(output, `"message":"`+message+`"`) {
		t.Errorf("Se esperaba el mensaje \"%s\" en la salida, pero no se encontró. Salida: %s", message, output)
	}
}

func TestLogWarn(t *testing.T) {
	// Establece la variable de entorno LOG_FORMAT para activar el formato JSON
	t.Setenv("LOG_FORMAT", "JSON")

	// Ejecuta la función LogWarn
	message := "Este es un mensaje de advertencia"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		LogWarn(message, fileName)
	})

	// Verifica que el nivel en el JSON sea "WARNING"
	if !strings.Contains(output, `"level":"WARNING"`) {
		t.Errorf("Se esperaba el nivel \"level\":\"WARNING\" en la salida, pero no se encontró. Salida: %s", output)
	}

	// Verifica que el mensaje esté presente
	if !strings.Contains(output, `"message":"`+message+`"`) {
		t.Errorf("Se esperaba el mensaje \"%s\" en la salida, pero no se encontró. Salida: %s", message, output)
	}
}

func TestLogError(t *testing.T) {
	// Establece la variable de entorno LOG_FORMAT para activar el formato JSON
	t.Setenv("LOG_FORMAT", "JSON")

	// Ejecuta la función LogError con un error
	message := "Este es un error con detalles"
	fileName := "logger_test.go"
	err := fmt.Errorf("file does not exist")
	output := captureOutput(func() {
		LogError(message, err, fileName)
	})

	// Verifica que el nivel en el JSON sea "ERROR"
	if !strings.Contains(output, `"level":"ERROR"`) {
		t.Errorf("Se esperaba el nivel \"level\":\"ERROR\" en la salida, pero no se encontró. Salida: %s", output)
	}

	// Verifica que el mensaje de error esté presente
	if !strings.Contains(output, `"message":"`+message+` - Error: file does not exist"`) {
		t.Errorf("Se esperaba el mensaje \"%s - Error: file does not exist\" en la salida, pero no se encontró. Salida: %s", message, output)
	}
}

func TestLogDebug(t *testing.T) {
	// Establece la variable de entorno LOG_LEVEL y LOG_FORMAT para activar el formato JSON
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("LOG_FORMAT", "JSON")

	// Ejecuta la función LogDebug
	message := "Este es un mensaje de depuración"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		LogDebug(message, fileName)
	})

	// Verifica que el nivel en el JSON sea "DEBUG"
	if !strings.Contains(output, `"level":"DEBUG"`) {
		t.Errorf("Se esperaba el nivel \"level\":\"DEBUG\" en la salida, pero no se encontró. Salida: %s", output)
	}

	// Verifica que el mensaje esté presente
	if !strings.Contains(output, `"message":"`+message+`"`) {
		t.Errorf("Se esperaba el mensaje \"%s\" en la salida, pero no se encontró. Salida: %s", message, output)
	}
}

func TestLoggerAdapterLogInfo(t *testing.T) {
	// Establece la variable de entorno LOG_FORMAT para activar el formato JSON
	t.Setenv("LOG_FORMAT", "JSON")

	// Ejecuta la función LogInfo usando el adaptador
	message := "Este es un mensaje de información desde el adaptador"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		Logger.LogInfo(message, fileName)
	})

	// Verifica que el nivel en el JSON sea "INFO"
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("Se esperaba el nivel \"level\":\"INFO\" en la salida, pero no se encontró. Salida: %s", output)
	}

	// Verifica que el mensaje esté presente
	if !strings.Contains(output, `"message":"`+message+`"`) {
		t.Errorf("Se esperaba el mensaje \"%s\" en la salida, pero no se encontró. Salida: %s", message, output)
	}
}

func TestLoggerAdapterLogWarn(t *testing.T) {
	// Establece la variable de entorno LOG_FORMAT para activar el formato JSON
	t.Setenv("LOG_FORMAT", "JSON")

	// Ejecuta la función LogWarn usando el adaptador
	message := "Este es un mensaje de advertencia desde el adaptador"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		Logger.LogWarn(message, fileName)
	})

	// Verifica que el nivel en el JSON sea "WARNING"
	if !strings.Contains(output, `"level":"WARNING"`) {
		t.Errorf("Se esperaba el nivel \"level\":\"WARNING\" en la salida, pero no se encontró. Salida: %s", output)
	}

	// Verifica que el mensaje esté presente
	if !strings.Contains(output, `"message":"`+message+`"`) {
		t.Errorf("Se esperaba el mensaje \"%s\" en la salida, pero no se encontró. Salida: %s", message, output)
	}
}

func TestLoggerAdapterLogError(t *testing.T) {
	// Establece la variable de entorno LOG_FORMAT para activar el formato JSON
	t.Setenv("LOG_FORMAT", "JSON")

	// Ejecuta la función LogError usando el adaptador con un error
	message := "Este es un error desde el adaptador"
	fileName := "logger_test.go"
	err := fmt.Errorf("file does not exist")
	output := captureOutput(func() {
		Logger.LogError(message, err, fileName)
	})

	// Verifica que el nivel en el JSON sea "ERROR"
	if !strings.Contains(output, `"level":"ERROR"`) {
		t.Errorf("Se esperaba el nivel \"level\":\"ERROR\" en la salida, pero no se encontró. Salida: %s", output)
	}

	// Verifica que el mensaje de error esté presente
	if !strings.Contains(output, `"message":"`+message+` - Error: file does not exist"`) {
		t.Errorf("Se esperaba el mensaje \"%s - Error: file does not exist\" en la salida, pero no se encontró. Salida: %s", message, output)
	}
}

func TestLoggerAdapterLogDebug(t *testing.T) {
	// Establece la variable de entorno LOG_LEVEL y LOG_FORMAT para activar el formato JSON
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("LOG_FORMAT", "JSON")

	// Ejecuta la función LogDebug usando el adaptador
	message := "Este es un mensaje de depuración desde el adaptador"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		Logger.LogDebug(message, fileName)
	})

	// Verifica que el nivel en el JSON sea "DEBUG"
	if !strings.Contains(output, `"level":"DEBUG"`) {
		t.Errorf("Se esperaba el nivel \"level\":\"DEBUG\" en la salida, pero no se encontró. Salida: %s", output)
	}

	// Verifica que el mensaje esté presente
	if !strings.Contains(output, `"message":"`+message+`"`) {
		t.Errorf("Se esperaba el mensaje \"%s\" en la salida, pero no se encontró. Salida: %s", message, output)
	}
}

func TestLogInfoStringFormat(t *testing.T) {
	t.Setenv("LOG_FORMAT", "STRING")

	message := "Mensaje de información en formato STRING"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		LogInfo(message, fileName)
	})

	// Verifica que el mensaje se registre como STRING
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Se esperaba que el mensaje incluyera [INFO], pero no se encontró. Salida: %s", output)
	}
	if !strings.Contains(output, message) {
		t.Errorf("Se esperaba el mensaje \"%s\", pero no se encontró. Salida: %s", message, output)
	}
	if !strings.Contains(output, fileName) {
		t.Errorf("Se esperaba que el nombre del archivo \"%s\" estuviera en la salida, pero no se encontró. Salida: %s", fileName, output)
	}
}

func TestLogDebugIgnored(t *testing.T) {
	t.Setenv("LOG_LEVEL", "INFO")
	t.Setenv("LOG_FORMAT", "STRING")

	message := "Mensaje de depuración que debe ser ignorado"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		LogDebug(message, fileName)
	})

	// Verifica que no se registre el mensaje
	if output != "" {
		t.Errorf("Se esperaba que el mensaje no se registrara, pero se encontró. Salida: %s", output)
	}
}

func TestUnknownLogFormat(t *testing.T) {
	t.Setenv("LOG_FORMAT", "UNKNOWN")

	message := "Mensaje con formato desconocido"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		LogInfo(message, fileName)
	})

	// Verifica que el logger maneje el caso correctamente (puede ser en formato STRING por defecto)
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Se esperaba que el mensaje se registrara en formato STRING como fallback. Salida: %s", output)
	}
	if !strings.Contains(output, message) {
		t.Errorf("Se esperaba el mensaje \"%s\", pero no se encontró. Salida: %s", message, output)
	}
}

func TestLogErrorWithoutDetails(t *testing.T) {
	t.Setenv("LOG_FORMAT", "JSON")

	message := "Error sin detalles adicionales"
	fileName := "logger_test.go"
	output := captureOutput(func() {
		LogError(message, nil, fileName)
	})

	// Verifica que no incluya el campo de error en el JSON
	if strings.Contains(output, "Error:") {
		t.Errorf("No se esperaba incluir detalles de error, pero se encontraron. Salida: %s", output)
	}
}

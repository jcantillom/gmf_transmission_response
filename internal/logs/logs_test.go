package logs_test

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gmf_transmission_response/internal/logs"
	"testing"
)

func TestLogError(t *testing.T) {
	// Crear un buffer para capturar la salida de los logs
	var logBuffer bytes.Buffer

	// Configurar logrus para escribir en el buffer en lugar de stdout
	logs.Log.SetOutput(&logBuffer)
	logs.Log.SetLevel(logrus.ErrorLevel)

	// Llamar a la función LogError
	logs.LogError(nil, "Este es un error de prueba con valor %d", 42)

	// Verificar que el log generado contiene el mensaje esperado
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Este es un error de prueba con valor 42")
	assert.Contains(t, logOutput, "level=error")
}

func TestLogInfo(t *testing.T) {
	// Crear un buffer para capturar la salida de los logs
	var logBuffer bytes.Buffer

	// Configurar logrus para escribir en el buffer en lugar de stdout
	logs.Log.SetOutput(&logBuffer)
	logs.Log.SetLevel(logrus.InfoLevel)

	// Llamar a la función LogInfo
	logs.LogInfo(context.Background(), "Este es un mensaje informativo con valor %d", 100)

	// Verificar que el log generado contiene el mensaje esperado
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Este es un mensaje informativo con valor 100")
	assert.Contains(t, logOutput, "level=info")
}

func TestLogMigracionTablaCompletada(t *testing.T) {
	// Crear un buffer para capturar la salida de los logs
	var logBuffer bytes.Buffer

	// Configurar logrus para escribir en el buffer en lugar de stdout
	logs.Log.SetOutput(&logBuffer)
	logs.Log.SetLevel(logrus.InfoLevel)

	// Llamar a la función LogMigracionTablaCompletada
	logs.LogMigracionTablaCompletada("usuarios")

	// Verificar que el log generado contiene el mensaje esperado
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Migración de la tabla usuarios completada correctamente 🐘")
	assert.Contains(t, logOutput, "level=info")
}

func TestLogErrorMigrandoTabla(t *testing.T) {
	// Crear un buffer para capturar la salida de los logs
	var logBuffer bytes.Buffer

	// Configurar logrus para escribir en el buffer en lugar de stdout
	logs.Log.SetOutput(&logBuffer)
	logs.Log.SetLevel(logrus.ErrorLevel)

	// Llamar a la función LogErrorMigrandoTabla
	logs.LogErrorMigrandoTabla("usuarios", assert.AnError)

	// Verificar que el log generado contiene el mensaje esperado
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Error al migrar la tabla usuarios: assert.AnError general error for testing")
	assert.Contains(t, logOutput, "level=error")
}

func TestLogWarn(t *testing.T) {
	// Crear un buffer para capturar la salida de los logs
	var logBuffer bytes.Buffer

	// Configurar logrus para escribir en el buffer en lugar de stdout
	logs.Log.SetOutput(&logBuffer)
	logs.Log.SetLevel(logrus.WarnLevel)

	// Llamar a la función LogWarn
	logs.LogWarn(context.Background(), "Este es un mensaje de advertencia con valor %d", 123)

	// Verificar que el log generado contiene el mensaje esperado
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Este es un mensaje de advertencia con valor 123")
	assert.Contains(t, logOutput, "level=warn")
}

func TestLogConexionBaseDatosEstablecida(t *testing.T) {
	// Crear un buffer para capturar la salida de los logs
	var logBuffer bytes.Buffer

	// Configurar logrus para escribir en el buffer en lugar de stdout
	logs.Log.SetOutput(&logBuffer)
	logs.Log.SetLevel(logrus.InfoLevel)

	// Llamar a la función LogConexionBaseDatosEstablecida
	logs.LogConexionBaseDatosEstablecida()

	// Verificar que el log generado contiene el mensaje esperado
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Conexión a la base de datos establecida correctamente 🐘")
	assert.Contains(t, logOutput, "level=info")
}

func TestLogErrorConexionBaseDatos(t *testing.T) {
	// Crear un buffer para capturar la salida de los logs
	var logBuffer bytes.Buffer

	// Configurar logrus para escribir en el buffer en lugar de stdout
	logs.Log.SetOutput(&logBuffer)
	logs.Log.SetLevel(logrus.ErrorLevel)

	// Llamar a la función LogErrorConexionBaseDatos
	logs.LogErrorConexionBaseDatos(assert.AnError)

	// Verificar que el log generado contiene el mensaje esperado
	logOutput := logBuffer.String()
	assert.Contains(
		t,
		logOutput,
		"Error al establecer la conexión a la base de datos: assert.AnError general error for testing",
	)
	assert.Contains(t, logOutput, "level=error")
}

func TestLogErrorCerrandoConexionBaseDatos(t *testing.T) {
	// Crear un buffer para capturar la salida de los logs
	var logBuffer bytes.Buffer

	// Configurar logrus para escribir en el buffer en lugar de stdout
	logs.Log.SetOutput(&logBuffer)
	logs.Log.SetLevel(logrus.ErrorLevel)

	// Llamar a la función LogErrorCerrandoConexionBaseDatos
	logs.LogErrorCerrandoConexionBaseDatos(assert.AnError)

	// Verificar que el log generado contiene el mensaje esperado
	logOutput := logBuffer.String()
	assert.Contains(
		t,
		logOutput,
		"Error al cerrar la conexión a la base de datos: assert.AnError general error for testing",
	)
	assert.Contains(t, logOutput, "level=error")
}

func TestLogConexionBaseDatosCerrada(t *testing.T) {
	// Crear un buffer para capturar la salida de los logs
	var logBuffer bytes.Buffer

	// Configurar logrus para escribir en el buffer en lugar de stdout
	logs.Log.SetOutput(&logBuffer)
	logs.Log.SetLevel(logrus.InfoLevel)

	// Llamar a la función LogConexionBaseDatosCerrada
	logs.LogConexionBaseDatosCerrada()

	// Verificar que el log generado contiene el mensaje esperado
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Conexión a la base de datos cerrada correctamente 🚪")
	assert.Contains(t, logOutput, "level=info")
}

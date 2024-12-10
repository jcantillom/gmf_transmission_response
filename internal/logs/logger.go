package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogInterface define una interfaz para el logger.
type LogInterface interface {
	LogError(message string, err error, fileName string)
	LogInfo(message string, fileName string)
	LogWarn(message string, fileName string, extraArgs ...string)
	LogDebug(message string, fileName string)
}

// LoggerAdapter implementa LogInterface utilizando funciones específicas.
type LoggerAdapter struct{}

func (l *LoggerAdapter) LogError(message string, err error, fileName string) {
	LogError(message, err, fileName)
}

func (l *LoggerAdapter) LogInfo(message, fileName string) {
	LogInfo(message, fileName)
}

func (l *LoggerAdapter) LogWarn(message string, fileName string, extraArgs ...string) {
	LogWarn(message, fileName, extraArgs...)
}

func (l *LoggerAdapter) LogDebug(message string, fileName string) {
	LogDebug(message, fileName)
}

var Logger LogInterface = &LoggerAdapter{}

// BasePath almacena la ruta base del proyecto.
var basePath string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		basePath = ""
	} else {
		basePath = filepath.ToSlash(wd)
	}
}

// getCurrentTimestamp obtiene la fecha y hora actual en la zona horaria de Colombia.
func getCurrentTimestamp() string {
	location, err := time.LoadLocation("America/Bogota")
	if err != nil {
		location = time.UTC
	}
	return time.Now().In(location).Format("2006-01-02 15:04:05")
}

// runtimeCaller es una variable que facilita las pruebas.
var runtimeCaller = runtime.Caller

// getCallerInfo obtiene el archivo y la línea desde donde se llamó el logger.
func getCallerInfo() string {
	// Incrementamos el nivel a 4 para capturar la llamada desde el código que invoca el log.
	_, file, line, ok := runtimeCaller(4)
	if !ok {
		return "???"
	}
	relativePath := strings.TrimPrefix(filepath.ToSlash(file), basePath)
	return fmt.Sprintf("%s:%d", relativePath, line)
}

// logMessage maneja el formato del log sin colores.
func logMessage(level string, message, fileName string) {
	timestamp := getCurrentTimestamp()
	callerInfo := getCallerInfo()
	if fileName != "" {
		fmt.Printf(
			"%s [%s] [%s] [File: %s] %s\n",
			timestamp,
			level,
			callerInfo,
			fileName,
			message,
		)
	} else {
		fmt.Printf("%s [%s] [%s] %s\n", timestamp, level, callerInfo, message)
	}
}

// LogInfo genera logs a nivel INFO.
func LogInfo(message, fileName string) {
	logMessage("INFO", message, fileName)
}

// LogWarn genera logs a nivel WARNING.
func LogWarn(message string, fileName string, extraArgs ...string) {
	if len(extraArgs) >= 2 {
		key := extraArgs[0]
		value := extraArgs[1]
		logMessage("WARNING", fmt.Sprintf("%s - %s: %s", message, key, value), fileName)
	} else {
		logMessage("WARNING", message, fileName)
	}
}

// LogError genera logs a nivel ERROR.
func LogError(message string, err error, fileName string) {
	if err != nil {
		logMessage("ERROR", fmt.Sprintf("%s - Error: %v", message, err), fileName)
	} else {
		logMessage("ERROR", message, fileName)
	}
}

// LogDebug genera logs a nivel DEBUG.
func LogDebug(message, fileName string) {
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel == "DEBUG" {
		logMessage("DEBUG", message, fileName)
	}
}

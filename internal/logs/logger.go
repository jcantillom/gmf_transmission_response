package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogInterface define una interfaz para el logger.
type LogInterface interface {
	LogError(message string, err error, fileName string)
	LogInfo(message, fileName string)
	LogWarn(message string, fileName string, extraArgs ...string)
	LogDebug(message, fileName string)
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

func (l *LoggerAdapter) LogDebug(message, fileName string) {
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

// getCurrentTimestamp obtiene la fecha y hora actual en la zona horaria de Colombia con milisegundos.
func getCurrentTimestamp() string {
	location, err := time.LoadLocation("America/Bogota")
	if err != nil {
		location = time.UTC
	}
	return time.Now().In(location).Format("2006-01-02 15:04:05.000")
}

// runtimeCaller es una variable que facilita las pruebas.
var runtimeCaller = runtime.Caller

// getCallerInfo obtiene el archivo y la línea desde donde se llamó el logger.
func getCallerInfo() (string, int) {
	_, file, line, ok := runtimeCaller(4)
	if !ok {
		return "???", 0
	}
	relativePath := filepath.Base(file)
	return relativePath, line
}

// logMessageJSON maneja el formato JSON del log.
func logMessageJSON(level, message, fileName string) string {
	timestamp := getCurrentTimestamp()
	moduleName, lineNumber := getCallerInfo()

	logData := struct {
		Timestamp  string `json:"timestamp"`
		Level      string `json:"level"`
		ModuleName string `json:"module_name"`
		LineNumber int    `json:"line_number"`
		RequestID  string `json:"request_id"`
		Message    string `json:"message"`
	}{
		Timestamp:  timestamp,
		Level:      level,
		ModuleName: moduleName,
		LineNumber: lineNumber,
		RequestID:  fileName,
		Message:    message,
	}

	jsonLog, err := json.Marshal(logData)
	if err != nil {
		return `{"error": "failed to marshal log to JSON"}`
	}
	return string(jsonLog)
}

// logMessagePlain maneja el formato detallado en texto plano.
func logMessagePlain(level, message, fileName string) string {
	timestamp := getCurrentTimestamp()
	moduleName, lineNumber := getCallerInfo()
	return fmt.Sprintf("%s [%s] [%s:%d] [FileName: %s] %s",
		timestamp, level, moduleName, lineNumber, fileName, message)
}

// LogInfo genera logs a nivel INFO.
func LogInfo(message, fileName string) {
	format := os.Getenv("LOG_FORMAT")
	if format == "JSON" {
		fmt.Println(logMessageJSON("INFO", message, fileName))
	} else {
		fmt.Println(logMessagePlain("INFO", message, fileName))
	}
}

// LogWarn genera logs a nivel WARNING.
func LogWarn(message string, fileName string, extraArgs ...string) {
	var formattedMessage string
	if len(extraArgs) >= 2 {
		key := extraArgs[0]
		value := extraArgs[1]
		formattedMessage = fmt.Sprintf("%s - %s: %s", message, key, value)
	} else {
		formattedMessage = message
	}

	format := os.Getenv("LOG_FORMAT")
	if format == "JSON" {
		fmt.Println(logMessageJSON("WARNING", formattedMessage, fileName))
	} else {
		fmt.Println(logMessagePlain("WARNING", formattedMessage, fileName))
	}
}

// LogError genera logs a nivel ERROR.
func LogError(message string, err error, fileName string) {
	var formattedMessage string
	if err != nil {
		formattedMessage = fmt.Sprintf("%s - Error: %v", message, err)
	} else {
		formattedMessage = message
	}

	format := os.Getenv("LOG_FORMAT")
	if format == "JSON" {
		fmt.Println(logMessageJSON("ERROR", formattedMessage, fileName))
	} else {
		fmt.Println(logMessagePlain("ERROR", formattedMessage, fileName))
	}
}

// LogDebug genera logs a nivel DEBUG.
func LogDebug(message, fileName string) {
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel == "DEBUG" {
		format := os.Getenv("LOG_FORMAT")
		if format == "JSON" {
			fmt.Println(logMessageJSON("DEBUG", message, fileName))
		} else {
			fmt.Println(logMessagePlain("DEBUG", message, fileName))
		}
	}
}

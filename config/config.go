package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gmf_transmission_response/internal/logs"
	"log"
)

// Manager maneja la carga de configuración de la aplicación.
type Manager struct{}

// NewConfigManager crea una nueva instancia de Manager.
func NewConfigManager() *Manager {
	return &Manager{}
}

// InitConfig carga la configuración desde el archivo .env y las variables de entorno.
func (cm *Manager) InitConfig() {
	// Cargar el archivo .env si existe
	if err := godotenv.Load(); err != nil {
		logs.Logger.LogWarn("Archivo .env no encontrado, usando solo variables de entorno", "CONFIG_INIT")
	}

	// Configurar Viper para usar el prefijo "APP" y cargar variables automáticamente del entorno
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Especificar el archivo .env y cargarlo con Viper
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Leer las configuraciones del archivo .env
	if err := viper.ReadInConfig(); err != nil {
		logs.Logger.LogWarn("Error al leer el archivo de configuración", "CONFIG_INIT", err.Error())
	}

	// Lista de variables de entorno obligatorias
	requiredEnvVars := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
	}

	// Validar que todas las variables de entorno requeridas estén presentes
	for _, envVar := range requiredEnvVars {
		value := viper.GetString(envVar)
		if value == "" {
			msg := "La variable de entorno " + envVar + " no está configurada"
			logs.Logger.LogError(msg, nil, "CONFIG_INIT")
			log.Fatalf("Error: %s", msg)
		}
	}

	// Enlazar las variables de entorno con Viper
	for _, envVar := range requiredEnvVars {
		viper.BindEnv(envVar)
	}

	logs.Logger.LogInfo("Configuración cargada correctamente", "CONFIG_INIT")
}

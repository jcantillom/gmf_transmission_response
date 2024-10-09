package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gmf_transmission_response/internal/logs"
)

// Manager ConfigManager maneja la carga de configuración de la aplicación.
type Manager struct{}

// NewConfigManager crea una nueva instancia de ConfigManager.
func NewConfigManager() *Manager {
	return &Manager{}
}

// InitConfig carga la configuración desde el archivo .env y las variables de entorno.
func (cm *Manager) InitConfig() {
	// Cargar el archivo .env si existe
	if err := godotenv.Load(); err != nil {
		logs.LogError(nil, "archivo .env no encontrado")
	}

	// Configurar Viper para usar el prefijo "APP" y cargar variables automáticamente del entorno
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Especificar el archivo .env y cargarlo con Viper
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Leer las configuraciones del archivo .env
	if err := viper.ReadInConfig(); err != nil {
		logs.LogError(nil, "error al leer el archivo de configuración: %v", err)
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
		// Verificar si la variable está presente
		value := viper.GetString(envVar)
		if value == "" {
			log.Fatalf(
				"Error: La variable de entorno %s no está configurada en el archivo .env o en el entorno.",
				envVar,
			)
		}
	}

	// Enlazar las variables de entorno
	for _, envVar := range requiredEnvVars {
		viper.BindEnv(envVar)
	}
}

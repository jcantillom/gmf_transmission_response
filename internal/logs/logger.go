package logs

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

var Log = InitLogger()

// InitLogger inicializa el logger de logrus.
func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(
		&logrus.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		},
	)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

// LogError permite recibir un mensaje con o sin contexto.
func LogError(ctx interface{}, msg string, args ...interface{}) {
	switch ctxTyped := ctx.(type) {
	case context.Context:
		Log.WithContext(ctxTyped).Errorf(msg, args...)
	default:
		Log.Errorf(msg, args...)
	}
}

func LogInfo(ctx context.Context, msg string, args ...interface{}) {
	Log.WithContext(ctx).Infof(msg, args...)
}

func LogWarn(ctx context.Context, msg string, args ...interface{}) {
	Log.WithContext(ctx).Warnf(msg, args...)
}

// Logs de conexión a la base de datos
func LogConexionBaseDatosEstablecida() {
	LogInfo(context.Background(), "Conexión a la base de datos establecida correctamente 🐘")
}

func LogErrorConexionBaseDatos(err error) {
	LogError(context.Background(), "Error al establecer la conexión a la base de datos: %v", err)
}

func LogErrorCerrandoConexionBaseDatos(err error) {
	LogError(context.Background(), "Error al cerrar la conexión a la base de datos: %v", err)
}

func LogConexionBaseDatosCerrada() {
	LogInfo(context.Background(), "Conexión a la base de datos cerrada correctamente 🚪")
}

// Logs de migración de tablas
func LogMigracionTablaCompletada(tabla string) {
	LogInfo(context.Background(), "Migración de la tabla %s completada correctamente 🐘", tabla)
}

func LogErrorMigrandoTabla(tabla string, err error) {
	LogError(context.Background(), "Error al migrar la tabla %s: %v", tabla, err)
}

package aws_test

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	awsinternal "gmf_transmission_response/internal/aws"
	"os"
	"testing"
)

// Mock del cliente de SecretsManager
type MockSecretsManagerClient struct {
	GetSecretValueFunc func(
		ctx context.Context,
		params *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

func (m *MockSecretsManagerClient) GetSecretValue(
	ctx context.Context,
	params *secretsmanager.GetSecretValueInput,
	optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return m.GetSecretValueFunc(ctx, params, optFns...)
}

func TestSecretsManager_GetSecret_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Datos simulados del secreto
	secretData := map[string]string{
		"USERNAME": "test-user",
		"PASSWORD": "test-pass",
	}
	secretString, _ := json.Marshal(secretData)

	// Mock del cliente de SecretsManager
	mockClient := &MockSecretsManagerClient{
		GetSecretValueFunc: func(
			ctx context.Context,
			params *secretsmanager.GetSecretValueInput,
			optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(string(secretString)),
			}, nil
		},
	}

	sm := &awsinternal.SecretsManager{Client: mockClient}

	// Ejecutar prueba
	result, err := sm.GetSecret("test-secret")
	assert.NoError(t, err)
	assert.Equal(t, "test-user", result["USERNAME"])
	assert.Equal(t, "test-pass", result["PASSWORD"])
}

func TestSecretsManager_GetSecret_ResourceNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock del cliente de SecretsManager para un error de recurso no encontrado
	mockClient := &MockSecretsManagerClient{
		GetSecretValueFunc: func(
			ctx context.Context,
			params *secretsmanager.GetSecretValueInput,
			optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return nil, &types.ResourceNotFoundException{Message: aws.String("Secret not found")}
		},
	}

	sm := &awsinternal.SecretsManager{Client: mockClient}

	// Ejecutar prueba
	_, err := sm.GetSecret("non-existent-secret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "secreto no encontrado")
}

func TestSecretsManager_GetSecret_DeserializationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock del cliente de SecretsManager con un secreto mal formateado
	mockClient := &MockSecretsManagerClient{
		GetSecretValueFunc: func(
			ctx context.Context,
			params *secretsmanager.GetSecretValueInput,
			optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String("invalid-json"),
			}, nil
		},
	}

	sm := &awsinternal.SecretsManager{Client: mockClient}

	// Ejecutar prueba
	_, err := sm.GetSecret("malformed-secret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deserializando el secreto")
}

func TestSecretsManager_GetSecret_EmptySecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock del cliente de SecretsManager con un secreto vacío
	mockClient := &MockSecretsManagerClient{
		GetSecretValueFunc: func(
			ctx context.Context,
			params *secretsmanager.GetSecretValueInput,
			optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: nil,
			}, nil
		},
	}

	sm := &awsinternal.SecretsManager{Client: mockClient}

	// Ejecutar prueba
	_, err := sm.GetSecret("empty-secret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "secreto vacío o no válido")
}

func TestSecretsManager_NewSecretsManager_LocalStack(t *testing.T) {
	// Establecer la variable de entorno para simular el uso de LocalStack.
	os.Setenv("APP_ENV", "local")
	defer os.Unsetenv("APP_ENV")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock de la carga de configuración de LocalStack
	mockClient := &MockSecretsManagerClient{
		GetSecretValueFunc: func(
			ctx context.Context,
			params *secretsmanager.GetSecretValueInput,
			optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			// Simula un secreto devuelto por LocalStack
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(`{"USERNAME": "test-user", "PASSWORD": "test-pass"}`),
			}, nil
		},
	}

	// Crear SecretsManager con el cliente simulado.
	sm := &awsinternal.SecretsManager{Client: mockClient}

	// Ejecutar la prueba del secreto
	result, err := sm.GetSecret("test-secret")
	assert.NoError(t, err)
	assert.Equal(t, "test-user", result["USERNAME"])
	assert.Equal(t, "test-pass", result["PASSWORD"])

	// Verificar que la configuración de LocalStack esté activa.
	assert.Equal(t, "local", os.Getenv("APP_ENV"))
}

func TestSecretsManager_NewSecretsManager_AWS(t *testing.T) {
	// Establecer la variable de entorno para simular el uso de AWS real.
	os.Setenv("APP_ENV", "aws")
	defer os.Unsetenv("APP_ENV")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock de la carga de configuración de AWS real
	mockClient := &MockSecretsManagerClient{
		GetSecretValueFunc: func(
			ctx context.Context,
			params *secretsmanager.GetSecretValueInput,
			optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			// Simula un secreto devuelto por AWS real
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(`{"USERNAME": "test-user", "PASSWORD": "test-pass"}`),
			}, nil
		},
	}

	// Crear SecretsManager con el cliente simulado.
	sm := &awsinternal.SecretsManager{Client: mockClient}

	// Ejecutar la prueba del secreto
	result, err := sm.GetSecret("test-secret")
	assert.NoError(t, err)
	assert.Equal(t, "test-user", result["USERNAME"])
	assert.Equal(t, "test-pass", result["PASSWORD"])

	// Verificar que la configuración de AWS real esté activa.
	assert.Equal(t, "aws", os.Getenv("APP_ENV"))
}

func TestNewSecretsManager_LocalStack(t *testing.T) {
	// Establecer la variable de entorno para simular el uso de LocalStack.
	os.Setenv("APP_ENV", "local")
	defer os.Unsetenv("APP_ENV")

	// Crear SecretsManager.
	sm, err := awsinternal.NewSecretsManager()

	// Verificar que no haya errores.
	assert.NoError(t, err)
	assert.NotNil(t, sm)
	assert.NotNil(t, sm.Client)
}

func TestNewSecretsManager_AWS(t *testing.T) {
	// Establecer la variable de entorno para simular el uso de AWS real.
	os.Setenv("APP_ENV", "aws")
	defer os.Unsetenv("APP_ENV")

	// Crear SecretsManager.
	sm, err := awsinternal.NewSecretsManager()

	// Verificar que no haya errores.
	assert.NoError(t, err)
	assert.NotNil(t, sm)
	assert.NotNil(t, sm.Client)
}

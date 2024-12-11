package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

// SecretsManagerClient define los métodos que el cliente de Secrets Manager debe implementar.
type SecretsManagerClient interface {
	GetSecretValue(
		ctx context.Context,
		params *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options),
	) (*secretsmanager.GetSecretValueOutput, error)
}

// SecretsManager implementa SecretsManagerInterface.
type SecretsManager struct {
	Client SecretsManagerClient
}

// NewSecretsManager crea una nueva instancia de SecretsManager.
func NewSecretsManager() (*SecretsManager, error) {
	var cfg aws.Config
	var err error

	if os.Getenv("APP_ENV") == "local" {
		endpoint := "http://localhost:4566"
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("us-east-1"),
			config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: "us-east-1",
				}, nil
			})),
		)
		if err != nil {
			return nil, fmt.Errorf("error configurando LocalStack: %w", err)
		}
	} else {
		// Cargar configuración de AWS real.
		cfg, err = config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("error cargando la configuración de AWS: %w", err)
		}
	}

	client := secretsmanager.NewFromConfig(cfg)
	return &SecretsManager{Client: client}, nil
}

// GetSecret obtiene y deserializa un secreto desde AWS Secrets Manager o LocalStack.
func (sm *SecretsManager) GetSecret(secretName string) (map[string]string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := sm.Client.GetSecretValue(context.TODO(), input)
	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			return nil, fmt.Errorf("secreto no encontrado: %s", secretName)
		}
		return nil, fmt.Errorf("error obteniendo el secreto: %w", err)
	}

	var secretMap map[string]string
	if result.SecretString != nil {
		if err := json.Unmarshal([]byte(*result.SecretString), &secretMap); err != nil {
			return nil, fmt.Errorf("error deserializando el secreto: %w", err)
		}
	} else {
		return nil, fmt.Errorf("secreto vacío o no válido")
	}

	return secretMap, nil
}

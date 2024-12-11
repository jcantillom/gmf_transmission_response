package connection

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockAWSSecretsManager es un mock de AWS Secrets Manager.
type MockAWSSecretsManager struct {
	mock.Mock
}

func (m *MockAWSSecretsManager) GetSecret(secretName string) (map[string]string, error) {
	args := m.Called(secretName)
	return args.Get(0).(map[string]string), args.Error(1)
}

func TestNewDBManager(t *testing.T) {
	dbManager := NewDBManager()
	assert.NotNil(t, dbManager)
	assert.Nil(t, dbManager.DB)
}

func TestDBManager_InitDB_Error(t *testing.T) {
	mockSecretsManager := new(MockAWSSecretsManager)
	mockSecretsManager.On(
		"GetSecret",
		"test_secret").Return(
		nil, errors.New("error getting secret"))

	os.Setenv("SECRETS_DB", "test_secret")

	dbManager := &DBManager{}
	err := dbManager.InitDB()
	assert.Error(t, err)
	assert.Nil(t, dbManager.DB)
}

func TestDBManager_GetDB(t *testing.T) {
	dbManager := &DBManager{
		DB: &gorm.DB{},
	}
	db := dbManager.GetDB()
	assert.NotNil(t, db)
}

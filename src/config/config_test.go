package config_test

import (
	"testing"
	"time"

	"github.com/kodmain/kitsune/src/config"
	"github.com/stretchr/testify/assert"
)

func TestDefaultValues(t *testing.T) {
	assert.Equal(t, "kitsune", config.DEFAULT_APP_NAME, "La valeur par défaut de DEFAULT_APP_NAME devrait être 'kitsune'")
	assert.Equal(t, "local", config.DEFAULT_VERSION, "La valeur par défaut de DEFAULT_VERSION devrait être 'local'")
	assert.Equal(t, "local", config.BUILD_VERSION, "La valeur de BUILD_VERSION devrait être initialisée à 'local'")
	assert.Equal(t, "kitsune", config.BUILD_APP_NAME, "La valeur de BUILD_APP_NAME devrait être initialisée à 'kitsune'")
}
func TestRetryDefaultValues(t *testing.T) {
	assert.Equal(t, time.Duration(1), config.DEFAULT_RETRY_INTERVAL, "La valeur par défaut de DEFAULT_RETRY_INTERVAL devrait être 1")
	assert.Equal(t, time.Duration(3), config.DEFAULT_RETRY_MAX, "La valeur par défaut de DEFAULT_RETRY_MAX devrait être 3")
	assert.Equal(t, time.Duration(15), config.DEFAULT_TIMEOUT, "La valeur par défaut de DEFAULT_TIMEOUT devrait être 15")
	assert.Equal(t, time.Duration(15), config.DEFAULT_CACHE, "La valeur par défaut de DEFAULT_CACHE devrait être 15")
}

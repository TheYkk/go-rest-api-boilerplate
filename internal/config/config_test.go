package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {

	t.Run("Config not found", func(t *testing.T) {

		_, err := GetAllValues(".", "NoFile")
		assert.Contains(t, err.Error(), "Not Found")
	})
	t.Run("Invalid Config file", func(t *testing.T) {

		_, err := GetAllValues("./test-fixtures", "invalid-test-config")
		assert.Contains(t, err.Error(), "While parsing config: yaml")
	})
	t.Run("Config found", func(t *testing.T) {

		cfg, err := GetAllValues("./test-fixtures", "test-config")
		assert.NotNil(t, cfg)
		assert.Nil(t, err)
	})

}

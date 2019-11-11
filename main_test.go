package main

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigWithDefaultPath(t *testing.T) {
	loadConfig()
	actualHost := viper.GetString("server.host")
	assert.Equal(t, "localhost", actualHost)
}

package config_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	t.Logf("Current test filename: %s", filename)
	path := filepath.Join(filename, "./../../../../examples/reposcrtl.yaml")

	testConfigPath := []string{path}
	cfg, err := config.NewConfigFromFiles(testConfigPath)
	assert.Nil(t, err, "No Error expected %s", err)
	assert.NotNil(t, cfg, "Expected pared Config")
}

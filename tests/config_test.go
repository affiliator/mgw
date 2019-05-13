package tests

import (
	"github.com/affiliator/mgw/config"
	"testing"
)

func TestInitialize(t *testing.T) {
	result := config.Initialize()

	if result == nil {
		t.Error("Config not initialized.")
	}
}

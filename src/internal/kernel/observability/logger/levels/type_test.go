package levels_test

import (
	"testing"

	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/stretchr/testify/assert"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    levels.TYPE
		expected string
	}{
		{levels.OFF, "OFF"},
		{levels.PANIC, "PANIC"},
		{levels.FATAL, "FATAL"},
		{levels.ERROR, "ERROR"},
		{levels.SUCCESS, "SUCCESS"},
		{levels.MESSAGE, "MESSAGE"},
		{levels.WARN, "WARN"},
		{levels.INFO, "INFO"},
		{levels.DEBUG, "DEBUG"},
		{levels.TRACE, "TRACE"},
		{levels.TYPE(99), "UNKNOWN"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.level.String(), "La chaîne de caractères devrait correspondre au niveau")
	}
}

func TestLevelColor(t *testing.T) {
	tests := []struct {
		level    levels.TYPE
		expected string
	}{
		{levels.PANIC, "9"},
		{levels.FATAL, "160"},
		{levels.ERROR, "1"},
		{levels.SUCCESS, "2"},
		{levels.MESSAGE, "7"},
		{levels.WARN, "3"},
		{levels.INFO, "4"},
		{levels.DEBUG, "6"},
		{levels.TRACE, "7"},
		{levels.TYPE(99), "UNKNOWN"}, // Test pour un niveau inconnu
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.level.Color(), "La couleur devrait correspondre au niveau")
	}
}

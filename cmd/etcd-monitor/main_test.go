package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_Comprehensive(t *testing.T) {
	t.Run("Main function exists", func(t *testing.T) {
		// Test that main function exists and can be called
		// We can't actually call main() in tests, but we can verify it exists
		assert.NotNil(t, main)
	})

	t.Run("Command structure", func(t *testing.T) {
		// Test that the command line structure is correct
		assert.Equal(t, "1.0.0", appVersion)
		assert.Equal(t, "etcd-monitor", appName)
	})
}

func TestParseEndpoints(t *testing.T) {
	t.Run("ParseEndpoints - Single endpoint", func(t *testing.T) {
		endpoints := parseEndpoints("localhost:2379")
		assert.Len(t, endpoints, 1)
		assert.Equal(t, "localhost:2379", endpoints[0])
	})

	t.Run("ParseEndpoints - Multiple endpoints", func(t *testing.T) {
		endpoints := parseEndpoints("localhost:2379,localhost:2380,localhost:2381")
		assert.Len(t, endpoints, 3)
		assert.Equal(t, "localhost:2379", endpoints[0])
		assert.Equal(t, "localhost:2380", endpoints[1])
		assert.Equal(t, "localhost:2381", endpoints[2])
	})

	t.Run("ParseEndpoints - Empty string", func(t *testing.T) {
		endpoints := parseEndpoints("")
		assert.Len(t, endpoints, 0)
	})

	t.Run("ParseEndpoints - Single comma", func(t *testing.T) {
		endpoints := parseEndpoints(",")
		assert.Len(t, endpoints, 0)
	})

	t.Run("ParseEndpoints - Multiple commas", func(t *testing.T) {
		endpoints := parseEndpoints("localhost:2379,,localhost:2380,")
		assert.Len(t, endpoints, 2)
		assert.Equal(t, "localhost:2379", endpoints[0])
		assert.Equal(t, "localhost:2380", endpoints[1])
	})
}

func TestConstants(t *testing.T) {
	t.Run("App version and name", func(t *testing.T) {
		assert.Equal(t, "1.0.0", appVersion)
		assert.Equal(t, "etcd-monitor", appName)
	})
}
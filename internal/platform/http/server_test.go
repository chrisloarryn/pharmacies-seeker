package platformhttp

import (
	"errors"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServerStoresApp(t *testing.T) {
	app := fiber.New()

	server := NewServer("8080", app)

	assert.Same(t, app, server.App())
}

func TestRunDelegatesToListenFunction(t *testing.T) {
	app := fiber.New()
	server := NewServer("9090", app)

	var addr string
	server.listen = func(gotAddr string, _ ...fiber.ListenConfig) error {
		addr = gotAddr
		return errors.New("boom")
	}

	err := server.Run()

	require.EqualError(t, err, "boom")
	assert.Equal(t, ":9090", addr)
}

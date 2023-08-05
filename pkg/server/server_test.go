//go:build unit

package server

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	configmock "go-rest-api-boilerplate/pkg/config/mock"
)

func TestServer(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	t.Run("should create server instance and return server instance", func(t *testing.T) {
		mockConfig := configmock.NewMockConfig(mockController)
		mockConfig.
			EXPECT().
			GetServerPort().
			Return(
				"8080",
			).
			Times(1)

		var handlers []Handler
		testServer := NewServer(mockConfig, handlers)

		assert.IsType(t, &server{}, testServer)
	})

	t.Run("should server start and stop", func(t *testing.T) {
		mockConfig := configmock.NewMockConfig(mockController)
		mockConfig.
			EXPECT().
			GetServerPort().
			Return(
				"8080",
			).
			Times(1)

		var handlers []Handler
		testServer := NewServer(mockConfig, handlers)

		go func() {
			err := testServer.Start()
			assert.NoError(t, err)
		}()

		err := testServer.Shutdown()
		assert.NoError(t, err)
	})
}

func TestServer_GetFiberInstance(t *testing.T) {
	testServer := &server{
		fiber: fiber.New(),
	}
	fiberInstance := testServer.GetFiberInstance()

	assert.IsType(t, fiberInstance, testServer.fiber)
}

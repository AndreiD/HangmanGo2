package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestAuthInterceptor(t *testing.T) {

	md := metadata.New(map[string]string{"login": "andy", "password": "2k3jfeij324fds2!2q433@#$3"})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	output, err := authenticateClient(ctx)

	t.Run("Test Auth Positive Test", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, output, "1337", "test user should be with id 1337")
	})

}

func TestInvalidCredentials(t *testing.T) {

	md := metadata.New(map[string]string{"login": "andy123", "password": "2k3jfeij324fds2!2q433@#$3"})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := authenticateClient(ctx)

	t.Run("Test Auth Unknown User", func(t *testing.T) {
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "unknown user")
		}
	})

	md = metadata.New(map[string]string{"login": "andy", "password": "1231231"})
	ctx = metadata.NewIncomingContext(context.Background(), md)

	_, err = authenticateClient(ctx)

	t.Run("Test Auth Bad Password", func(t *testing.T) {
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "bad password")
		}
	})

}

func TestMissingCredentials(t *testing.T) {

	ctx := context.Background()

	_, err := authenticateClient(ctx)

	t.Run("Test Missing Credentials", func(t *testing.T) {
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "missing credentials")
		}
	})

}

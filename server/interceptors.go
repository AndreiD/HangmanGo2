package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// private type for Context keys
type contextKey int

const (
	clientIDKey contextKey = iota
)

// authenticateAgent check the client credentials
// private, used by interceptor
func authenticateClient(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")

		// user, pass and clientID should sit in a database
		if clientLogin != "andy" {
			return "", fmt.Errorf("unknown user %s", clientLogin)
		}
		if clientPassword != "2k3jfeij324fds2!2q433@#$3" {
			return "", fmt.Errorf("bad password %s", clientPassword)
		}
		log.Printf("client: %s requested action", clientLogin)
		return "1337", nil
	}
	return "", fmt.Errorf("missing credentials")
}

// UnaryInterceptor calls authenticateClient with current context
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	clientID, err := authenticateClient(ctx)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)
}

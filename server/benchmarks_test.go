package main

import (
	"context"
	"testing"

	"github.com/AndreiD/HangmanGo2/api"
)

// silly benchmark
func BenchmarkStartServers(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hangm := &hangman{}
		hangm.NewGame(context.Background(), &api.GameRequest{Id: -1})
		hangm.NewGame(context.Background(), &api.GameRequest{RetryLimit: 5})

	}

}

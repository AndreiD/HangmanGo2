package cmd

import (
	"github.com/AndreiD/HangmanGo2/api"
)

// SaveGame saves the current game
func SaveGame(client api.HangmanClient, g *api.Game) error {
	ctx, cancel := AppContext()
	defer cancel()
	_, err := client.SaveGame(ctx, &api.GameRequest{Id: g.Id})
	if err != nil {
		return err
	}
	return nil
}

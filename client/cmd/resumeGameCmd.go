package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/AndreiD/HangmanGo2/api"
	"github.com/spf13/cobra"
)

//used as a flag
var gameid string

// resumeGameCmd represents the resumeGame command
var resumeGameCmd = &cobra.Command{
	Use:   "resumeGame",
	Short: "Resums a gamewith id.",
	Long:  `resums a game with id.`,
	Run: func(cmd *cobra.Command, args []string) {

		// check if the gameid was passed as a flag
		if cmd.Flags().Changed("game") {

			game, err := resumeGame(hangmanClient, gameid)
			if err != nil {
				fmt.Printf("could not restart the game: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("Game Restarted.  Hint it's a %d letters word \n", len(game.WordMasked))

			OnGoingGame(game)

		} else {
			cmd.Usage()
		}

	},
}

func init() {
	// add the gameid flag
	resumeGameCmd.PersistentFlags().StringVar(&gameid, "game", "", "use --game=X to access the game with id X")
	rootCmd.AddCommand(resumeGameCmd)
}

// resumes a game with id
func resumeGame(client api.HangmanClient, id string) (*api.Game, error) {
	_gameID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("Game id should be a number")
	}

	ctx, cancel := AppContext()
	defer cancel()
	game, err := client.ResumeGame(ctx, &api.GameRequest{Id: int32(_gameID)})
	if err != nil {
		return nil, err
	}
	return game, nil
}

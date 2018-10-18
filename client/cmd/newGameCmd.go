package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AndreiD/HangmanGo2/api"
	"github.com/spf13/cobra"
)

var newGameCmd = &cobra.Command{
	Use:   "newGame",
	Short: "Starts a new game",
	Long:  `Starts a new game of hangman`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("newGame called")

		game, err := newGame(hangmanClient)
		if err != nil {
			log.Printf("could not start a new game %s", err)
			os.Exit(1)
		}

		fmt.Printf("Game started.  Hint it's a %d letters word \n", len(game.WordMasked))

		OnGoingGame(game)

	},
}

func init() {
	rootCmd.AddCommand(newGameCmd)

}

func newGame(client api.HangmanClient) (*api.Game, error) {
	ctx, cancel := AppContext()
	defer cancel()
	r, err := client.NewGame(ctx, &api.GameRequest{RetryLimit: int32(Config.GetInt("game_difficulty"))})
	if err != nil {
		return nil, err
	}
	fmt.Println("Client ready!")
	fmt.Printf("Client difficulty configured to %d retries\n", Config.GetInt("game_difficulty"))
	return r, nil
}

//checks if you lost
func checkIfYouLost(err error) {
	if fmt.Sprintf("%s", err) == "You Lost" {
		fmt.Println("▓██   ██▓ ▒█████   █    ██    ▓█████▄  ██▓▓█████ ▓█████▄ ")
		fmt.Println(" ▒██  ██▒▒██▒  ██▒ ██  ▓██▒   ▒██▀ ██▌▓██▒▓█   ▀ ▒██▀ ██▌")
		fmt.Println("  ▒██ ██░▒██░  ██▒▓██  ▒██░   ░██   █▌▒██▒▒███   ░██   █▌")
		fmt.Println("  ░ ▐██▓░▒██   ██░▓▓█  ░██░   ░▓█▄   ▌░██░▒▓█  ▄ ░▓█▄   ▌")
		fmt.Println("  ░ ██▒▓░░ ████▓▒░▒▒█████▓    ░▒████▓ ░██░░▒████▒░▒████▓ ")
		fmt.Println("   ██▒▒▒ ░ ▒░▒░▒░ ░▒▓▒ ▒ ▒     ▒▒▓  ▒ ░▓  ░░ ▒░ ░ ▒▒▓  ▒ ")
		fmt.Println(" ▓██ ░▒░   ░ ▒ ▒░ ░░▒░ ░ ░     ░ ▒  ▒  ▒ ░ ░ ░  ░ ░ ▒  ▒ ")
		fmt.Println(" ▒ ▒ ░░  ░ ░ ░ ▒   ░░░ ░ ░     ░ ░  ░  ▒ ░   ░    ░ ░  ░ ")
		fmt.Println(" ░ ░         ░ ░     ░           ░     ░     ░  ░   ░    ")
		fmt.Println(" ░ ░                           ░                  ░      ")
		os.Exit(0)
	}
}

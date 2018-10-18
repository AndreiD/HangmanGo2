package cmd

import (
	"errors"
	"fmt"

	"github.com/AndreiD/HangmanGo2/api"
	"github.com/spf13/cobra"
)

var listGamesCmd = &cobra.Command{
	Use:   "listGames",
	Short: "Shows your saved games",
	Long:  `Shows your saved games`,
	Run: func(cmd *cobra.Command, args []string) {

		err := listGames(hangmanClient)
		if err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(listGamesCmd)
}

// prints the games on the screen
func listGames(client api.HangmanClient) error {
	var output string

	ctx, cancel := AppContext()
	defer cancel()

	r, err := client.ListGames(ctx, &api.GameRequest{Id: -1})
	if err != nil {
		return err
	}

	if len(r.Game) > 0 {
		output += "ID	Ongoing		Attempts Left	Hint \n"
		for _, v := range r.Game {
			status := "          "
			if v.Status {
				status = "in progress"
			}
			output += fmt.Sprint(v.Id, "	", status, "	", v.RetryLeft, "		", v.WordMasked, "\n")
		}
	} else {
		return errors.New("No saved games on the server")
	}

	fmt.Println("███████╗ █████╗ ██╗   ██╗███████╗██████╗      ██████╗  █████╗ ███╗   ███╗███████╗███████╗")
	fmt.Println("██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗    ██╔════╝ ██╔══██╗████╗ ████║██╔════╝██╔════╝")
	fmt.Println("███████╗███████║██║   ██║█████╗  ██║  ██║    ██║  ███╗███████║██╔████╔██║█████╗  ███████╗")
	fmt.Println("╚════██║██╔══██║╚██╗ ██╔╝██╔══╝  ██║  ██║    ██║   ██║██╔══██║██║╚██╔╝██║██╔══╝  ╚════██║")
	fmt.Println("███████║██║  ██║ ╚████╔╝ ███████╗██████╔╝    ╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗███████║")
	fmt.Println("╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═════╝      ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝")
	fmt.Println("                                                                                         ")

	fmt.Println(output)
	return nil
}

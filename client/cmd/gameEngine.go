package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AndreiD/HangmanGo2/api"
	"github.com/AndreiD/HangmanGo2/client/utils"
)

// OnGoingGame this waiting for you to win or lose it
func OnGoingGame(game *api.Game) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter a letter: ")
		letter, _ := reader.ReadString('\n')
		letter = strings.TrimSpace(letter)
		fmt.Println(">> you entered: " + letter)

		if utils.IsLetter(letter) {
			output, err := GuessLetter(hangmanClient, game, letter)
			if err != nil {
				// should move them away from errors!
				game.Status = "lost"
				checkIfYouLost(err)
				game.Status = "won"
				checkIfYouWin(err)
				game.Status = "ongoing"
				fmt.Println(err)
			}
			// save the game on errors
			if er := SaveGame(hangmanClient, game); er != nil {
				fmt.Println(err)
			}

			fmt.Println(output)
			continue
		} else {
			fmt.Println("You did not entered a letter.")
		}
	}

}

// GuessLetter is the main logic of the game.
func GuessLetter(client api.HangmanClient, g *api.Game, l string) (string, error) {

	var reply string

	if len(l) != 1 {
		return "Please provide a single letter", nil
	}
	ctx, cancel := AppContext()
	defer cancel()

	if g.Id > 0 {
		gg, err := client.GuessLetter(ctx, &api.GuessRequest{GameID: g.Id, Letter: l})
		if err != nil {
			return "", err
		}

		//refactor this to a normal message, not error
		if gg.RetryLeft < 1 {
			fmt.Printf("\nYou failed to guess: %s\n\n", g.Word)
			return "", errors.New("you lost")
		}

		if strings.Index(gg.WordMasked, "_") == -1 {
			return "", fmt.Errorf("you won")
		}

		reply += hangingArt[(len(hangingArt) - int(gg.RetryLeft) - 1)]
		reply += fmt.Sprintf("\nRemaining attempts: %v", gg.RetryLeft)
		reply += ("\nIncorrect attempts: ")
		for _, v := range gg.IncorrectGuesses {
			reply += fmt.Sprint(v.Letter, " ")
		}
		reply += fmt.Sprint("\nWord hint:", gg.WordMasked)
	} else {
		return "", errors.New("Invalid Game ID")
	}
	return reply, nil
}

var hangingArt = []string{
	`
    _________
    |/      |
    |
    |
    |
    |
    |
____|____`,
	`
    _________
    |/      |
    |      (_)
    |
    |
    |
    |
____|____`,
	`
    __________
    |/      |
    |      (_)
    |       |
    |       |
    |
    |
____|____`,
	`
    __________
    |/      |
    |      (_)
    |      \|
    |       |
    |
    |
____|____`,
	`
    __________
    |/      |
    |      (_)
    |      \|/
    |       |
    |
    |
____|____`,
	`
    __________
    |/      |
    |      (_)
    |      \|/
    |       |
    |      /
    |
____|____`,
	`
    __________
    |/      |
    |      (_)
    |      \|/
    |       |
    |      / \
    |
____|____`,
}

//checks if you lost
func checkIfYouLost(err error) {
	if strings.Contains(fmt.Sprintf("%s", err), "you lost") {
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

//checks if you won
func checkIfYouWin(err error) {
	if strings.Contains(fmt.Sprintf("%s", err), "you won") {
		fmt.Println("██╗   ██╗ ██████╗ ██╗   ██╗    ██╗    ██╗ ██████╗ ███╗   ██╗")
		fmt.Println("╚██╗ ██╔╝██╔═══██╗██║   ██║    ██║    ██║██╔═══██╗████╗  ██║")
		fmt.Println(" ╚████╔╝ ██║   ██║██║   ██║    ██║ █╗ ██║██║   ██║██╔██╗ ██║")
		fmt.Println("  ╚██╔╝  ██║   ██║██║   ██║    ██║███╗██║██║   ██║██║╚██╗██║")
		fmt.Println("   ██║   ╚██████╔╝╚██████╔╝    ╚███╔███╔╝╚██████╔╝██║ ╚████║")
		fmt.Println("   ╚═╝    ╚═════╝  ╚═════╝      ╚══╝╚══╝  ╚═════╝ ╚═╝  ╚═══╝")
		fmt.Println("                                                            ")
		os.Exit(0)
	}
}

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AndreiD/HangmanGo2/api"
)

// OnGoingGame this waiting for you to win or lose it
func OnGoingGame(game *api.Game) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter a letter: ")
		letter, _ := reader.ReadString('\n')
		letter = strings.TrimSpace(letter)

		if isLetter(letter) {
			output, err := GuessALetter(hangmanClient, game, letter)
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

// GuessALetter is the main logic of the game.
func GuessALetter(client api.HangmanClient, game *api.Game, letter string) (string, error) {

	var reply string

	ctx, cancel := AppContext()
	defer cancel()

	if game.Id > 0 {
		nGame, err := client.GuessLetter(ctx, &api.GuessRequest{GameID: game.Id, Letter: letter})
		if err != nil {
			return "", err
		}

		//refactor this to a normal message, not error
		if nGame.RetryLeft < 1 {
			fmt.Printf("\nYou failed to guess: %s  See more here: https://en.wikipedia.org/wiki/%s \n\n", game.Word, game.Word)
			return "", errors.New("you lost")
		}

		// should move this one to a single method
		if strings.Index(nGame.WordMasked, "_") == -1 {
			return "", fmt.Errorf("you won")
		}

		reply += hangingArt[(len(hangingArt) - int(nGame.RetryLeft) - 1)]
		reply += fmt.Sprintf("\nRemaining attempts: %v", nGame.RetryLeft)
		reply += ("\nIncorrect attempts: ")
		for _, v := range nGame.IncorrectGuesses {
			reply += fmt.Sprint(v.Letter, " ")
		}
		reply += fmt.Sprint("\nWord hint:", nGame.WordMasked)
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

// IsLetter checks if it's a letter
// this is also checked on the server side too, but well behaving client should not burden the server
func isLetter(s string) bool {

	//should be just one
	if len(s) > 1 {
		return false
	}

	//should not be a strange character
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

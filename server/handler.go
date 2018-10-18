package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/AndreiD/HangmanGo2/api"
)

type hangman struct {
	game []*api.Game
}

// NewGame appends a new game in the memory
func (s *hangman) NewGame(ctx context.Context, r *api.GameRequest) (*api.Game, error) {
	if r.RetryLimit < 1 {
		return nil, errors.New("Please specify retry limit for this hangman")
	}
	// pick a random word
	word, err := pickRandomWord("wordlist.txt")
	if err != nil {
		return nil, fmt.Errorf("could not pick a word %s", err)
	}
	wordMAsked := strings.Repeat("_", utf8.RuneCountInString(word))
	GameID := int32(len(s.game))

	if GameID == 0 {
		s.game = append(s.game, &api.Game{Id: 0, Status: true})
		GameID++
	}
	s.game = append(s.game, &api.Game{Id: GameID, Word: word, WordMasked: wordMAsked, RetryLimit: r.RetryLimit, RetryLeft: r.RetryLimit, Status: true})
	g := *s.game[GameID]
	g.Word = ""
	return &g, nil
}

// ListGames lists the games
func (s *hangman) ListGames(context.Context, *api.GameRequest) (*api.GameArray, error) {
	d := &api.GameArray{Game: s.game}
	if d.Game == nil {
		return &api.GameArray{}, nil
	}
	d.Game = d.Game[1:]
	return d, nil
}

// ResumeGame resumes the game
func (s *hangman) ResumeGame(ctx context.Context, r *api.GameRequest) (*api.Game, error) {
	if r.Id > 0 && int32(len(s.game)) > r.Id {
		if s.game[r.Id].RetryLeft < 1 {
			return nil, errors.New("This game is over")
		}
		if s.game[r.Id].Status {
			return nil, errors.New("Game is played by someone else")
		}

		s.game[r.Id].Status = true
		d := *s.game[r.Id]
		d.Word = ""
		return &d, nil
	}
	return nil, errors.New("Invalid Game ID")
}

// SaveGame saves the current game everytime there a new "action"
func (s *hangman) SaveGame(ctx context.Context, r *api.GameRequest) (*api.Game, error) {

	if r.Id > 0 && int32(len(s.game)) > r.Id {
		s.game[r.Id].Status = false
		gg := *s.game[r.Id]
		gg.Word = ""
		return &gg, nil
	}
	return nil, errors.New("Invalid Game ID")
}

// GuessLetter guesses a letter
func (s *hangman) GuessLetter(ctx context.Context, r *api.GuessRequest) (*api.Game, error) {

	if r.GameID > 0 && int32(len(s.game)) > r.GameID {
		r.Letter = strings.ToLower(r.Letter)
		g := s.game[r.GameID]
		if g.RetryLeft < 1 {
			return nil, errors.New("This game is Over")
		}

		for k, v := range g.Word { // expose all letter occurencies
			if v == rune(r.Letter[0]) {
				g.WordMasked = g.WordMasked[:k] + r.Letter + g.WordMasked[k+1:]
			}
		}
		if strings.Index(g.Word, r.Letter) == -1 {
			contains := false
			for _, v := range g.IncorrectGuesses {
				if r.Letter == v.Letter {
					contains = true
				} else {

				}
			}
			if !contains {
				g.IncorrectGuesses = append(g.IncorrectGuesses, &api.GuessRequest{Letter: r.Letter})
				g.RetryLeft = g.RetryLeft - 1
			}
		}
		gg := *g     // need to dereference so we don't change the original struct
		gg.Word = "" // don't sent the naked word to the client , to avoid cheating clients :)
		return &gg, nil
	}
	return nil, errors.New("Invalid Game ID")
}

// PickRandomWord reads a random word from a wordlist file
func pickRandomWord(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read the file %s", err)
	}

	lines := strings.Split(string(content), "\n")

	rand.Seed(time.Now().UnixNano())
	randInt := rand.Intn(len(lines))

	return lines[randInt], nil
}

package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/AndreiD/HangmanGo2/api"
	"github.com/AndreiD/HangmanGo2/server/data"
	"github.com/golang/protobuf/proto"
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
	fmt.Printf("The word for this game is %s \n", word)
	if err != nil {
		return nil, fmt.Errorf("could not pick a word %s", err)
	}
	wordMAsked := strings.Repeat("_", utf8.RuneCountInString(word))
	GameID := int32(len(s.game))

	if GameID == 0 {
		s.game = append(s.game, &api.Game{Id: 0, Status: "ongoing"})
		GameID++
	}
	s.game = append(s.game, &api.Game{Id: GameID, Word: word, WordMasked: wordMAsked, RetryLimit: r.RetryLimit, RetryLeft: r.RetryLimit, Status: "ongoing", PlayerId: r.PlayerId})
	g := *s.game[GameID]
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
		if s.game[r.Id].Status == "won" || s.game[r.Id].Status == "lost" {
			return nil, errors.New("you can't resume a game that ended")
		}

		s.game[r.Id].Status = "ongoing"

		// everytime you resume, you lose a retry point
		s.game[r.Id].RetryLeft = s.game[r.Id].RetryLeft - 1

		d := *s.game[r.Id]
		return &d, nil
	}
	return nil, errors.New("Invalid Game ID")
}

// SaveGame saves the current game everytime there a new "action"
func (s *hangman) SaveGame(ctx context.Context, r *api.GameRequest) (*api.Game, error) {

	if r.Id > 0 && int32(len(s.game)) > r.Id {
		s.game[r.Id].Status = "ongoing"
		gg := *s.game[r.Id]

		// prepare for storage
		marshaled, err := proto.Marshal(s.game[r.Id])
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}
		DB := data.Init()
		DB.Save(fmt.Sprint(r.Id), marshaled)

		return &gg, nil
	}
	return nil, errors.New("Invalid Game ID")
}

// GuessLetter guesses a lett
func (s *hangman) GuessLetter(ctx context.Context, r *api.GuessRequest) (*api.Game, error) {

	if r.GameID > 0 && int32(len(s.game)) > r.GameID {

		// convert letter to lowercase
		r.Letter = strings.ToLower(r.Letter)
		g := s.game[r.GameID]
		if g.RetryLeft < 1 {
			return nil, errors.New("This game is OVER")
		}

		for k, v := range g.Word {
			if v == rune(r.Letter[0]) {
				g.WordMasked = g.WordMasked[:k] + r.Letter + g.WordMasked[k+1:]
			}
		}
		if strings.Index(g.Word, r.Letter) == -1 {
			contains := false
			for _, v := range g.IncorrectGuesses {
				if r.Letter == v.Letter {
					contains = true
				}
			}
			if contains == false {
				g.IncorrectGuesses = append(g.IncorrectGuesses, &api.GuessRequest{Letter: r.Letter})
				g.RetryLeft = g.RetryLeft - 1
			}
		}

		return g, nil
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

	// returns the lowercase word
	return strings.TrimSpace(strings.ToLower(lines[randInt])), nil
}

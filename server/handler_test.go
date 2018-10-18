package main

import (
	"context"
	"log"
	"testing"

	"github.com/AndreiD/HangmanGo2/api"
)

var Game *api.Game
var hangm *hangman

func init() {
	var err error
	hangm = &hangman{}
	if _, err = hangm.NewGame(context.Background(), &api.GameRequest{Id: -1}); err == nil {
		log.Panic("Game initialization without Retry limit should fail")
	}
	Game, err = hangm.NewGame(context.Background(), &api.GameRequest{RetryLimit: 5})
	if err != nil {
		log.Panic("Game initialization returned an error:", err)
	}
}

func TestPickRandomWord(t *testing.T) {
	word, err := pickRandomWord("wordlist.txt")
	if err != nil {
		t.Log("picking random word failed:", err)
		t.Fail()
	}
	if len(word) < 1 {
		t.Log("word length is less than 1!")
		t.Fail()
	}
}

func TestNewGame(t *testing.T) {
	if Game.Id != 1 {
		t.Logf("Game initialization expected ID:%v, actual ID:%v", 1, Game.Id)
		t.Fail()
	}
}

func TestListGames(t *testing.T) {
	l, err := hangm.ListGames(context.Background(), &api.GameRequest{Id: -1})
	if err != nil {
		t.Log("Game listing error:", err)
		t.Fail()
	}
	if len(l.Game) == 0 {
		t.Log("Game listing returned unexpected 0 length")
		t.Fail()
	}
	if l.Game[0].Id != 1 {
		t.Logf("Game listing returned unexpected Id:%v for the first element", l.Game[0].Id)
		t.Fail()
	}
}

func TestResumeGame(t *testing.T) {
	if _, err := hangm.ResumeGame(context.Background(), &api.GameRequest{Id: 1}); err == nil {
		t.Log("Game resume should fail for locked Games")
		t.Fail()
	}
	if _, err := hangm.SaveGame(context.Background(), &api.GameRequest{Id: 1}); err != nil {
		t.Logf("Game save error:%v", err)
		t.Fail()
	}

	g1, err := hangm.ResumeGame(context.Background(), &api.GameRequest{Id: 1})
	if err != nil {
		t.Logf("Game resume error:%v", err)
		t.Fail()
	}

	if g1.Id != 1 {
		t.Logf("Game ID expected:%v, actual:%v", 1, g1.Id)
		t.Fail()
	}
	if _, err := hangm.ResumeGame(context.Background(), &api.GameRequest{Id: -1}); err == nil {
		t.Log("Game didn't fail with an invalid Game ID")
		t.Fail()
	}
}

func TestGuessALetter(t *testing.T) {
	g, err := hangm.SaveGame(context.Background(), &api.GameRequest{Id: Game.Id})
	if err != nil {
		t.Logf("Game save error:%v", err)
		t.Fail()

	}
	gg, err := hangm.GuessLetter(context.Background(), &api.GuessRequest{GameID: 1, Letter: "~"})
	if err != nil {
		t.Logf("Game letter guess error:%v", err)
		t.Fail()
	}
	if g.RetryLeft-gg.RetryLeft != 1 {
		t.Logf("Retry Limit decrease expected:1 actual:%v", (g.RetryLeft - gg.RetryLeft))
		t.Fail()
	}
	if _, err := hangm.GuessLetter(context.Background(), &api.GuessRequest{GameID: -1, Letter: "~"}); err == nil {
		t.Log("Letter guess didn't fail with an invalid Game ID")
		t.Fail()
	}
}

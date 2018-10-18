package utils

import (
	"testing"
)

func TestIsLetter(t *testing.T) {
	isLetter := IsLetter("abc")

	if isLetter == true {
		t.Error("a letter is not a word...")
	}

	isLetter = IsLetter("%")

	if isLetter == true {
		t.Error("a letter is not a special char...")
	}
}

package twitterBot

import (
	"testing"
)

func assertEqualBooleans(t *testing.T, got, want bool) {
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertEqualStrings(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

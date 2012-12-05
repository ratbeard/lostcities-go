package main

import(
	"testing"
)

func TestGame(t *testing.T) {
	g := NewGame()
	if len(g.player1Hand) != 5 {
		t.Error("Wrong number of starting cards for player1")
	}
	if len(g.player2Hand) != 5 {
		t.Error("Wrong number of starting cards for player2")
	}
	
}
package main

import (
	"testing"
)

func TestGame(t *testing.T) {
	g := NewGame()
	
	if len(g.player1Hand) != 5 {
		t.Error("Wrong number of starting cards for player1", len(g.player1Hand))
	}
	
	if len(g.player2Hand) != 5 {
		t.Error("Wrong number of starting cards for player2", len(g.player2Hand))
	}
	
	if len(g.deck) != (65 - 10) {
		t.Error("Wrong number of cards left in deck", len(g.deck))
	}

}

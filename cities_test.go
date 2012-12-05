package main

import (
	"testing"
	"math/rand"
	"time"
)

func TestDealing(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Seed(0)
	g := NewGame()
	
	
	if len(g.player1Hand) != 5 {
		t.Error("Wrong number of starting cards for player1", len(g.player1Hand))
	}
	
	if len(g.player2Hand) != 5 {
		t.Error("Wrong number of starting cards for player2", len(g.player2Hand))
	}
	
	if len(g.deck) != (55) {
		t.Error("Wrong number of cards left in deck", len(g.deck))
	}
	
	if g.currentTurn != "player1" {
		t.Error("player1 should always start ")
	}

}

func TestPlaying(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	g := NewGame()
	_ = g

}
package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// Using rand.seed(0), game should look like:
// player1Hand: [{green 7} {blue 4} {yellow 10} {red 9} {red s}]
// player2Hand: [{yellow 7} {blue 1} {blue 8} {blue 10} {white 1}]
// Remaining cards (reversed here, so top of deck is first):
// 		[ {yellow 2} {white s} {red 1} {white 9} {red 3} {yellow 1} {white s} {yellow s} {red s} {white 3} {red 5} {red 7} {yellow 4} {yellow 6} {white 7} {yellow 8} {blue 5} {blue 3} {blue 9} {blue s} {green s} {red 8} {white 10} {blue s} {green s} {blue 2} {green 10} {green 8} {yellow 3} {green 4} {green 6} {green 5} {green 2} {red 4} {blue s} {red 10} {white 2} {white 5} {green 1} {yellow 9} {white 4} {blue 6} {white 8} {red 6} {white 6} {red 2} {white s} {red s} {yellow s} {yellow 5} {green 3} {yellow s} {green 9} {blue 7} {green s}]
// 
func TestNewGame(t *testing.T) {
	_ = fmt.Println
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

	//fmt.Printf("%#v", g.player1Hand)
	//fmt.Println(g.player2Hand)
	//fmt.Println(g.deck)
}

func TestExpectedDeckOrder(t *testing.T) {
	rand.Seed(0)
	game := NewGame()
	var card Card

	card = Card{"green", "7"}
	if game.player1Hand[0] != card {
		t.Error("player1 unexpected card")
	}

	card = Card{"blue", "8"}
	if game.player2Hand[2] != card {
		t.Error("player2 unexpected card")
	}

	card = Card{"yellow", "s"}
	if game.deck[6] != card {
		t.Error("deck unexpected card")
	}
}

func TestValidMove(t *testing.T) {
	rand.Seed(0)
	game := NewGame()
	var move *Move

	move = &Move{"player1", Card{"green", "7"}, Play, "deck"}
	if !game.validMove(move) {
		t.Error("is valid", move)
	}

	move = &Move{"player1", Card{"green", "7"}, Play, "yellow"}
	if game.validMove(move) {
		t.Error("trying to draw from an empty discard pile is invalid", move)
	}
}

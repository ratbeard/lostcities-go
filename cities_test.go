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

func TestHasCard(t *testing.T) {
	rand.Seed(0)
	game := NewGame()
	var card Card

	card = Card{"green", "7"}
	if !hasCard(game.player1Hand, card) {
		t.Error("Card should be in hand", card)
	}

	if hasCard(game.player2Hand, card) {
		t.Error("Card should not be in hand", card)
	}
}

func TestHighestCard(t *testing.T) {
	var cards []Card
	var card, card1, card2 Card

	card1 = Card{"green", "8"}
	card2 = Card{"green", "7"}
	if !card1.higherThan(card2) {
		t.Error("8 is higher than 7")
	}

	card1 = Card{"green", "10"}
	card2 = Card{"green", "7"}
	if !card1.higherThan(card2) {
		t.Error("10 is higher than 7")
	}

	card1 = Card{"green", "1"}
	card2 = Card{"green", "s"}
	if !card1.higherThan(card2) {
		t.Error("s is higher than 1")
	}

	card1 = Card{"green", "s"}
	card2 = Card{"green", "s"}
	if !card1.higherThan(card2) {
		t.Error("s is higher than s")
	}

	cards = []Card{{"green", "s"}, {"green", "8"}}
	card = Card{"green", "10"}
	if !highestCard(cards, card) {
		t.Error("10 is higher than 8")
	}
}

func TestValidMove(t *testing.T) {
	rand.Seed(0)
	game := NewGame()

	// Playing or Discarding from hand then drawing from deck is a legit first move
	for _, card := range game.player1Hand {
		assertValidMove(t, game, &Move{"player1", card, PlayAction, "deck"})
		assertValidMove(t, game, &Move{"player1", card, DiscardAction, "deck"})
	}

	// Wrong turn, not legit
	assertInvalidMove(t, game, &Move{"player2", game.player2Hand[0], PlayAction, "deck"})

	// Bad player, not legit
	assertInvalidMove(t, game, &Move{"xxx", game.player2Hand[0], PlayAction, "deck"})

	// Playing card not in hand is not legit
	assertInvalidMove(t, game, &Move{"player1", game.player2Hand[0], PlayAction, "deck"})
	assertInvalidMove(t, game, &Move{"player1", game.deck[0], PlayAction, "deck"})

	// Trying something other than Play or Discard is not legit
	assertInvalidMove(t, game, &Move{"player1", game.player1Hand[0], 3, "deck"})

	// Drawing from an empty discard pile is not legit
	assertInvalidMove(t, game, &Move{"player1", game.player1Hand[0], PlayAction, "yellow"})

	// Drawing from a non-empty discard pile is legit
	game.discards["yellow"] = []Card{Card{"yellow", "1"}}
	assertValidMove(t, game, &Move{"player1", game.player1Hand[0], PlayAction, "yellow"})
	game.discards["yellow"] = []Card{}

	// Playing a card thats lower than a card you've already played is not legit
	game.player1Plays["green"] = []Card{Card{"green", "8"}}
	assertInvalidMove(t, game, &Move{"player1", Card{"green", "7"}, PlayAction, "deck"})
	game.player1Plays["green"] = []Card{}

	// Playing a card thats higher than a card you've already played is legit
	game.player1Plays["green"] = []Card{Card{"green", "6"}}
	assertValidMove(t, game, &Move{"player1", Card{"green", "7"}, PlayAction, "deck"})
	game.player1Plays["green"] = []Card{}

	// Playing a card thats lower than a card thats been discarded is legit
	game.discards["green"] = []Card{Card{"green", "8"}}
	assertValidMove(t, game, &Move{"player1", Card{"green", "7"}, PlayAction, "deck"})
	game.discards["green"] = []Card{}

	// Playing a card thats lower than a card that your opponent has played is legit
	game.player2Plays["green"] = []Card{Card{"green", "8"}}
	assertValidMove(t, game, &Move{"player1", Card{"green", "7"}, PlayAction, "deck"})
	game.player2Plays["green"] = []Card{}

}

func TestPlayingInvalidMove(t *testing.T) {
	rand.Seed(0)
	game := NewGame()
	move := &Move{"player2", game.player1Hand[0], PlayAction, "yellow"}

	err := game.PlayMove(move)
	if err == nil {
		t.Error("Playing an invalid move should return an error")
	}
	if game.currentTurn != "player1" {
		t.Error("Playing an invalid move should not advance game turn")
	}
}

func TestPlayingValidMove(t *testing.T) {
	rand.Seed(0)
	game := NewGame()
	move := &Move{"player1", game.player1Hand[0], PlayAction, "deck"}

	err := game.PlayMove(move)
	if err != nil {
		t.Error("Playing a valid move should not return an error", err)
	}
	if game.currentTurn != "player2" {
		t.Error("Playing a valid move should advance game turn")
	}
}

func TestTurnsAndGameEnd(t *testing.T) {
	rand.Seed(0)
	game := NewGame()
	var turn string
	var hand []Card
	var move *Move

	for i := 0; i < 55; i++ {
		if i%2 == 0 {
			turn, hand = "player1", game.player1Hand
		} else {
			turn, hand = "player2", game.player2Hand
		}

		// Check its the correct turn
		if game.currentTurn != turn {
			t.Errorf("Turn #%d should be %s's turn", i, turn)
		}

		// Play a move
		move = &Move{turn, hand[0], PlayAction, "deck"}
		game.PlayMove(move)
	}

	// Check end game state
	if len(game.deck) > 0 {
		t.Errorf("Deck should be empty, %d cards left", len(game.deck))
	}
	if !game.done {
		t.Error("Game should be done")
	}
}

// Helpers
func assertValidMove(t *testing.T, game *Game, move *Move) {
	if err := game.CheckMove(move); err != nil {
		t.Error("Should be a valid move:", move, "got:", err)
	}
}

func assertInvalidMove(t *testing.T, game *Game, move *Move) {
	if err := game.CheckMove(move); err == nil {
		t.Error("Should be an invalid move:", move)
	}
}

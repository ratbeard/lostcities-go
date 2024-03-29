package main

import (
	"fmt"
	"testing"
)

// Using rand.seed(1), game should look like:
// player1Hand: [{yellow 6} {blue 9} {red 9} {green 6} {blue 5} {white s} {yellow 9} {red 3}]
// player2Hand: [{yellow 10} {red 2} {green 9} {blue s} {white 4} {white s} {blue 4} {white 9}]
// Remaining cards (reversed here, so top of deck is first):
// 		[{green 3} {yellow s} {white 10} {red 5} {yellow s} {blue 8} {yellow 4} {red 8} {red 10} {blue 6} {white 3} {green 5} {white s} {yellow 7} {red 7} {green 7} {green s} {green 2} {red 6} {green s} {white 2} {blue s} {blue 3} {blue 2} {yellow 2} {red s} {white 7} {blue s} {red 4} {yellow 5} {red s} {yellow 8} {white 5} {green 4} {green s} {blue 10} {green 8} {white 6} {blue 7} {green 10} {red s} {yellow s} {yellow 3} {white 8}]
// 
func TestNewGame(t *testing.T) {
	_ = fmt.Println
	g := NewGame(1)

	if len(g.player1Hand.Cards) != 8 {
		t.Error("Wrong number of starting cards for player1", len(g.player1Hand.Cards))
	}

	if len(g.player2Hand.Cards) != 8 {
		t.Error("Wrong number of starting cards for player2", len(g.player2Hand.Cards))
	}

	if len(g.deck.Cards) != 44 {
		t.Error("Wrong number of cards left in deck", len(g.deck.Cards))
	}

	if g.currentTurn != "player1" {
		t.Error("player1 should always start ")
	}

  // Print out cards and deck:
	//fmt.Println(g.player1Hand.Cards)
	//fmt.Println(g.player2Hand.Cards)
	//fmt.Println(g.deck.Cards)
  //for i := len(g.deck.Cards)-1; i > -1; i-- {
  //  fmt.Print(g.deck.Cards[i], " ")
  //}
	//panic("smap")
}

func TestExpectedDeckOrder(t *testing.T) {
	game := NewGame(1)
	var card Card

	card = Card{"yellow", "6"}
	if game.player1Hand.Cards[0] != card {
		t.Error("player1 unexpected card", game.player1Hand.Cards)
	}

	card = Card{"green", "9"}
	if game.player2Hand.Cards[2] != card {
		t.Error("player2 unexpected card")
	}

	card = Card{"white", "6"}
	if game.deck.Cards[6] != card {
		t.Error("deck unexpected card")
	}
}

func TestHasCard(t *testing.T) {
	game := NewGame(1)
	var card Card

	card = Card{"green", "6"}
	if !game.player1Hand.Has(card) {
		t.Error("Card should be in hand", card)
	}

	if game.player2Hand.Has(card) {
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
	p := Pile{Cards: cards}
	if !p.IsHighestCard(card) {
		t.Error("10 is higher than 8")
	}
}

func TestValidMove(t *testing.T) {
	game := NewGame(1)

	// Playing or Discarding from hand then drawing from deck is a legit first move
	for _, card := range game.player1Hand.Cards {
		assertValidMove(t, game, &Move{"player1", card, PlayAction, "deck"})
		assertValidMove(t, game, &Move{"player1", card, DiscardAction, "deck"})
	}

	// Wrong turn, not legit
	assertInvalidMove(t, game, &Move{"player2", game.player2Hand.Cards[0], PlayAction, "deck"})

	// Bad player, not legit
	assertInvalidMove(t, game, &Move{"xxx", game.player2Hand.Cards[0], PlayAction, "deck"})

	// Playing card not in hand is not legit
	assertInvalidMove(t, game, &Move{"player1", game.player2Hand.Cards[0], PlayAction, "deck"})
	assertInvalidMove(t, game, &Move{"player1", game.deck.Cards[0], PlayAction, "deck"})

	// Trying something other than Play or Discard is not legit
	assertInvalidMove(t, game, &Move{"player1", game.player1Hand.Cards[0], 3, "deck"})

	// Drawing from an empty discard pile is not legit
	assertInvalidMove(t, game, &Move{"player1", game.player1Hand.Cards[0], PlayAction, "yellow"})

	// Drawing from a non-empty discard pile is legit
	var p *Pile
	p = game.discards["yellow"]
	p.Add(Card{"yellow", "1"})
	assertValidMove(t, game, &Move{"player1", game.player1Hand.Cards[0], PlayAction, "yellow"})
	_, _ = p.Pop()

	// Playing a card thats lower than a card you've already played is not legit
	p = game.player1Plays["green"]
	p.Add(Card{"green", "8"})
	assertInvalidMove(t, game, &Move{"player1", Card{"green", "6"}, PlayAction, "deck"})
	_, _ = p.Pop()

	// Playing a card thats higher than a card you've already played is legit
	p = game.player1Plays["green"]
	p.Add(Card{"green", "5"})
	assertValidMove(t, game, &Move{"player1", Card{"green", "6"}, PlayAction, "deck"})
	_, _ = p.Pop()

	// Playing a card thats lower than a card thats been discarded is legit
	p = game.discards["green"]
	p.Add(Card{"green", "8"})
	assertValidMove(t, game, &Move{"player1", Card{"green", "6"}, PlayAction, "deck"})
	_, _ = p.Pop()

	// Discarding a card thats lower than a played card is legit
	p = game.player1Plays["green"]
	p.Add(Card{"green", "8"})
	assertValidMove(t, game, &Move{"player1", Card{"green", "6"}, DiscardAction, "deck"})
	_, _ = p.Pop()

	// Playing a card thats lower than a card that your opponent has played is legit
	p = game.player2Plays["green"]
	p.Add(Card{"green", "8"})
	assertValidMove(t, game, &Move{"player1", Card{"green", "6"}, PlayAction, "deck"})
	_, _ = p.Pop()
}

func TestPlayingInvalidMove(t *testing.T) {
	game := NewGame(1)
	move := &Move{"player2", game.player1Hand.Cards[0], PlayAction, "yellow"}

	err := game.PlayMove(move)
	if err == nil {
		t.Error("Playing an invalid move should return an error")
	}
	if game.currentTurn != "player1" {
		t.Error("Playing an invalid move should not advance game turn")
	}
}

func TestPlayingValidMove(t *testing.T) {
	game := NewGame(1)
	var card Card
	var move *Move

	// Turn 1 - Player1 plays {green 6}
	card = Card{"green", "6"}
	move = &Move{"player1", card, PlayAction, "deck"}
	err := game.PlayMove(move)
	if err != nil {
		t.Fatal("Playing a valid move should not return an error", err)
	}
	if p := game.player1Plays[card.suit]; !p.Has(card) {
		t.Fatal("Playing a card should put it in the play pile", card, p)
	}
	if p := game.player1Hand; p.Has(card) {
		t.Fatal("Playing a card should remove it from the hand", card, p)
	}

	// Turn 2 - Player2 discards {yellow 10}
	card = game.player2Hand.Cards[0]
	move = &Move{"player2", card, DiscardAction, "deck"}
	err = game.PlayMove(move)
	if err != nil {
		t.Error("Playing a valid move should not return an error", err)
	}
	if len(game.discards[card.suit].Cards) != 1 {
		t.Fatal("Discarding a card should put it in the discard pile")
	}
	if game.discards[card.suit].Cards[0] != card {
		t.Error("Discarding a card should put it in the discard pile")
	}

	// Turn 3 - Player1 discards {yellow 6}
	card = Card{"yellow", "6"}
	move = &Move{"player1", card, DiscardAction, "deck"}
	err = game.PlayMove(move)
	if err != nil {
		t.Error("Playing a valid move should not return an error", err)
	}
	if len(game.discards[card.suit].Cards) != 2 {
		t.Fatal("Discarding a card should put it in the discard pile")
	}
	if game.discards[card.suit].Cards[1] != card {
		t.Error("Discarding a card should put it in the discard pile")
	}

	// Turn 4 - Player2 plays {red 2}, draws from yellow discard
	card = Card{"red", "2"}
	move = &Move{"player2", card, PlayAction, "yellow"}
	err = game.PlayMove(move)
	if err != nil {
		t.Error("Playing a valid move should not return an error", err)
	}
	if len(game.player2Plays[card.suit].Cards) != 1 {
		t.Fatal("Playing a card should put it in the play pile")
	}
	if game.player2Plays[card.suit].Cards[0] != card {
		t.Error("Playing a card should put it in the play pile")
	}
	if len(game.discards[move.drawPile].Cards) != 1 {
		t.Fatal("Drawing from discard pile removes the top card from the pile")
	}
	if game.discards[move.drawPile].Cards[0] == card {
		t.Fatal("Drawing from discard pile removes the top card from the pile")
	}

}

func TestTurnsAndGameEnd(t *testing.T) {
	game := NewGame(1)
	var turn string
	var hand Pile
	var move *Move

	for i := 1; i < 45; i++ {
		if i%2 == 0 {
			turn, hand = "player2", game.player2Hand
		} else {
			turn, hand = "player1", game.player1Hand
		}

		// Check its the correct turn
		if game.done {
			t.Fatal("Didnt expect game to be done")
		}
		if game.currentTurn != turn {
			t.Errorf("Turn #%d should be %s's turn", i, turn)
		}

		// Play a move
		move = &Move{turn, hand.Cards[0], DiscardAction, "deck"}
		e := game.PlayMove(move)
		if e != nil {
			t.Error("Unexpected error while discarding")
		}
	}

	// Check end game state
	if len(game.deck.Cards) > 0 {
		t.Errorf("Deck should be empty, %d cards left", len(game.deck.Cards))
	}
	if !game.done {
		t.Error("Game should be done")
	}
}

// Helpers
func assertValidMove(t *testing.T, game *Game, move *Move) {
	if err := game.CheckMove(move); err != nil {
		t.Fatal("Should be a valid move:", move, "got:", err)
		panic("")
	}
}

func assertInvalidMove(t *testing.T, game *Game, move *Move) {
	if err := game.CheckMove(move); err == nil {
		t.Fatal("Should be an invalid move:", move)
		// panic('x')
	}
}

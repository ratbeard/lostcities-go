package main

import (
	"errors"
	"math/rand"
	"time"
)

const (
	cardInHandCount = 8
	cardCount       = 12 * 5

	PlayAction    = 1
	DiscardAction = 2
)

var (
	Pips  = [...]string{"s", "s", "s", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	Suits = [...]string{"yellow", "white", "blue", "green", "red"}
)

type Card struct {
	suit, pip string
}

type Move struct {
	player   string
	card     Card
	action   int    // Play || Discard
	drawPile string // "deck", "yellow", "green", "etc"
}

// Game -->
type Game struct {
	// Secret state
	deck, player1Hand, player2Hand Pile

	// Board state.  Discards are semi-secret depending on your memory
	discards, player1Plays, player2Plays map[string]*Pile

	// Who's turn is it
	currentTurn string

	// Is the game over
	done bool
}

func NewGame(seed int64) (game *Game) {
	// Randomize.  Kinda janky to use 1 as a sentinel value
	// but I used 0 in the tests.
	if seed == 0 {
		seed = time.Now().UTC().UnixNano()
	}
	rand.Seed(seed)

	game = &Game{
		currentTurn:  "player1",
		deck:         buildShuffledDeck(),
		player1Plays: make(map[string]*Pile),
		player2Plays: make(map[string]*Pile),
		discards:     make(map[string]*Pile),
	}

	// Initialize Pointers
	for _, suit := range Suits {
		game.player1Plays[suit] = &Pile{}
		game.player2Plays[suit] = &Pile{}
		game.discards[suit] = &Pile{}
	}

	// Deal out opening hand
	for i := 0; i < cardInHandCount; i++ {
		game.draw("player1", "deck")
		game.draw("player2", "deck")
	}
	return
}

func (game *Game) CheckMove(move *Move) error {
	if game.done {
		return errors.New("The game is over")
	}

	if move.player != game.currentTurn {
		return errors.New("Wrong turn")
	}

	if !game.handFor(move.player).Has(move.card) {
		return errors.New("Card not in hand")
	}

	if move.action != PlayAction && move.action != DiscardAction {
		return errors.New("Invalid action.  Must be play or discard")
	}

	if move.action == PlayAction {
		p := game.playPileFor(move.player, move.card.suit)
		if !p.IsHighestCard(move.card) {
			return errors.New("A higher card has been played in that pile")
		}
	}

	if p := game.drawPileFor(move.drawPile); len(p.Cards) == 0 {
		return errors.New("Cannot draw from empty pile")
	}

	return nil
}

func (game *Game) PlayMove(move *Move) error {
	if err := game.CheckMove(move); err != nil {
		return err
	}

	// Play/Discard the card
	suit := move.card.suit
	var p *Pile
	if move.action == DiscardAction {
		p = game.discards[suit]
	} else {
		p = game.playPileFor(move.player, suit)
	}
	p.Add(move.card)

	hand := game.handFor(move.player)
	hand.Remove(move.card)

	// Perform the draw
	game.draw(move.player, move.drawPile)

	// Switch turns
	if game.currentTurn == "player1" {
		game.currentTurn = "player2"
	} else {
		game.currentTurn = "player1"
	}

	// Check if game over
	if len(game.deck.Cards) == 0 {
		game.done = true
	}

	return nil
}

func (game *Game) draw(player, pileName string) {
	pile := game.drawPileFor(pileName)
	hand := game.handFor(player)
	card, _ := pile.Pop()
	//fmt.Println("draw", card, pileName, len(p.Cards))
	hand.Add(card)
}

// Kinda janky string comparison.
// Need to rewrite "10" and "s" to next/prev ascii char
// "s".higherThan("s") => true
func (card *Card) higherThan(other Card) bool {
	a, b := card.pip, other.pip
	if a == "s" {
		a = "0"
	}
	if b == "s" {
		b = "0"
	}
	if a == "10" {
		a = ":"
	}
	if b == "10" {
		b = ":"
	}
	return a >= b
}

// Helper getters
func (game *Game) playPileFor(player, suit string) *Pile {
	h := game.player1Plays
	if player == "player2" {
		h = game.player2Plays
	}
	return h[suit]
}

func (game *Game) handFor(name string) *Pile {
	if name == "player1" {
		return &game.player1Hand
	} else if name == "player2" {
		return &game.player2Hand
	}
	return nil
}

func (game *Game) drawPileFor(name string) *Pile {
	if name == "deck" {
		return &game.deck
	}
	return game.discards[name]
}

func calculateScore(hand map[string]*Pile) (score int) {
	for _, pile := range hand {
		score += pile.Score()
	}
	return
}

func buildShuffledDeck() Pile {
	// Build that deck
	unshuffled := make([]Card, cardCount)
	for i, suit := range Suits {
		for j, pip := range Pips {
			unshuffled[i*len(Pips)+j] = Card{suit, pip}
		}
	}

	// Shuffle that deck
	randIndices := rand.Perm(cardCount)
	shuffled := make([]Card, cardCount)
	for i, index := range randIndices {
		shuffled[i] = unshuffled[index]
	}

	return Pile{Cards: shuffled}
}

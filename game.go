package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	cardInHandCount = 5
	cardCount       = 13 * 5

	PlayAction    = 1
	DiscardAction = 2
)

var (
	Pips  = [...]string{"s", "s", "s", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	Suits = [...]string{"yellow", "white", "blue", "green", "red"}
)

type Move struct {
	player   string
	card     Card
	action   int    // Play || Discard
	drawPile string // "yellow", etc, "deck"
}

type Card struct {
	suit, pip string
}

type Deck struct {
	cards []Card
}

type Player struct {
	name string
	hand []Card
}

// Game -->

type Game struct {
	// Secret state
	deck, player1Hand, player2Hand Pile

	// Board, discards are semi-secret depending on your memory
	discards, player1Plays, player2Plays map[string]*Pile

	currentTurn string
	done        bool

	// Questionable:
	player1 *Player
	player2 *Player
}

func NewGame() (game *Game) {
	game = new(Game)

	// Initialize Data structures
	game.player1Plays = make(map[string]*Pile)
	game.player2Plays = make(map[string]*Pile)
	game.discards = make(map[string]*Pile)
	for _, suit := range Suits {
		game.player1Plays[suit] = &Pile{}
		game.player2Plays[suit] = &Pile{}
		game.discards[suit] = &Pile{}
	}

	// Get shuffled deck
	game.deck = buildShuffledDeck()

	// Initialize the players
	game.player1 = new(Player)
	game.player2 = new(Player)

	// Deal out opening hand
	for i := 0; i < cardInHandCount; i++ {
		game.draw("player1", "deck")
		game.draw("player2", "deck")
	}

	// Initialize state
	game.currentTurn = "player1"

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

	//fmt.Println(shuffled)
	/*
		reversed := make([]Card, cardCount)
		for i := 0; i < cardCount; i++ {
			reversed[cardCount-1-i] = shuffled[i] 
		}
		fmt.Println("Top to bottom:", reversed)
	*/

	return Pile{Cards: shuffled}
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

	pile := game.player1Plays[move.card.suit]
	if move.player == "player2" {
		pile = game.player2Plays[move.card.suit]
	}
	if !pile.IsHighestCard(move.card) {
		return errors.New("A higher card has been played in that pile")
	}

	if p := game.pileFor(move.drawPile); len(p.Cards) == 0 {
		return errors.New("Cannot draw from empty pile")
	}

	return nil
}

func (game *Game) PlayMove(move *Move) error {
	if err := game.CheckMove(move); err != nil {
		return err
	}

	// Perform the play/discard
	card := move.card
	var p *Pile
	if move.action == DiscardAction {
		p = game.discards[card.suit]
	} else if move.player == "player1" {
		p = game.player1Plays[card.suit]
	} else if move.player == "player2" {
		p = game.player2Plays[card.suit]
	}
	//fmt.Println("Playing move!", move, p)
	p.Add(card)
	//fmt.Println("now:", p)

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

func (game *Game) handFor(name string) *Pile {
	if name == "player1" {
		return &game.player1Hand
	} else if name == "player2" {
		return &game.player2Hand
	}
	return nil
}

func (game *Game) pileFor(name string) *Pile {
	if name == "deck" {
		return &game.deck
	}
	return game.discards[name]
}

func (game *Game) draw(player, pileName string) {
	p := game.pileFor(pileName)
	card, _ := p.Pop()
	//fmt.Println("draw", card, pileName, len(p.Cards))

	if player == "player1" {
		game.player1Hand.Add(card)
	} else {
		game.player2Hand.Add(card)
	}
}

func calculateScore(hand map[string]*Pile) (score int) {
	for _, pile := range hand {
		score += pile.Score()
	}
	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	/*
		game := new(Game)
		game.player1Plays = map[string][]Card{
			"yellow": []Card{{"yellow", "s"}, {"yellow", "s"}, {"yellow", "s"}, {"yellow", "1"}, {"yellow", "2"}, {"yellow", "3"}, {"yellow", "4"}, {"yellow", "5"}, {"yellow", "6"}, {"yellow", "7"}, {"yellow", "8"}, {"yellow", "9"}, {"yellow", "10"}},
			"white":  []Card{{"white", "s"}, {"white", "s"}, {"white", "10"}},
			"blue":   []Card{{"blue", "3"}, {"blue", "4"}},
			"green":  []Card{},
			"red":    []Card{{"red", "8"}, {"red", "9"}},
		}
		game.discards = map[string][]Card{
			"yellow": []Card{},
			"white":  []Card{{"white", "1"}},
			"blue":   []Card{{"blue", "s"}},
			"green":  []Card{{"green", "s"}, {"green", "1"}, {"green", "2"}},
			"red":    []Card{{"red", "10"}},
		}
		game.player2Plays = map[string][]Card{
			"yellow": []Card{},
			"white":  []Card{{"white", "4"}, {"white", "6"}},
			"blue":   []Card{{"blue", "s"}, {"blue", "1"}, {"blue", "2"}, {"blue", "8"}, {"blue", "10"}},
			"green":  []Card{{"green", "5"}, {"green", "6"}, {"green", "7"}, {"green", "10"}},
			"red":    []Card{{"red", "s"}, {"red", "3"}, {"red", "7"}},
		}

		printScreen(game)
		printScores(game)
	*/
	// fmt.Scan(&i)
	fmt.Println()
}

// Console Game
// ============
func printScores(game *Game) {
	fmt.Print("Player1 score: ")
	printScore(game.player1Plays)
	fmt.Println()

	fmt.Print("Player2 score: ")
	printScore(game.player2Plays)
	fmt.Println()
}

func printScore(plays map[string]*Pile) {
	var score int
	score = calculateScore(plays)
	fmt.Print(justifyRight(strconv.Itoa(score), 4), "  =  ")

	for _, color := range Suits {
		p := plays[color]
		score = p.Score()
		fmt.Print(colorStr(justifyRight(strconv.Itoa(score), 4), shellColors[color]), " ")
	}
}

func printScreen(game *Game) {
	// Clear the screen
	fmt.Print("\033[2J\033[H")

	fmt.Print("             Player 1                     Player 2\n\n")

	for _, color := range Suits {
		printR(game, color)
	}

	fmt.Print("\nYour hand: ", colorStr("6", White), " ", colorStr("10", White), " ", colorStr("8", Green))
	fmt.Print("\n")
	fmt.Print("Pick a card to draw: (w)hite, (y)ellow, (r)ed or (d)eck:\n")
	fmt.Print("Pick a card to play: (1-5)\n")
	fmt.Print("\n")
}

func printR(game *Game, color string) {
	width := 26
	cards := FormatCards(*game.player1Plays[color])
	justified := justifyRight(cards, width)
	colored := colorStr(justified, shellColors[color])
	fmt.Print(colored)

	discards := game.discards[color].Cards
	topDiscard := "   "
	if len(discards) > 0 {
		topDiscard = discards[len(discards)-1].pip
		topDiscard = colorStr(justifyRight(topDiscard, 3), shellColors[color])
	}
	fmt.Print("  |", topDiscard, " |  ")

	fmt.Print(colorStr(FormatCards(*game.player2Plays[color]), shellColors[color]))
	fmt.Println()
}

//
// Coloring / Formatting
// ======
const (
	Reset  = ""
	Yellow = "1;33"
	Blue   = "1;34"
	White  = "1;37"
	Green  = "32"
	Red    = "31"
)

var shellColors = map[string]string{
	"yellow": "1;33",
	"blue":   "1;34",
	"white":  "1;37",
	"green":  "32",
	"red":    "31",
}

/*

func (c Card) String() string {
	//return c.pip
	return colorStr(c.pip, shellColors[c.suit])
}
*/

func FormatCards(p Pile) string {
	c := p.Cards
	s := fmt.Sprint(c)
	return s[1 : len(s)-1]
}

func justifyRight(s string, width int) string {
	return s
	return strings.Repeat(" ", width-len(s)) + s
}

func colorEscape(str string) string {
	return "\033[" + str + "m"
}

func colorStr(str string, c string) string {
	return colorEscape(c) + str + colorEscape(Reset)
}

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var Pips = []string{"s", "s", "s", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
var Suits = []string{"yellow", "white", "blue", "green", "red"}

const cardCount = 13 * 5
const cardInHandCount = 5

const (
	PlayAction    = 1
	DiscardAction = 2
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
	x     int
}

type Player struct {
	name string
	hand []Card
}

type Pile struct {
	Cards []Card
	//Sort bool
}

func (pile *Pile) Has(card Card) (found bool) {
	for _, c := range pile.Cards {
		if c == card {
			return true
		}
	}
	return
}

func (pile *Pile) Pop() (card Card, ok bool) {
	size := len(pile.Cards)
	if size == 0 {
		return Card{}, false
	}
	
	card, pile.Cards = pile.Cards[size-1], pile.Cards[:size-1]
	return card, true
}

func (pile *Pile) Add(card Card) () {
	pile.Cards = append(pile.Cards, card)
}

func (pile *Pile) Remove(card Card) {
	
}

func (pile *Pile) MoveTopCard(other *Pile) {
	card, _ := pile.Pop()
	other.Add(card)
}




type Game struct {
	deck Pile
	
	// Secret state:
	// deck        []Card
	player1Hand []Card
	player2Hand []Card

	// Board (discards are semi-hidden):
	player1Plays map[string][]Card
	player2Plays map[string][]Card
	discards     map[string][]Card

	currentTurn string
	done        bool

	// Questionable:
	player1 *Player
	player2 *Player
}

func NewGame() (game *Game) {
	game = new(Game)

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
	game.player1Plays = make(map[string][]Card)
	game.player2Plays = make(map[string][]Card)
	game.discards = make(map[string][]Card)
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

	return Pile{ Cards: shuffled }
}

func (game *Game) CheckMove(move *Move) error {
	if game.done {
		return errors.New("The game is over")
	}

	if move.player != game.currentTurn {
		return errors.New("Wrong turn")
	}

	if !hasCard(game.handFor(move.player), move.card) {
		return errors.New("Card not in hand")
	}

	if move.action != PlayAction && move.action != DiscardAction {
		return errors.New("Invalid action.  Must be play or discard")
	}

	playPile := game.player1Plays[move.card.suit]
	if !highestCard(playPile, move.card) {
		return errors.New("A higher card has been played in that pile")
	}

	if pile := game.pileFor(move.drawPile); len(pile) == 0 {
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
	if move.action == DiscardAction {
		game.discards[card.suit] = append(game.discards[card.suit], card)
	} else {
		if move.player == "player1" {
			game.player1Plays[card.suit] = append(game.player1Plays[card.suit], card)
		} else {
			game.player2Plays[card.suit] = append(game.player2Plays[card.suit], card)
		}
	}

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

/*
func moveCard(card Card, from, to []Card) ([]Card, []Card){
	return from, to
}
*/

func hasCard(cards []Card, card Card) bool {
	for _, c := range cards {
		if c == card {
			return true
		}
	}
	return false
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

func highestCard(cards []Card, card Card) bool {
	for _, c := range cards {
		if c.higherThan(card) {
			// fmt.Println(card, "Not higher than", c)
			return false
		}
	}
	return true
}

func (game *Game) handFor(name string) []Card {
	if name == "player1" {
		return game.player1Hand
	} else if name == "player2" {
		return game.player2Hand
	}
	return nil
}

func (game *Game) pileFor(name string) []Card {
	if name == "deck" {
		return game.deck.Cards
	}
	return game.discards[name]
}

func (game *Game) draw(player, pileName string) {
	var card Card
	
	if pileName == "deck" {
		card, _ = game.deck.Pop()
	} else {
		pile := game.pileFor(pileName)
		card, pile, _ = pop(pile)
		game.discards[pileName] = pile
	}

	if player == "player1" {
		game.player1Hand = append(game.player1Hand, card)
	} else {
		game.player2Hand = append(game.player2Hand, card)
	}
}

// helper
func pop(cards []Card) (Card, []Card, bool) {
	size := len(cards)
	if size <= 0 {
		return Card{}, cards, false
	}

	return cards[size-1], cards[:size-1], true
}

// Score Calculation
// ========
func calculateScore(hand map[string][]Card) (score int) {
	for _, cards := range hand {
		score += calculateScoreForSuit(cards)
	}
	return
}

func calculateScoreForSuit(cards []Card) (score int) {
	if len(cards) == 0 {
		return
	}
	score = -20
	multiplier := 1
	for _, card := range cards {
		if card.pip == "s" {
			multiplier += 1
		} else {
			value, _ := strconv.Atoi(card.pip)
			score += value
		}
	}
	score *= multiplier
	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

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

func printScore(plays map[string][]Card) {
	var score int
	score = calculateScore(plays)
	fmt.Print(justifyRight(strconv.Itoa(score), 4), "  =  ")

	for _, color := range Suits {
		score = calculateScoreForSuit(plays[color])
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
	cards := FormatCards(game.player1Plays[color])
	justified := justifyRight(cards, width)
	colored := colorStr(justified, shellColors[color])
	fmt.Print(colored)

	discards := game.discards[color]
	topDiscard := "   "
	if len(discards) > 0 {
		topDiscard = discards[len(discards)-1].pip
		topDiscard = colorStr(justifyRight(topDiscard, 3), shellColors[color])
	}
	fmt.Print("  |", topDiscard, " |  ")

	fmt.Print(colorStr(FormatCards(game.player2Plays[color]), shellColors[color]))
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

func FormatCards(c []Card) string {
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

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var Pips = []string{"s", "s", "s", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
var Suits = []string{"yellow", "white", "blue", "green", "red"}

const cardsInHandCount = 5

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

type Game struct {
	// Hidden state:
	deck         []Card
	player1Hand  []Card
	player2Hand  []Card	
	
	// Board (discards are semi-hidden):
	player1Plays map[string][]Card
	player2Plays map[string][]Card
	discards     map[string][]Card
	
	currentTurn string
	
	// Questionable:
	player1      *Player
	player2      *Player
}

func NewGame() (game *Game) {
	game = new(Game)

	// Get Shuffled Deck
	game.deck = buildShuffledDeck()

	// Initialize the Players
	game.player1 = new(Player)
	game.player2 = new(Player)

	// Deal out opening hand
	for i := 0; i < cardsInHandCount; i++ {
		_ = game.drawFromDeck("player1")
		_ = game.drawFromDeck("player2")
	}
	
	game.currentTurn = "player1"
	
	return
}

func buildShuffledDeck() ([]Card) {
	cardCount := len(Pips)*len(Suits)
	unshuffled := make([]Card, cardCount)

	for i, suit := range Suits {
		for j, pip := range Pips {
			unshuffled[i*len(Pips)+j] = Card{suit, pip}
			// fmt.Println(i*len(Pips)+j, Card{suit, pip})
		}
	}

	randIndices := rand.Perm(cardCount)
	shuffled := make([]Card, cardCount)
	for i, index := range randIndices {
		shuffled[index] = unshuffled[i]
	}
	fmt.Println(shuffled)
	
	return shuffled
}

func (game *Game) drawFromDeck(player string) (bool) {
	card, deck, ok := pop(game.deck)
	game.deck = deck
	if (!ok) {
		return ok
	}
	
	fmt.Println("drawFromDeck:", player, card)
	
	if player == "player1" {
		game.player1Hand = append(game.player1Hand, card)
	} else {
		game.player2Hand = append(game.player2Hand, card)
	}
	return true
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


func (c Card) String() string {
	//return c.pip
	return colorStr(c.pip, shellColors[c.suit])
}

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

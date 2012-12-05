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

type Card struct {
	suit, pip string
}

var allCards []Card
/*
var allCards = []Card{
	Card{"yellow", "s"},
	Card{"yellow", "s"},
	Card{"yellow", "s"},
	Card{"yellow", "1"},
	Card{"yellow", "2"},
	Card{"yellow", "3"},
	Card{"yellow", "4"},
	Card{"yellow", "5"},
	Card{"yellow", "6"},
	Card{"yellow", "7"},
	Card{"yellow", "8"},
	Card{"yellow", "9"},
	Card{"yellow", "10"},
	Card{"white", "s"},
	Card{"white", "s"},
	Card{"white", "s"},
	Card{"white", "1"},
	Card{"white", "2"},
	Card{"white", "3"},
	Card{"white", "4"},
	Card{"white", "5"},
	Card{"white", "6"},
	Card{"white", "7"},
	Card{"white", "8"},
	Card{"white", "9"},
	Card{"white", "10"},
	Card{"blue", "s"},
	Card{"blue", "s"},
	Card{"blue", "s"},
	Card{"blue", "1"},
	Card{"blue", "2"},
	Card{"blue", "3"},
	Card{"blue", "4"},
	Card{"blue", "5"},
	Card{"blue", "6"},
	Card{"blue", "7"},
	Card{"blue", "8"},
	Card{"blue", "9"},
	Card{"blue", "10"},
	Card{"green", "s"},
	Card{"green", "s"},
	Card{"green", "s"},
	Card{"green", "1"},
	Card{"green", "2"},
	Card{"green", "3"},
	Card{"green", "4"},
	Card{"green", "5"},
	Card{"green", "6"},
	Card{"green", "7"},
	Card{"green", "8"},
	Card{"green", "9"},
	Card{"green", "10"},
	Card{"red", "s"},
	Card{"red", "s"},
	Card{"red", "s"},
	Card{"red", "1"},
	Card{"red", "2"},
	Card{"red", "3"},
	Card{"red", "4"},
	Card{"red", "5"},
	Card{"red", "6"},
	Card{"red", "7"},
	Card{"red", "8"},
	Card{"red", "9"},
	Card{"red", "10"},
}
*/

type Hand [5]Card
type PlayPile []Card
type DiscardPile []Card
type Deck struct {
	cards []Card
	x     int
}
type Player struct {
	name string
	hand Hand
}
type Game struct {
	player1      *Player
	player2      *Player
	deck         []Card
	player1Plays map[string][]Card
	player2Plays map[string][]Card
	discards     map[string][]Card
	x            int
}

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

func (g Game) init() {
	g.player1 = new(Player)
}

func NewGame() *Game {
	game := new(Game)

	// Shuffle deck
	cards := make([]Card, len(allCards))
	/*
	cards := make([]Card, len(Pips) * len(Suits))
	for _, suit range Suits {
		for _, pip range Pips {
			cards
		}
	}
	*/
	randIndices := rand.Perm(len(allCards))
	for i, index := range randIndices {
		cards[index] = allCards[i]
	}
	game.deck = cards

	player1 := new(Player)
	// Deal out opening hand
	for i := 0; i < 5; i++ {
		game.draw(player1)
	}
	return game
}

func (game *Game) draw(player *Player) {
	// card := game.deck[0]
	// game.deck.unshift()
}

func cool() {
	game := NewGame()
	game.x = 4
	fmt.Println(game)
	// game := new(Game)
	// game.init()
}

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


func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	cool()
	
	game := new(Game)
	game.player1Plays = map[string][]Card{
		"yellow": []Card{{"yellow", "s"}, {"yellow", "s"}, {"yellow", "s"}, {"yellow", "1"}, {"yellow", "2"}, {"yellow", "3"}, {"yellow", "4"}, {"yellow", "5"}, {"yellow", "6"}, {"yellow", "7"},{"yellow", "8"},{"yellow", "9"},{"yellow", "10"}, },
		"white": []Card{{"white", "s"}, {"white", "s"}, {"white", "10"}},
		"blue":  []Card{{"blue", "3"}, {"blue", "4"}},
		"green":  []Card{},
		"red": []Card{{"red", "8"}, {"red", "9"}},
	}
	game.discards = map[string][]Card{
		"yellow": []Card{},
		"white": []Card{{"white", "1"}},
		"blue": []Card{{"blue", "s"}},
		"green": []Card{{"green", "s"}, {"green", "1"}, {"green", "2"}},
		"red": []Card{{"red", "10"}},
	}
	game.player2Plays = map[string][]Card{
		"yellow": []Card{},
		"white": []Card{{"white", "4"}, {"white", "6"}},
		"blue":  []Card{{"blue", "s"}, {"blue", "1"}, {"blue", "2"}, {"blue", "8"}, {"blue", "10"}},
		"green": []Card{{"green", "5"}, {"green", "6"}, {"green", "7"}, {"green", "10"}},
		"red": []Card{{"red", "s"}, {"red", "3"}, {"red", "7"}},
	}
	
	
	printScreen(game)

	printScores(game)

	// fmt.Scan(&i)
	fmt.Println()
}

func printScores(game *Game) {
	fmt.Print("Player1 score: ")
	printScore(game.player1Plays)
	fmt.Println()

	fmt.Print("Player2 score: ")
	printScore(game.player2Plays)
	fmt.Println()
}

func printScore(plays map[string][]Card){
	var score int
	score = calculateScore(plays)
	fmt.Print(justifyRight(strconv.Itoa(score), 4), "  =  ")
	
	for _, color := range Suits {
		score = calculateScoreForSuit(plays[color])
		fmt.Print(colorStr(justifyRight(strconv.Itoa(score), 4), shellColors[color]), " ")
	}
}


func (c Card) String() string {
	return c.pip
	return colorStr(c.pip, shellColors[c.suit])
}

func FormatCards(c []Card) string {
	s := fmt.Sprint(c)
	return s[1:len(s)-1]
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
		topDiscard = discards[len(discards) - 1].pip
		topDiscard = colorStr(justifyRight(topDiscard, 3), shellColors[color])
	}
	fmt.Print("  |", topDiscard, " |  ")
	
	fmt.Print(colorStr(FormatCards(game.player2Plays[color]), shellColors[color]))
	fmt.Println()
}


func justifyRight(s string, width int) string {
	return strings.Repeat(" ", width-len(s)) + s
}

func colorEscape(str string) string {
	return "\033[" + str + "m"
}

func colorStr(str string, c string) string {
	return colorEscape(c) + str + colorEscape(Reset)
}




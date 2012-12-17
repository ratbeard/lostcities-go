package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ConsoleGame struct {
	game        *Game
	currentMove Move
}

func NewConsoleGame() *ConsoleGame {
	return &ConsoleGame{game: NewGame(1)}
}

func (cg *ConsoleGame) Start() {
	for {
		cg.update()
	}
}

func (cg *ConsoleGame) update() {
	// Clear the screen
	fmt.Print("\033[2J\033[H")

	// Print player heading, including active player
	player1 := "Player 1"
	player2 := "Player 2"
	if cg.game.currentTurn == "player1" {
		player1 = underline(player1)
	} else {
		player2 = underline(player2)
	}
	fmt.Printf("%22s%28s\n\n", player1, player2)

	// Print board
	for _, color := range Suits {
		cg.printRow(color)
	}
	fmt.Println()

	// Print hand
	fmt.Print("Your hand: ")
	hand := cg.game.handFor(cg.game.currentTurn)
	for _, card := range hand.Cards {
		fmt.Print(" ", colorCard(card))
	}
	fmt.Println()

	if cg.game.done {
		fmt.Println("Game Over!")
		return
	}
	
	// Update current move's player
	cg.currentMove.player = cg.game.currentTurn

	// Print Prompts
	var s string
	var x int
	switch {
	case cg.currentMove.card.pip == "":
		fmt.Print("Pick a card to play/discard (1-5):  ")
		fmt.Scan(&x)
		cg.currentMove.card = hand.Cards[x-1]
	case cg.currentMove.action == 0:
		fmt.Print("Play or discard (p,d):  ")
		fmt.Scan(&s)
		if s == "p" {
			cg.currentMove.action = PlayAction
		} else if s == "d" {
			cg.currentMove.action = DiscardAction
		}
	case cg.currentMove.drawPile == "":
		var pile string
		fmt.Print("Pick a card to draw (d y b w g r):  \n")
		fmt.Scan(&s)
		switch s {
		case "d":
			pile = "deck"
		case "y":
			pile = "yellow"
		case "b":
			pile = "blue"
		case "w":
			pile = "white"
		case "g":
			pile = "green"
		case "r":
			pile = "red"
		default:
			_ = "error"
		}
		cg.currentMove.drawPile = pile
	default:
		// Try to play the move
		if err := cg.game.PlayMove(&cg.currentMove); err != nil {
			fmt.Printf("\nError!: %s.  (Press enter)\n\n", err)
			fmt.Scanf("%s", &s)
		}
		// Reset current move
		cg.currentMove = Move{}
	}
}

func (cg *ConsoleGame) printRow(color string) {
	width := 26
	cards := FormatCards(*cg.game.player1Plays[color])
	justified := justifyRight(cards, width)
	colored := colorStr(justified, shellColors[color])
	fmt.Print(colored)

	discards := cg.game.discards[color].Cards
	topDiscard := "   "
	if len(discards) > 0 {
		topDiscard = discards[len(discards)-1].pip
		topDiscard = colorStr(justifyRight(topDiscard, 3), shellColors[color])
	}
	fmt.Print("  |", topDiscard, " |  ")

	fmt.Print(colorStr(FormatCards(*cg.game.player2Plays[color]), shellColors[color]))
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



func FormatCards(p Pile) string {
	c := p.Cards
	s := fmt.Sprint(c)
	return s[1 : len(s)-1]
}

func justifyRight(s string, width int) string {
	return s
	return strings.Repeat(" ", width-len(s)) + s
}



// Shell Colors
// ====
var shellColors = map[string]string{
	"yellow": "1;33",
	"blue":   "1;34",
	"white":  "1;37",
	"green":  "32",
	"red":    "31",
}

func underline(str string) string {
	return "<" + str + ">"
}

func colorCard(card Card) string {
	return colorStr(card.pip, card.suit)
}

func colorStr(str, color string) string {
	return colorEscape(shellColors[color]) + str + colorEscape("")
}

func colorEscape(str string) string {
	return "\033[" + str + "m"
}

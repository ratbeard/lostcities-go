package main

import (
	"fmt"
	"time"
	"math/rand"
)

type ConsoleGame struct {
	game        *Game
	currentMove Move
}

func NewConsoleGame() *ConsoleGame {
	rand.Seed(time.Now().UTC().UnixNano())
	g := &ConsoleGame{
		game: NewGame(),
	}
	return g
}

func (cg *ConsoleGame) Start() {
	for !cg.game.done {
		cg.printState()
		cg.getInput()
	}
}

func (cg *ConsoleGame) printState() {
	// Clear the screen
	fmt.Print("\033[2J\033[H")

	// Print player heading, and active
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
	for _, card := range cg.game.handFor(cg.game.currentTurn).Cards {
		fmt.Print(" ", colorCard(card))
	}
	fmt.Println()
	
	// Print Prompts
	if true {
		fmt.Print("Pick a card to play/discard (1-5):  ")
	}
	if !true {
		fmt.Print("Play or discard (p,d):  ")
	}
	if !true {
		fmt.Print("Pick a card to draw (y b w g r):  \n")
	}
	
	fmt.Println()
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

func (cg *ConsoleGame) getInput() {
	var c string
	fmt.Scan(&c)

}

func colorCard(card Card) string {
	return colorStr(card.pip, card.suit)
}

func underline(str string) string {
	return "<" + str + ">"
}

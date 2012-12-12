package main

import (
	"fmt"
)

type ConsoleGame struct {
}

func NewConsoleGame() *ConsoleGame {
	g := &ConsoleGame{}

	return g

}

func (this *ConsoleGame) hi() {
	fmt.Print("hi")
}

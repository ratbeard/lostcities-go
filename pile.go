package main

import (
	// "fmt"
	"strconv"
)

type Pile struct {
	Cards []Card
	//Sort bool
}

func (pile *Pile) index(card Card) int {
	for i, c := range pile.Cards {
		if c == card {
			return i
		}
	}
	return -1
}

func (pile *Pile) Has(card Card) bool {
	return pile.index(card) != -1
}

func (pile *Pile) Remove(card Card) {
	i := pile.index(card)
	pile.Cards = append(pile.Cards[:i], pile.Cards[i+1:]...)
}

func (pile *Pile) Pop() (card Card, ok bool) {
	size := len(pile.Cards)
	if size == 0 {
		return Card{}, false
	}

	card, pile.Cards = pile.Cards[size-1], pile.Cards[:size-1]
	return card, true
}

func (pile *Pile) Add(card Card) {
	pile.Cards = append(pile.Cards, card)
}

func (pile *Pile) MoveTopCardTo(other *Pile) {
	card, _ := pile.Pop()
	other.Add(card)
}

func (pile *Pile) Score() (score int) {
	cards := pile.Cards
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

func (pile *Pile) IsHighestCard(card Card) bool {
	for _, c := range pile.Cards {
		if c.higherThan(card) {
			// fmt.Println(card, "Not higher than", c)
			return false
		}
	}
	return true
}

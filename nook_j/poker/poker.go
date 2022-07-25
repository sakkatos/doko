package poker

import (
	"fmt"
	"strings"
)

func IsPlayer1Win(cards string) bool {
	splitCards := strings.Split(cards, " ")
	//player1 := (splitCards[:5])
	player2 := (splitCards[5:])
	fmt.Println(player2[4][0])
	if player2[4][0] == 'A' {
		return false
	}
	if cards[8 * 3] == 'A' {
		return false
	}
	// if strings.Contains(player1, "A") {
	// 	return true
	// }
	return true
}

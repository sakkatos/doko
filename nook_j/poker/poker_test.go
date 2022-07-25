package poker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Player1WinWithHighCard(t *testing.T) {
	answer := IsPlayer1Win("8C TS KC 9H AC 7D 2S 5D 3S 4S")
	expected := true
	assert.Equal(t, answer, expected)
}

func Test_Player1LooseWithHighCard(t *testing.T) {
	answer := IsPlayer1Win("8C TS KC 9H 4S 7D 2S 5D 3S AC")
	expected := false
	assert.Equal(t, answer, expected)
}

func Test_Player1LooseWithHighCardAgain(t *testing.T) {
	answer := IsPlayer1Win("8C TS KC 9H 4S 7D 2S 5D 4S AC")
	expected := false
	assert.Equal(t, answer, expected)
}

func Test_Player1LooseWithHighCardAgain2(t *testing.T) {
	answer := IsPlayer1Win("8C TS KC 9H 4S 7D 2S 5D AS 2C")
	expected := false
	assert.Equal(t, answer, expected)
}

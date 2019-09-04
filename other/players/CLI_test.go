package players_test

import (
	"strings"
	"testing"

	"github.com/mgxian/tdd-practice/other/players"
)

func TestCLI(t *testing.T) {
	t.Run("record Chris win from user int", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &players.StubPlayerStore{}
		cli := players.NewCLI(playerStore, in)
		cli.PlayPoker()
		players.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record Cleo win from user int", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &players.StubPlayerStore{}
		cli := players.NewCLI(playerStore, in)
		cli.PlayPoker()
		players.AssertPlayerWin(t, playerStore, "Cleo")
	})
}

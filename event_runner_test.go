package frcsim

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestLoadTeams(t *testing.T) {
	s := BuildSimulation("fixtures/teams.csv", "scoring-path", "fixtures/sorting.csv", 2763633600, 12)
	s.Run()
	spew.Dump(s.AverageRanks())
}

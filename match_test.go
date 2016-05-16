package frcsim

import "testing"

func TestMatchScoring(t *testing.T) {
	r1 := Team{ID: 1}
	r1.AddAttribute("score", 1)
	r2 := Team{ID: 2}
	r2.AddAttribute("score", 2)
	scoringFunction := func(attr []string, alliance map[string]float64) map[string]float64 {
		scores := make(map[string]float64)
		scores["total"] = alliance["score"]
		return scores
	}
	match := Match{ID: 1, ScoringFunction: scoringFunction, TeamAttributes: []string{"score"}}
	match.AddRedTeam(r1)
	match.AddRedTeam(r2)
	redScore, _ := match.Score()
	if redScore["total"] != float64(3) {
		t.Error("Scoring isn't correct:", "Expected", 3, "got", redScore["total"])
	}

}

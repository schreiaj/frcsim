package frcsim

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

type EventSimulation struct {
	Teams               []Team                                                `json:"teams"`
	Runs                []Event                                               `json:"runs"`
	MatchesPerTeam      int                                                   `json:"matchesPerTeam"`
	NumRuns             int                                                   `json:"numRuns"`
	QualScoringFunction func([]string, map[string]float64) map[string]float64 `json:"-"`
	// ElimScoringFunction func([]string, map[string]float64) map[string]float64
	TeamAttributes []string `json:"attributes"`
	SortOrder      []string `json:"sortOrder"`
}

func BuildSimulation(teamCsvPath string, scoringJSPath string, sortingCsvPath string, runs int, matches int) EventSimulation {
	sim := EventSimulation{}
	sim.Teams, sim.TeamAttributes = LoadTeams(teamCsvPath)
	sim.MatchesPerTeam = matches
	sim.NumRuns = runs
	sim.SortOrder = LoadSortingOrder(sortingCsvPath)
	sim.QualScoringFunction = func(attr []string, alliance map[string]float64) map[string]float64 {
		scores := make(map[string]float64)
		scores["total"] = alliance["score_1"] + alliance["score_2"]
		scores["score_2"] = alliance["score_2"]
		return scores
	}
	sim.Runs = make([]Event, sim.NumRuns)
	for i := range sim.Runs {
		sim.Runs[i].QualScoringFunction = sim.QualScoringFunction
		sim.Runs[i].Teams = sim.permuteTeams()
		sim.Runs[i].SortOrder = sim.SortOrder
		sim.Runs[i].TeamAttributes = sim.TeamAttributes
		sim.Runs[i].BuildSchedule(sim.MatchesPerTeam)
	}
	return sim
}

func (s *EventSimulation) Run() {

	for i := range s.Runs {
		s.Runs[i].ScoreEvent()
		s.Runs[i].RankEvent()
	}
}

func (s *EventSimulation) permuteTeams() []Team {
	order := rand.Perm(len(s.Teams))
	permTeams := make([]Team, len(s.Teams))
	for i := range order {
		permTeams[i] = s.Teams[order[i]]
	}
	return permTeams
}

func LoadTeams(path string) ([]Team, []string) {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("No file exists!"))
	}
	defer file.Close()
	reader := csv.NewReader(file)
	csvLines, err := reader.ReadAll()
	teams := make([]Team, len(csvLines)-1)
	teamAttrs := csvLines[0]
	for i, line := range csvLines[1:] {
		teams[i].ID = i
		for j, attr := range teamAttrs {
			val, _ := strconv.ParseFloat(line[j], 64)
			teams[i].AddAttribute(attr, val)
		}
	}
	return teams, teamAttrs
}

func LoadSortingOrder(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("No file exists!"))
	}
	defer file.Close()
	reader := csv.NewReader(file)
	csvLines, err := reader.ReadAll()
	return csvLines[0]
}

func (s *EventSimulation) AverageRanks() map[string]float64 {
	teams := make(map[string]float64)
	for _, run := range s.Runs {
		for i, team := range run.Rankings {
			teams[strconv.Itoa(team.Team.ID)] = teams[strconv.Itoa(team.Team.ID)] + float64(i+1)/float64(s.NumRuns)
		}
	}
	return teams

}

package frcsim

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/olekukonko/tablewriter"
)

type MatchResult struct {
	Team    Team               `json:"team"`
	Results map[string]float64 `json:"result"`
	Match   Match              `json:"match"`
	Win     bool               `json:"win"`
	Tie     bool               `json:"tie"`
}

type TeamRanking struct {
	Team      Team               `json:"team"`
	Breakdown map[string]float64 `json:"breakdown"`
}

type Event struct {
	Teams               []Team                                                `json:"teams"`
	Schedule            []Match                                               `json:"schedule"`
	Results             []MatchResult                                         `json:"results"`
	QualScoringFunction func([]string, map[string]float64) map[string]float64 `json:"-"`
	ElimScoringFunction func([]string, map[string]float64) map[string]float64 `json:"-"`
	TeamAttributes      []string                                              `json:"-"`
	SortOrder           []string                                              `json:"-"`
	Rankings            []TeamRanking                                         `json:"rankings"`
}

func (e *Event) BuildSchedule(matches int) {
	file, err := Asset(fmt.Sprintf("%s/%d_%d.csv", "schedules", len(e.Teams), matches))
	if err != nil {
		panic(fmt.Errorf("No schedule template exists for %d teams and %d matches", len(e.Teams), matches))
	}
	reader := csv.NewReader(bytes.NewReader(file))
	csvLines, err := reader.ReadAll()
	e.Schedule = []Match{}
	for match := 0; match < len(csvLines); match++ {
		currentMatch := Match{ID: match, ScoringFunction: e.QualScoringFunction, TeamAttributes: e.TeamAttributes}
		for team := 0; team < 3; team++ {
			redIndex, _ := strconv.Atoi(csvLines[match][team*2])
			blueIndex, _ := strconv.Atoi(csvLines[match][team*2+6])
			currentMatch.AddRedTeam(e.Teams[redIndex-1])
			currentMatch.AddBlueTeam(e.Teams[blueIndex-1])
		}
		e.Schedule = append(e.Schedule, currentMatch)
	}
}

func (e *Event) ScoreEvent() {
	e.Results = []MatchResult{}
	for _, match := range e.Schedule {
		redScores, blueScores := match.Score()
		for _, team := range match.RedAlliance {
			result := MatchResult{Team: team, Match: match, Results: redScores, Win: redScores["total"] > blueScores["total"], Tie: redScores["total"] == blueScores["total"]}
			e.Results = append(e.Results, result)
		}
		for _, team := range match.BlueAlliance {
			result := MatchResult{Team: team, Match: match, Results: blueScores, Win: blueScores["total"] > redScores["total"], Tie: redScores["total"] == blueScores["total"]}
			e.Results = append(e.Results, result)
		}
	}
}

func (e *Event) RankEvent() {
	teamData := make(map[string]map[string]float64)
	teamsHash := make(map[string]Team)
	for _, res := range e.Results {
		teamNo := strconv.Itoa(res.Team.ID)
		teamsHash[teamNo] = res.Team
		currentData := teamData[teamNo]
		for _, key := range e.SortOrder {
			if teamData[teamNo] == nil {
				teamData[teamNo] = make(map[string]float64)
			}
			teamData[teamNo][key] = currentData[key] + res.Results[key]
			// spew.Dump(teamData[teamNo])
		}

	}
	e.Rankings = make([]TeamRanking, len(e.Teams))
	i := 0
	for name, team := range teamData {
		e.Rankings[i] = TeamRanking{Team: teamsHash[name], Breakdown: team}
		i++
	}
	rankingKeys := make([]string, len(e.Rankings))
	rankingHash := make(map[string]TeamRanking, len(e.Rankings))
	for i, team := range e.Rankings {
		rankingString := e.breakdownToString(team.Breakdown, team.Team)
		rankingKeys[i] = rankingString
		rankingHash[rankingString] = team
	}
	sort.Strings(rankingKeys)
	tempRankings := make([]TeamRanking, len(e.Teams))
	numTeams := len(e.Teams) - 1
	for i := range rankingKeys {
		// Yeah, just going backwards through this
		tempRankings[i] = rankingHash[rankingKeys[numTeams-i]]
	}
	e.Rankings = tempRankings
}

func (e *Event) breakdownToString(breakdown map[string]float64, t Team) string {
	var s string
	for _, attr := range e.SortOrder {
		s = fmt.Sprintf("%s|%09.4f", s, breakdown[attr])
	}
	s = fmt.Sprintf("%s|%04d", s, t.ID)
	return s
}

func (e *Event) PrintSchedule() {
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"Red1", "Red2", "Red3", "Blue1", "Blue2", "Blue3"}
	table.SetHeader(header)
	for _, m := range e.Schedule {
		table.Append(m.TeamIDs())
	}
	table.Render()
}

func (e *Event) PrintRankings() {
	table := tablewriter.NewWriter(os.Stdout)
	header := append([]string{"team"}, e.SortOrder...)
	table.SetHeader(header)
	for _, m := range e.Rankings {
		data := make([]string, len(e.SortOrder)+1)
		data[0] = strconv.Itoa(m.Team.ID)
		i := 1
		for _, val := range e.SortOrder {
			data[i] = fmt.Sprintf("%09.4f", m.Breakdown[val])
			i++
		}

		table.Append(data)
	}
	table.Render()
}

func (e *Event) PrintResults() {
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"Red1", "Red2", "Red3", "Blue1", "Blue2", "Blue3", "Red Score", "Blue Score"}
	table.SetHeader(header)
	for _, m := range e.Schedule {
		redScore, blueScore := m.Score()
		spew.Dump(redScore)
		redString := fmt.Sprintf("%09.4f", redScore["score_2"])
		blueString := fmt.Sprintf("%09.4f", blueScore["score_2"])

		table.Append(append(m.TeamIDs(), redString, blueString))
	}
	table.Render()
}

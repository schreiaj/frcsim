package frcsim

// func TestBuildSchedule(t *testing.T) {
// 	assert := assert.New(t)
// 	teams := make([]Team, 18)
// 	for i := range teams {
// 		teams[i].ID = i + 1
// 		teams[i].AddAttribute("score", float64(i))
// 	}
// 	event := Event{Teams: teams}
// 	event.QualScoringFunction = func(attr []string, alliance map[string]float64) map[string]float64 {
// 		scores := make(map[string]float64)
// 		scores["total"] = alliance["score"]
// 		return scores
// 	}
// 	event.TeamAttributes = []string{"score"}
// 	event.SortOrder = []string{"total"}
//
// 	event.BuildSchedule(1)
// 	assert.NotNil(event.Schedule[0], "Schedule was not generated")
// 	assert.NotNil(event.Schedule[0].RedAlliance, "Team was not placed in schedule")
// 	assert.Equal(event.Schedule[0].RedAlliance[0].ID, 4, "Team was incorrectly placed")
// 	event.ScoreEvent()
// 	event.RankEvent()
// }

package frcsim

import "strconv"

type Match struct {
	ID              int                                                   `json:"matchId"`
	RedAlliance     []Team                                                `json:"blueAlliance"`
	BlueAlliance    []Team                                                `json:"redAlliance"`
	TeamAttributes  []string                                              `json:"-"`
	ScoringFunction func([]string, map[string]float64) map[string]float64 `json:"-"`
}

func (m *Match) AddRedTeam(t Team) {
	m.RedAlliance = append(m.RedAlliance, t)
}

func (m *Match) AddBlueTeam(t Team) {
	m.BlueAlliance = append(m.BlueAlliance, t)
}

func (m *Match) RedAttributes() map[string]float64 {
	return m.extractAllianceAttributes(m.RedAlliance)
}

func (m *Match) BlueAttributes() map[string]float64 {
	return m.extractAllianceAttributes(m.BlueAlliance)
}

func (m *Match) extractAllianceAttributes(alliance []Team) map[string]float64 {
	attributes := make(map[string]float64)

	for _, key := range m.TeamAttributes {
		tempAttr := 0.0

		for _, team := range alliance {
			tempAttr += team.GetAttribute(key)
		}

		attributes[key] = tempAttr
	}

	return attributes
}

func (m *Match) Score() (map[string]float64, map[string]float64) {
	red := m.ScoringFunction(m.TeamAttributes, m.RedAttributes())
	blue := m.ScoringFunction(m.TeamAttributes, m.BlueAttributes())
	return red, blue

}

func (m *Match) TeamIDs() []string {
	teams := make([]string, len(m.RedAlliance)+len(m.BlueAlliance))
	for i, t := range m.RedAlliance {
		teams[i] = strconv.Itoa(t.ID)
	}
	for i, t := range m.BlueAlliance {
		teams[i+3] = strconv.Itoa(t.ID)
	}
	return teams
}

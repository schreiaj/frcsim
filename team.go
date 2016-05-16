package frcsim

type Team struct {
	ID         int                `json:"id"`
	Attributes map[string]float64 `json:"attributes"`
}

func (t *Team) GetAttribute(name string) float64 {
	return t.Attributes[name]
}

func (t *Team) AddAttribute(name string, value float64) {
	if t.Attributes == nil {
		t.Attributes = make(map[string]float64)
	}
	t.Attributes[name] = value
}

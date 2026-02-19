package engine

import (
	"encoding/json"
	"time"
)

// ─── Parsed effet ──────────────────────────────────────

// ParsedEffet represents a parsed bonus/malus effect from JSON.
type ParsedEffet struct {
	Stat  string  `json:"stat"`
	Op    string  `json:"op"` // add, multiply, coeff
	Value float64 `json:"value"`
}

// ParseEffet parses a JSON effet string. Returns nil on failure.
func ParseEffet(raw *string) *ParsedEffet {
	if raw == nil || *raw == "" {
		return nil
	}
	var pe ParsedEffet
	if err := json.Unmarshal([]byte(*raw), &pe); err != nil {
		return nil
	}
	if pe.Stat == "" || pe.Op == "" {
		return nil
	}
	return &pe
}

// ─── Active Modifier ───────────────────────────────────

// ActiveModifier represents a runtime modifier on a tama.
type ActiveModifier struct {
	ID             string      `json:"id"`
	SourceName     string      `json:"sourceName"`
	SourceType     string      `json:"sourceType"`     // bonus, malus
	SourceCategory string      `json:"sourceCategory"` // trait, race, sickness, event, lifechoice
	Effet          ParsedEffet `json:"effet"`
	StartTime      time.Time   `json:"startTime"`
	EndTime        *time.Time  `json:"endTime"`
}

// IsExpired checks if a modifier has expired.
func (m *ActiveModifier) IsExpired(now time.Time) bool {
	if m.EndTime == nil {
		return false
	}
	return now.After(*m.EndTime) || now.Equal(*m.EndTime)
}

// ─── Stat Modifiers (aggregated) ───────────────────────

// StatModifiers holds aggregated coefficient, additive, and multiplier maps.
type StatModifiers struct {
	Coefficients map[string]float64 `json:"coefficients"`
	Additives    map[string]float64 `json:"additives"`
	Multipliers  map[string]float64 `json:"multipliers"`
}

// NewStatModifiers creates a zeroed StatModifiers.
func NewStatModifiers() *StatModifiers {
	return &StatModifiers{
		Coefficients: make(map[string]float64),
		Additives:    make(map[string]float64),
		Multipliers:  make(map[string]float64),
	}
}

// AggregateModifiers aggregates active modifiers into a StatModifiers.
func AggregateModifiers(mods []ActiveModifier, now time.Time) *StatModifiers {
	result := NewStatModifiers()
	for _, mod := range mods {
		if mod.IsExpired(now) {
			continue
		}
		e := mod.Effet
		switch e.Op {
		case "coeff":
			prev, ok := result.Coefficients[e.Stat]
			if !ok {
				prev = 1.0
			}
			result.Coefficients[e.Stat] = prev * e.Value
		case "add":
			result.Additives[e.Stat] += e.Value
		case "multiply":
			prev, ok := result.Multipliers[e.Stat]
			if !ok {
				prev = 1.0
			}
			result.Multipliers[e.Stat] = prev * e.Value
		}
	}
	return result
}

// GetCoeff returns the coefficient for a stat (default 1.0).
func GetCoeff(mods *StatModifiers, stat string) float64 {
	if mods == nil {
		return 1.0
	}
	v, ok := mods.Coefficients[stat]
	if !ok {
		return 1.0
	}
	return v
}

// GetAdditive returns the additive for a stat (default 0).
func GetAdditive(mods *StatModifiers, stat string) float64 {
	if mods == nil {
		return 0
	}
	return mods.Additives[stat]
}

// GetMultiplier returns the multiplier for a stat (default 1.0).
func GetMultiplier(mods *StatModifiers, stat string) float64 {
	if mods == nil {
		return 1.0
	}
	v, ok := mods.Multipliers[stat]
	if !ok {
		return 1.0
	}
	return v
}

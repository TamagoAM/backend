package engine

import (
	"fmt"
	"math/rand"
)

// ─── Sickness rolling ──────────────────────────────────

// SicknessResult holds the result of a sickness roll.
type SicknessResult struct {
	BecameSick   bool
	SicknessName *string
	Sickness     *DBSickness
	Events       []TickEvent
}

// RollSicknessCheck rolls for sickness (DB-driven with fallback).
func RollSicknessCheck(stats *TamaStats, alreadySick bool, dbSicknesses []DBSickness) SicknessResult {
	if alreadySick {
		return SicknessResult{}
	}

	// DB-driven sickness roll
	if len(dbSicknesses) > 0 {
		sick := rollAcquiredSickness(dbSicknesses, stats.Hygiene, stats.Hunger, false)
		if sick != nil {
			msg := fmt.Sprintf("Oh no! Your tama caught %s! %s", sick.Name, sick.Desc)
			return SicknessResult{
				BecameSick:   true,
				SicknessName: &sick.Name,
				Sickness:     sick,
				Events:       []TickEvent{{Type: "sickness", Message: msg}},
			}
		}
		return SicknessResult{}
	}

	// Fallback to hardcoded sicknesses
	chance := SicknessChanceRandom
	if stats.Hygiene < SicknessHygieneThreshold {
		chance += SicknessChanceLowHygiene
	}
	if stats.Hunger < 20 {
		chance += 0.05
	}

	if rand.Float64() < chance {
		name := RandomSicknesses[rand.Intn(len(RandomSicknesses))]
		msg := fmt.Sprintf("Oh no! Your tama caught %s!", name)
		return SicknessResult{
			BecameSick:   true,
			SicknessName: &name,
			Events:       []TickEvent{{Type: "sickness", Message: msg}},
		}
	}

	return SicknessResult{}
}

// rollAcquiredSickness rolls for an acquired DB sickness.
func rollAcquiredSickness(sicknesses []DBSickness, hygiene, hunger float64, alreadySick bool) *DBSickness {
	if alreadySick {
		return nil
	}

	chance := 0.02
	if hygiene < 25 {
		chance += 0.15
	}
	if hunger < 20 {
		chance += 0.05
	}

	if rand.Float64() > chance {
		return nil
	}

	var acquired []DBSickness
	for _, s := range sicknesses {
		if s.Type == "acquired" || s.Type == "both" {
			acquired = append(acquired, s)
		}
	}
	if len(acquired) == 0 {
		return nil
	}

	return weightedPickSickness(acquired)
}

// weightedPickSickness picks a sickness weighted by severity.
func weightedPickSickness(sicknesses []DBSickness) *DBSickness {
	weights := make([]int, len(sicknesses))
	total := 0
	for i, s := range sicknesses {
		switch s.Severity {
		case "mild":
			weights[i] = 6
		case "moderate":
			weights[i] = 3
		case "severe":
			weights[i] = 1
		default:
			weights[i] = 3
		}
		total += weights[i]
	}

	roll := rand.Intn(total)
	for i, w := range weights {
		roll -= w
		if roll < 0 {
			return &sicknesses[i]
		}
	}
	return &sicknesses[len(sicknesses)-1]
}

// ─── Car accident rolling ──────────────────────────────

// RollCarAccident rolls for a random car accident.
func RollCarAccident(stats TamaStats) (TamaStats, []TickEvent) {
	if rand.Float64() < CarAccidentChance {
		stats.CarAccident++
		stats.Money -= CarAccidentMoneyCost
		if stats.Money < 0 {
			stats.Money = 0
		}
		stats.Hygiene = Clamp(stats.Hygiene - 10)
		msg := fmt.Sprintf("Car accident! Medical bills cost %d coins.", CarAccidentMoneyCost)
		return stats, []TickEvent{{Type: "car_accident", Message: msg}}
	}
	return stats, nil
}

// ─── DB Event rolling ──────────────────────────────────

// RollDBEvent rolls for a random DB event appropriate for the stage.
func RollDBEvent(events []DBEvent, currentStage LifeStage) *DBEvent {
	if rand.Float64() > EventChancePerTick {
		return nil
	}

	var eligible []DBEvent
	for _, e := range events {
		if e.Scope != "individual" {
			continue
		}
		if e.MinStage != nil && StageIndex(currentStage) < StageIndex(*e.MinStage) {
			continue
		}
		eligible = append(eligible, e)
	}
	if len(eligible) == 0 {
		return nil
	}

	return weightedPickEvent(eligible)
}

func weightedPickEvent(events []DBEvent) *DBEvent {
	weights := make([]int, len(events))
	total := 0
	for i, e := range events {
		switch e.Severity {
		case "minor":
			weights[i] = 6
		case "major":
			weights[i] = 3
		case "catastrophic":
			weights[i] = 1
		default:
			weights[i] = 3
		}
		total += weights[i]
	}

	roll := rand.Intn(total)
	for i, w := range weights {
		roll -= w
		if roll < 0 {
			return &events[i]
		}
	}
	return &events[len(events)-1]
}

// ─── Life Choice rolling ───────────────────────────────

// RollLifeChoiceCheck rolls for a life choice presentation.
func RollLifeChoiceCheck(choices []DBLifeChoice, currentStage LifeStage, choicesMade map[int]bool) *DBLifeChoice {
	// Check if any choice for this stage was already made
	for _, c := range choices {
		if c.Stage == currentStage {
			if choicesMade[c.ID] {
				return nil
			}
		}
	}

	if rand.Float64() > LifeChoiceChancePerTick {
		return nil
	}

	var available []DBLifeChoice
	for _, c := range choices {
		if !choicesMade[c.ID] && c.Stage == currentStage {
			available = append(available, c)
		}
	}
	if len(available) == 0 {
		return nil
	}

	return weightedPickChoice(available)
}

func weightedPickChoice(choices []DBLifeChoice) *DBLifeChoice {
	weights := make([]int, len(choices))
	total := 0
	for i, c := range choices {
		switch c.Rarity {
		case "common":
			weights[i] = 6
		case "uncommon":
			weights[i] = 3
		case "rare":
			weights[i] = 1
		default:
			weights[i] = 3
		}
		total += weights[i]
	}

	roll := rand.Intn(total)
	for i, w := range weights {
		roll -= w
		if roll < 0 {
			return &choices[i]
		}
	}
	return &choices[len(choices)-1]
}

// ─── Process all random events ─────────────────────────

// EventResults holds the aggregated results from processing random events.
type EventResults struct {
	Stats           TamaStats
	Events          []TickEvent
	BecameSick      bool
	SicknessName    *string
	Sickness        *DBSickness
	TriggeredEvent  *DBEvent
	TriggeredChoice *DBLifeChoice
}

// ProcessRandomEvents processes all random events for a tick.
func ProcessRandomEvents(stats TamaStats, alreadySick bool, ctx *GameContext) EventResults {
	result := EventResults{Stats: stats}

	var dbSicknesses []DBSickness
	var dbEvents []DBEvent
	var dbChoices []DBLifeChoice
	var currentStage LifeStage
	var choicesMade map[int]bool

	if ctx != nil {
		dbSicknesses = ctx.DBSicknesses
		dbEvents = ctx.DBEvents
		dbChoices = ctx.DBChoices
		currentStage = ctx.CurrentStage
		choicesMade = ctx.ChoicesMade
	}

	// Sickness check
	sickResult := RollSicknessCheck(&result.Stats, alreadySick, dbSicknesses)
	result.Events = append(result.Events, sickResult.Events...)
	result.BecameSick = sickResult.BecameSick
	result.SicknessName = sickResult.SicknessName
	result.Sickness = sickResult.Sickness

	// Car accident
	updatedStats, accEvents := RollCarAccident(result.Stats)
	result.Stats = updatedStats
	result.Events = append(result.Events, accEvents...)

	// DB events
	if len(dbEvents) > 0 && currentStage != "" {
		event := RollDBEvent(dbEvents, currentStage)
		if event != nil {
			result.TriggeredEvent = event
			msg := fmt.Sprintf("📰 %s: %s", event.Name, event.Desc)
			result.Events = append(result.Events, TickEvent{Type: "event_triggered", Message: msg})
		}
	}

	// Life choices
	if len(dbChoices) > 0 && currentStage != "" && choicesMade != nil {
		choice := RollLifeChoiceCheck(dbChoices, currentStage, choicesMade)
		if choice != nil {
			result.TriggeredChoice = choice
			msg := fmt.Sprintf("🎯 A life choice approaches: %s", choice.Name)
			result.Events = append(result.Events, TickEvent{Type: "life_choice", Message: msg})
		}
	}

	return result
}

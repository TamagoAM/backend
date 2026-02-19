package engine

import (
	"math"
	"time"
)

// Clamp constrains a value between StatMin and StatMax.
func Clamp(value float64) float64 {
	return math.Min(math.Max(value, StatMin), StatMax)
}

// ComputeElapsedHours returns the hours elapsed since lastActivity, capped at MaxCatchupMs.
func ComputeElapsedHours(lastActivity time.Time, now time.Time) float64 {
	elapsedMs := float64(now.Sub(lastActivity).Milliseconds())
	if elapsedMs > MaxCatchupMs {
		elapsedMs = MaxCatchupMs
	}
	if elapsedMs <= 0 {
		return 0
	}
	return elapsedMs / (60 * 60 * 1000)
}

// GetLastActivityTime returns the most recent activity timestamp from stats.
func GetLastActivityTime(stats *TamaStats) time.Time {
	var latest time.Time
	candidates := []*time.Time{stats.LastFed, stats.LastPlayed, stats.LastCleaned, stats.LastWorked}
	for _, c := range candidates {
		if c != nil && c.After(latest) {
			latest = *c
		}
	}
	return latest
}

// ApplyDecay applies time-based decay to stats using the same formulas as the frontend.
//
// Modifier integration:
//   - "allDecay" coeff: multiplies ALL decay rates
//   - "hunger"/"boredom"/"hygiene" coeff: per-stat multiplier
//   - additive adjustments: added per-hour on top of base decay
func ApplyDecay(stats TamaStats, hours float64, isSick bool, mods *StatModifiers) TamaStats {
	if hours <= 0 {
		return stats
	}

	sickMult := 1.0
	if isSick {
		sickMult = SicknessDecayMultiplier
	}

	allDecayCoeff := GetCoeff(mods, "allDecay")

	hungerCoeff := GetCoeff(mods, "hunger")
	boredomCoeff := GetCoeff(mods, "boredom")
	hygieneCoeff := GetCoeff(mods, "hygiene")

	hungerAdd := GetAdditive(mods, "hunger")
	boredomAdd := GetAdditive(mods, "boredom")
	hygieneAdd := GetAdditive(mods, "hygiene")

	hungerLoss := (HungerDecayPerHour * hungerCoeff * allDecayCoeff * sickMult * hours) - (hungerAdd * hours)
	boredomGain := (BoredomRisePerHour * boredomCoeff * allDecayCoeff * sickMult * hours) - (boredomAdd * hours)
	hygieneLoss := (HygieneDecayPerHour * hygieneCoeff * allDecayCoeff * sickMult * hours) - (hygieneAdd * hours)

	// Work accident recovery
	accidentRecovery := int(math.Floor(hours / WorkAccidentRecoveryHours))
	newWorkAccident := stats.WorkAccident - accidentRecovery
	if newWorkAccident < 0 {
		newWorkAccident = 0
	}
	newCarAccident := stats.CarAccident - accidentRecovery
	if newCarAccident < 0 {
		newCarAccident = 0
	}

	// Work effort decay
	effortDecay := int(math.Floor(hours / WorkEffortDecayHours))
	newWorked := stats.Worked - effortDecay
	if newWorked < 0 {
		newWorked = 0
	}

	// Satisfaction passive decay
	satisLoss := SatisDecayPerHour * hours

	result := stats
	result.Hunger = Clamp(stats.Hunger - hungerLoss)
	result.Boredom = Clamp(stats.Boredom + boredomGain)
	result.Hygiene = Clamp(stats.Hygiene - hygieneLoss)
	result.SocialSatis = Clamp(stats.SocialSatis - satisLoss)
	result.WorkSatis = Clamp(stats.WorkSatis - satisLoss)
	result.PersonalSatis = Clamp(stats.PersonalSatis - satisLoss)
	result.WorkAccident = newWorkAccident
	result.CarAccident = newCarAccident
	result.Worked = newWorked
	// Money does NOT decay

	return result
}

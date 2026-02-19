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

// ComputeNightSplit splits elapsed time between sleeping (lights off + night),
// awake-at-night (lights on + night), and daytime hours.
// This uses a simplified model: if it's currently night, and lights are off,
// all elapsed time since last tick is considered sleeping (approximation).
// For more accuracy over long catch-ups, it would need hour-by-hour simulation.
func ComputeNightSplit(hours float64, ctx *GameContext, loc *time.Location, now time.Time) NightCycleResult {
	if ctx == nil || hours <= 0 {
		return NightCycleResult{DayHours: hours}
	}

	// Determine current time in user's timezone
	localNow := now.In(loc)
	currentHour := localNow.Hour()
	isNightNow := IsNightHour(currentHour)

	if !isNightNow {
		// It's daytime now. For a rough split, estimate how many of the
		// elapsed hours were night vs day by walking backwards.
		nightHours := estimateNightHours(localNow, hours)
		dayHours := hours - nightHours
		if dayHours < 0 {
			dayHours = 0
		}
		if nightHours > 0 && ctx.LightsOff {
			return NightCycleResult{
				NightHoursSleeping: nightHours,
				DayHours:           dayHours,
			}
		} else if nightHours > 0 {
			return NightCycleResult{
				NightHoursAwake: nightHours,
				DayHours:        dayHours,
			}
		}
		return NightCycleResult{DayHours: hours}
	}

	// It's nighttime now
	if ctx.LightsOff {
		// All elapsed time treated as sleeping
		return NightCycleResult{NightHoursSleeping: hours}
	}
	// Lights are ON at night → penalty
	return NightCycleResult{NightHoursAwake: hours}
}

// estimateNightHours walks backwards from 'now' and counts how many
// of the last 'hours' fell within the night window.
func estimateNightHours(localNow time.Time, hours float64) float64 {
	totalMinutes := hours * 60
	nightMinutes := 0.0
	cursor := localNow.Add(-time.Duration(totalMinutes) * time.Minute)
	for m := 0.0; m < totalMinutes; m += 5 { // 5-minute granularity
		t := cursor.Add(time.Duration(m) * time.Minute)
		if IsNightHour(t.Hour()) {
			nightMinutes += 5
		}
	}
	return nightMinutes / 60.0
}

// ApplyDecay applies time-based decay to stats using the same formulas as the frontend.
// Now includes night cycle: sleeping pauses decay + regens happiness,
// lights-on-at-night applies 1.5x happiness decay.
//
// Modifier integration:
//   - "allDecay" coeff: multiplies ALL decay rates
//   - "hunger"/"boredom"/"hygiene" coeff: per-stat multiplier
//   - additive adjustments: added per-hour on top of base decay
func ApplyDecay(stats TamaStats, hours float64, isSick bool, mods *StatModifiers) TamaStats {
	return ApplyDecayWithNight(stats, hours, isSick, mods, nil, nil, time.Now())
}

// ApplyDecayWithNight applies decay with optional night cycle processing.
func ApplyDecayWithNight(stats TamaStats, hours float64, isSick bool, mods *StatModifiers, ctx *GameContext, loc *time.Location, now time.Time) TamaStats {
	if hours <= 0 {
		return stats
	}

	// Compute night/day split
	var nightSplit NightCycleResult
	if ctx != nil && loc != nil {
		nightSplit = ComputeNightSplit(hours, ctx, loc, now)
	} else {
		nightSplit = NightCycleResult{DayHours: hours}
	}

	result := stats

	// 1. Process sleeping hours (lights off at night): no decay, happiness regens
	if nightSplit.NightHoursSleeping > 0 {
		result.Happiness = Clamp(result.Happiness + SleepHappinessRegenPerHour*nightSplit.NightHoursSleeping)
	}

	// 2. Process awake-at-night hours (lights on at night): normal decay + 1.5x happiness penalty
	if nightSplit.NightHoursAwake > 0 {
		result = applyBaseDecay(result, nightSplit.NightHoursAwake, isSick, mods, NightPenaltyHappinessMultiplier)
	}

	// 3. Process daytime hours: normal decay
	if nightSplit.DayHours > 0 {
		result = applyBaseDecay(result, nightSplit.DayHours, isSick, mods, 1.0)
	}

	return result
}

// applyBaseDecay applies standard decay for a given number of hours with an optional
// happiness multiplier (1.0 = normal, 1.5 = night penalty).
func applyBaseDecay(stats TamaStats, hours float64, isSick bool, mods *StatModifiers, happinessMultiplier float64) TamaStats {
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

	// Satisfaction passive decay (with optional night penalty on happiness)
	satisLossHappiness := SatisDecayPerHour * hours * happinessMultiplier

	result := stats
	result.Hunger = Clamp(stats.Hunger - hungerLoss)
	result.Boredom = Clamp(stats.Boredom + boredomGain)
	result.Hygiene = Clamp(stats.Hygiene - hygieneLoss)
	result.SocialSatis = Clamp(stats.SocialSatis - satisLossHappiness)
	result.WorkSatis = Clamp(stats.WorkSatis - satisLossHappiness)
	result.PersonalSatis = Clamp(stats.PersonalSatis - satisLossHappiness)
	result.WorkAccident = newWorkAccident
	result.CarAccident = newCarAccident
	result.Worked = newWorked
	// Money does NOT decay

	return result
}

package engine

import "time"

// ProcessTick processes a time-based tick for a tama.
// It applies decay, rolls for random events, computes satisfaction/happiness,
// and checks for death.
func ProcessTick(stats TamaStats, isSick bool, now time.Time, ctx *GameContext) TickResult {
	var events []TickEvent

	// 1. Elapsed hours from last tick (or last activity if no tick recorded)
	lastActivity := GetLastActivityTime(&stats)
	hours := ComputeElapsedHours(lastActivity, now)

	// 2. Get modifiers
	var mods *StatModifiers
	if ctx != nil {
		mods = ctx.Mods
	}

	// 3. Apply decay
	current := ApplyDecay(stats, hours, isSick, mods)

	// 4. Random events
	eventResult := ProcessRandomEvents(current, isSick, ctx)
	current = eventResult.Stats
	events = append(events, eventResult.Events...)

	nowSick := isSick || eventResult.BecameSick
	var sickName *string
	if eventResult.SicknessName != nil {
		sickName = eventResult.SicknessName
	} else if isSick && ctx != nil && ctx.CurrentSickness != nil {
		sickName = &ctx.CurrentSickness.Name
	}

	// 5. Compute satisfaction & happiness
	var friends *FriendContext
	if ctx != nil {
		friends = ctx.Friends
	}
	social, work, personal := ComputeAllSatisfaction(&current, mods, friends)
	current.SocialSatis = social
	current.WorkSatis = work
	current.PersonalSatis = personal
	happiness := HappinessFromSatis(social, work, personal)
	current.Happiness = happiness

	// 6. Death check
	isDead := happiness <= DeathHappinessThreshold && current.Hunger <= 0

	if isDead {
		events = append(events, TickEvent{
			Type:    "death",
			Message: "Your tama has passed away… 💀",
		})
	}

	return TickResult{
		Stats:        current,
		IsDead:       isDead,
		IsSick:       nowSick,
		SicknessName: sickName,
		Happiness:    happiness,
		Events:       events,
	}
}

// ProcessTickForHours processes a specific number of hours of decay for a tama.
// Used by the background ticker where we know the exact elapsed time.
func ProcessTickForHours(stats TamaStats, hours float64, isSick bool, ctx *GameContext) TickResult {
	var events []TickEvent

	var mods *StatModifiers
	if ctx != nil {
		mods = ctx.Mods
	}

	// Apply decay for the given hours
	current := ApplyDecay(stats, hours, isSick, mods)

	// Random events
	eventResult := ProcessRandomEvents(current, isSick, ctx)
	current = eventResult.Stats
	events = append(events, eventResult.Events...)

	nowSick := isSick || eventResult.BecameSick
	var sickName *string
	if eventResult.SicknessName != nil {
		sickName = eventResult.SicknessName
	} else if isSick && ctx != nil && ctx.CurrentSickness != nil {
		sickName = &ctx.CurrentSickness.Name
	}

	// Satisfaction & happiness
	var friends *FriendContext
	if ctx != nil {
		friends = ctx.Friends
	}
	social, work, personal := ComputeAllSatisfaction(&current, mods, friends)
	current.SocialSatis = social
	current.WorkSatis = work
	current.PersonalSatis = personal
	happiness := HappinessFromSatis(social, work, personal)
	current.Happiness = happiness

	// Death check
	isDead := happiness <= DeathHappinessThreshold && current.Hunger <= 0
	if isDead {
		events = append(events, TickEvent{
			Type:    "death",
			Message: "Your tama has passed away… 💀",
		})
	}

	return TickResult{
		Stats:        current,
		IsDead:       isDead,
		IsSick:       nowSick,
		SicknessName: sickName,
		Happiness:    happiness,
		Events:       events,
	}
}

// ComputeLifeStage computes the current life stage from a birthday.
func ComputeLifeStage(birthday time.Time, now time.Time) LifeStage {
	ageHours := now.Sub(birthday).Hours()
	current := StageInfancy
	for _, b := range StageBoundaries {
		if ageHours >= b.StartHour {
			current = b.Stage
		}
	}
	return current
}

package engine

import "time"

// TamaStats mirrors the core mutable stats stored in Tama_stats.
type TamaStats struct {
	Fed         int        `json:"fed"`
	LastFed     *time.Time `json:"lastFed"`
	Played      int        `json:"played"`
	LastPlayed  *time.Time `json:"lastPlayed"`
	Cleaned     int        `json:"cleaned"`
	LastCleaned *time.Time `json:"lastCleaned"`
	Worked      int        `json:"worked"`
	LastWorked  *time.Time `json:"lastWorked"`

	Hunger  float64 `json:"hunger"`
	Boredom float64 `json:"boredom"`
	Hygiene float64 `json:"hygiene"`
	Money   int     `json:"money"`

	CarAccident  int `json:"carAccident"`
	WorkAccident int `json:"workAccident"`

	SocialSatis   float64 `json:"socialSatis"`
	WorkSatis     float64 `json:"workSatis"`
	PersonalSatis float64 `json:"personalSatis"`
	Happiness     float64 `json:"happiness"`
}

// TickEvent represents a loggable event from a tick.
type TickEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// TickResult is the result of a single engine tick.
type TickResult struct {
	Stats        TamaStats   `json:"stats"`
	IsDead       bool        `json:"isDead"`
	IsSick       bool        `json:"isSick"`
	SicknessName *string     `json:"sicknessName"`
	Happiness    float64     `json:"happiness"`
	Events       []TickEvent `json:"events"`
}

// FriendContext holds friend data for social satisfaction.
type FriendContext struct {
	AliveFriends int `json:"aliveFriends"`
	DeadFriends  int `json:"deadFriends"`
}

// DBSickness mirrors the Sickness DB model used by the engine.
type DBSickness struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Desc           string  `json:"desc"`
	Type           string  `json:"type"`     // congenital, acquired, both
	Severity       string  `json:"severity"` // mild, moderate, severe
	ExpirationDays *int    `json:"expirationDays"`
	CureCost       *int    `json:"cureCost"`
	Bonus          *string `json:"bonus"`
	Malus          *string `json:"malus"`
}

// DBEvent mirrors the Event DB model.
type DBEvent struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Desc     string     `json:"desc"`
	Severity string     `json:"severity"` // minor, major, catastrophic
	Scope    string     `json:"scope"`    // individual, global
	MinStage *LifeStage `json:"minStage"`
	Bonus    *string    `json:"bonus"`
	Malus    *string    `json:"malus"`
}

// DBLifeChoice mirrors the LifeChoice DB model.
type DBLifeChoice struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	Stage      LifeStage `json:"stage"`
	Rarity     string    `json:"rarity"` // common, uncommon, rare
	ChoiceType string    `json:"choiceType"`
	Traits     *string   `json:"traits"`
	Bonus      *string   `json:"bonus"`
	Malus      *string   `json:"malus"`
}

// GameContext holds all DB-driven data needed for a tick.
type GameContext struct {
	Mods            *StatModifiers
	DBSicknesses    []DBSickness
	DBEvents        []DBEvent
	DBChoices       []DBLifeChoice
	CurrentStage    LifeStage
	ChoicesMade     map[int]bool // set of life choice IDs
	CurrentSickness *DBSickness
	Friends         *FriendContext

	// Night cycle
	LightsOff bool   // user turned lights off
	IsNight   bool   // current time is in night window (22:00-10:00 user local)
	Timezone  string // IANA timezone string
}

// NightCycleResult stores what happened during night processing.
type NightCycleResult struct {
	NightHoursSleeping float64 // hours with lights off (paused + regen)
	NightHoursAwake    float64 // hours with lights on at night (penalty)
	DayHours           float64 // normal daytime hours
}

// IsNightHour checks if a given hour (0-23) is within the night window.
func IsNightHour(hour int) bool {
	// Night window: 22:00 - 10:00 (wraps midnight)
	if NightStartHour > NightEndHour {
		return hour >= NightStartHour || hour < NightEndHour
	}
	return hour >= NightStartHour && hour < NightEndHour
}

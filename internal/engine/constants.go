package engine

// ─── Time ──────────────────────────────────────────────
// TickIntervalMs is the tick interval in milliseconds (5 minutes).
const TickIntervalMs = 5 * 60 * 1000

// MaxCatchupMs is the maximum elapsed time to process in one catch-up (24h).
const MaxCatchupMs = 24 * 60 * 60 * 1000

// ─── Decay rates (per hour) ────────────────────────────
const HungerDecayPerHour = 1.04
const BoredomRisePerHour = 1.04
const HygieneDecayPerHour = 0.83
const SatisDecayPerHour = 0.5

// ─── Sickness multiplier ──────────────────────────────
const SicknessDecayMultiplier = 1.5

// ─── Action effects ───────────────────────────────────
type ActionEffect struct {
	Hunger             float64
	Hygiene            float64
	Boredom            float64
	MoneyCost          int
	MoneyEarned        int
	WorkAccidentChance float64
}

var ActionEffects = map[string]ActionEffect{
	"feed":         {Hunger: 15, Hygiene: -3, Boredom: 0, MoneyCost: 5, MoneyEarned: 0, WorkAccidentChance: 0},
	"feed_cookie":  {Hunger: 10, Hygiene: -2, Boredom: -8, MoneyCost: 3, MoneyEarned: 0, WorkAccidentChance: 0},
	"feed_steak":   {Hunger: 25, Hygiene: -5, Boredom: 0, MoneyCost: 8, MoneyEarned: 0, WorkAccidentChance: 0},
	"play":         {Hunger: -5, Hygiene: -3, Boredom: -20, MoneyCost: 0, MoneyEarned: 0, WorkAccidentChance: 0},
	"play_alone":   {Hunger: -3, Hygiene: -2, Boredom: -15, MoneyCost: 0, MoneyEarned: 0, WorkAccidentChance: 0},
	"play_friend":  {Hunger: -5, Hygiene: -3, Boredom: -25, MoneyCost: 0, MoneyEarned: 0, WorkAccidentChance: 0},
	"play_outside": {Hunger: -8, Hygiene: 5, Boredom: -20, MoneyCost: 0, MoneyEarned: 0, WorkAccidentChance: 0},
	"clean":        {Hunger: 0, Hygiene: 20, Boredom: 5, MoneyCost: 0, MoneyEarned: 0, WorkAccidentChance: 0},
	"work":         {Hunger: -8, Hygiene: 0, Boredom: 8, MoneyCost: 0, MoneyEarned: 25, WorkAccidentChance: 0.03},
}

// ─── Sickness thresholds ──────────────────────────────
const SicknessHygieneThreshold = 25
const SicknessChanceLowHygiene = 0.15
const SicknessChanceRandom = 0.02

// ─── Death ────────────────────────────────────────────
const DeathHappinessThreshold = 0.0

// ─── Satisfaction constants ───────────────────────────
const WorkSatisMoneyCAP = 500

// ─── Random events ────────────────────────────────────
const CarAccidentChance = 0.005
const WorkCommuteAccidentChance = 0.01
const CarAccidentMoneyCost = 50

// ─── Work decay ───────────────────────────────────────
const WorkAccidentRecoveryHours = 48.0
const WorkEffortDecayHours = 24.0

// ─── Stat bounds ──────────────────────────────────────
const StatMin = 0.0
const StatMax = 100.0

// ─── Friend social bonus/malus ────────────────────────
const FriendAliveBonus = 5.0
const FriendDeadMalus = 3.0
const FriendBonusCap = 25.0
const NoFriendsMalus = -8.0

// ─── DB event probabilities ──────────────────────────
const EventChancePerTick = 0.008
const LifeChoiceChancePerTick = 0.015
const CongenitalChance = 0.15

// ─── Night Cycle ─────────────────────────────────
const NightStartHour = 22 // 10 PM local time
const NightEndHour = 10   // 10 AM local time

// When lights are off at night: decay paused + happiness regen
const SleepHappinessRegenPerHour = 0.5

// When lights are ON at night (user forgot): 1.5x happiness decay
const NightPenaltyHappinessMultiplier = 1.5

// ─── Notification stat thresholds ────────────────
const NotifLowHungerThreshold = 20.0
const NotifLowHappinessThreshold = 20.0
const NotifHighBoredomThreshold = 80.0
const NotifLowHygieneThreshold = 20.0

// ─── Random sicknesses (fallback) ────────────────────
var RandomSicknesses = []string{
	"Common Cold", "Tummy Ache", "Pixel Flu", "Sniffles",
	"Rash", "Dizzy Spell", "Fever", "Hiccups",
}

// ─── Life Stages ──────────────────────────────────────
type LifeStage string

const (
	StageInfancy    LifeStage = "infancy"
	StageChildhood  LifeStage = "childhood"
	StageTeenage    LifeStage = "teenage"
	StageEarlyAdult LifeStage = "earlyAdult"
	StageMidAdult   LifeStage = "midAdult"
	StageLateAdult  LifeStage = "lateAdult"
	StageElderly    LifeStage = "elderly"
)

type StageBoundary struct {
	Stage     LifeStage
	StartHour float64
}

var StageBoundaries = []StageBoundary{
	{StageInfancy, 0},
	{StageChildhood, 24},
	{StageTeenage, 96},
	{StageEarlyAdult, 192},
	{StageMidAdult, 336},
	{StageLateAdult, 480},
	{StageElderly, 624},
}

const TotalLifespanHours = 3360

func StageIndex(stage LifeStage) int {
	for i, b := range StageBoundaries {
		if b.Stage == stage {
			return i
		}
	}
	return 0
}

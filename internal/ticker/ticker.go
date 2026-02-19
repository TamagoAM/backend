package ticker

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"tamagoam/internal/engine"
)

// AliveTama holds a joined row of Tama + TamaStat for the background ticker.
type AliveTama struct {
	// Tama fields
	TamaID   int        `db:"TamaId"`
	UserID   int        `db:"UserId"`
	Name     string     `db:"Name"`
	Race     string     `db:"Race"`
	Sickness *string    `db:"Sickness"`
	Birthday *time.Time `db:"Birthday"`
	Traits   *string    `db:"Traits"`

	// TamaStat fields
	TamaStatID    int        `db:"TamaStatId"`
	Fed           int        `db:"Fed"`
	LastFed       *time.Time `db:"LastFed"`
	Played        int        `db:"Played"`
	LastPlayed    *time.Time `db:"LastPlayed"`
	Cleaned       int        `db:"Cleaned"`
	LastCleaned   *time.Time `db:"LastCleaned"`
	Worked        int        `db:"Worked"`
	LastWorked    *time.Time `db:"LastWorked"`
	Hunger        int        `db:"Hunger"`
	Boredom       int        `db:"Boredom"`
	Hygiene       int        `db:"Hygiene"`
	Money         int        `db:"Money"`
	CarAccident   int        `db:"CarAccident"`
	WorkAccident  int        `db:"WorkAccident"`
	SocialSatis   float64    `db:"SocialSatis"`
	WorkSatis     float64    `db:"WorkSatis"`
	PersonalSatis float64    `db:"PersonalSatis"`
	Happiness     float64    `db:"Happiness"`
	LastTickAt    *time.Time `db:"LastTickAt"`
}

// DBSicknessRow mirrors the Sickness table for loading game data.
type DBSicknessRow struct {
	SicknessID     int     `db:"SicknessId"`
	Name           string  `db:"Name"`
	Desc           *string `db:"Desc"`
	Type           string  `db:"Type"`
	Severity       string  `db:"Severity"`
	ExpirationDays *int    `db:"ExpirationDays"`
	CureCost       *int    `db:"CureCost"`
	Bonus          *string `db:"Bonus"`
	Malus          *string `db:"Malus"`
}

// DBEventRow mirrors the Event table.
type DBEventRow struct {
	EventID  int     `db:"EventId"`
	Name     string  `db:"Name"`
	Desc     *string `db:"Desc"`
	Severity string  `db:"Severity"`
	Scope    string  `db:"Scope"`
	MinStage *string `db:"MinStage"`
	Bonus    *string `db:"Bonus"`
	Malus    *string `db:"Malus"`
}

// DBLifeChoiceRow mirrors the LifeChoices table.
type DBLifeChoiceRow struct {
	LifeChoicesID int     `db:"LifeChoicesId"`
	Name          string  `db:"Name"`
	Desc          *string `db:"Desc"`
	Stage         string  `db:"Stage"`
	Rarity        string  `db:"Rarity"`
	ChoiceType    string  `db:"ChoiceType"`
	Traits        *string `db:"Traits"`
	Bonus         *string `db:"Bonus"`
	Malus         *string `db:"Malus"`
}

// Ticker runs the background game engine loop.
type Ticker struct {
	db       *sqlx.DB
	interval time.Duration
	stop     chan struct{}
}

// New creates a new Ticker.
func New(db *sqlx.DB, interval time.Duration) *Ticker {
	return &Ticker{
		db:       db,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins the background ticker goroutine.
func (t *Ticker) Start() {
	go func() {
		log.Printf("[ticker] started — interval %v", t.interval)
		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				t.tick()
			case <-t.stop:
				log.Println("[ticker] stopped")
				return
			}
		}
	}()
}

// Stop stops the background ticker.
func (t *Ticker) Stop() {
	close(t.stop)
}

// tick processes one tick cycle for all alive tamas.
func (t *Ticker) tick() {
	ctx := context.Background()
	now := time.Now()

	// 1. Load all alive tamas with their stats
	tamas, err := t.loadAliveTamas(ctx)
	if err != nil {
		log.Printf("[ticker] error loading alive tamas: %v", err)
		return
	}
	if len(tamas) == 0 {
		return
	}

	// 2. Load game data (sicknesses, events, life choices)
	sicknesses := t.loadSicknesses(ctx)
	events := t.loadEvents(ctx)
	choices := t.loadLifeChoices(ctx)

	// 3. Process each tama
	processed := 0
	deaths := 0
	sickened := 0

	for _, tama := range tamas {
		result, err := t.processTama(ctx, tama, now, sicknesses, events, choices)
		if err != nil {
			log.Printf("[ticker] error processing tama %d (%s): %v", tama.TamaID, tama.Name, err)
			continue
		}
		processed++
		if result.IsDead {
			deaths++
		}
		if result.IsSick && (tama.Sickness == nil || *tama.Sickness == "") {
			sickened++
		}
	}

	if deaths > 0 || sickened > 0 {
		log.Printf("[ticker] processed %d tamas — %d deaths, %d new sicknesses", processed, deaths, sickened)
	}
}

// processTama processes a single tama's tick.
func (t *Ticker) processTama(
	ctx context.Context,
	tama AliveTama,
	now time.Time,
	sicknesses []engine.DBSickness,
	events []engine.DBEvent,
	choices []engine.DBLifeChoice,
) (*engine.TickResult, error) {
	// Compute elapsed hours since last tick
	var lastTick time.Time
	if tama.LastTickAt != nil {
		lastTick = *tama.LastTickAt
	} else {
		// Fallback: use last activity time or birthday
		lastTick = engine.GetLastActivityTime(&engine.TamaStats{
			LastFed:     tama.LastFed,
			LastPlayed:  tama.LastPlayed,
			LastCleaned: tama.LastCleaned,
			LastWorked:  tama.LastWorked,
		})
		if lastTick.IsZero() && tama.Birthday != nil {
			lastTick = *tama.Birthday
		}
	}

	hours := engine.ComputeElapsedHours(lastTick, now)
	if hours <= 0 {
		return &engine.TickResult{Stats: toEngineStats(tama)}, nil
	}

	isSick := tama.Sickness != nil && *tama.Sickness != ""

	// Build game context
	var currentStage engine.LifeStage
	if tama.Birthday != nil {
		currentStage = engine.ComputeLifeStage(*tama.Birthday, now)
	}

	// Load life choice history for this tama
	choicesMade := t.loadChoicesMade(ctx, tama.TamaID)

	// Load friend context
	friends := t.loadFriendContext(ctx, tama.UserID)

	gameCtx := &engine.GameContext{
		DBSicknesses: sicknesses,
		DBEvents:     events,
		DBChoices:    choices,
		CurrentStage: currentStage,
		ChoicesMade:  choicesMade,
		Friends:      friends,
	}

	if isSick {
		// Find the current sickness from DB
		for i, s := range sicknesses {
			if s.Name == *tama.Sickness {
				gameCtx.CurrentSickness = &sicknesses[i]
				break
			}
		}
	}

	// Process the tick
	result := engine.ProcessTickForHours(toEngineStats(tama), hours, isSick, gameCtx)

	// Write back to DB
	if err := t.updateStats(ctx, tama.TamaStatID, &result.Stats, now); err != nil {
		return nil, err
	}

	// Handle death
	if result.IsDead {
		if err := t.markDead(ctx, tama.TamaID, now); err != nil {
			log.Printf("[ticker] error marking tama %d dead: %v", tama.TamaID, err)
		}
	}

	// Handle new sickness
	if result.IsSick && (tama.Sickness == nil || *tama.Sickness == "") && result.SicknessName != nil {
		if err := t.setSickness(ctx, tama.TamaID, *result.SicknessName); err != nil {
			log.Printf("[ticker] error setting sickness for tama %d: %v", tama.TamaID, err)
		}
	}

	return &result, nil
}

// ─── Helpers ───────────────────────────────────────────

func toEngineStats(t AliveTama) engine.TamaStats {
	return engine.TamaStats{
		Fed:           t.Fed,
		LastFed:       t.LastFed,
		Played:        t.Played,
		LastPlayed:    t.LastPlayed,
		Cleaned:       t.Cleaned,
		LastCleaned:   t.LastCleaned,
		Worked:        t.Worked,
		LastWorked:    t.LastWorked,
		Hunger:        float64(t.Hunger),
		Boredom:       float64(t.Boredom),
		Hygiene:       float64(t.Hygiene),
		Money:         t.Money,
		CarAccident:   t.CarAccident,
		WorkAccident:  t.WorkAccident,
		SocialSatis:   t.SocialSatis,
		WorkSatis:     t.WorkSatis,
		PersonalSatis: t.PersonalSatis,
		Happiness:     t.Happiness,
	}
}

// ─── DB Queries ────────────────────────────────────────

func (t *Ticker) loadAliveTamas(ctx context.Context) ([]AliveTama, error) {
	var tamas []AliveTama
	err := t.db.SelectContext(ctx, &tamas, `
		SELECT
			t.TamaId, t.UserId, t.Name, t.Race, t.Sickness, t.Birthday, t.Traits,
			ts.TamaStatId, ts.Fed, ts.LastFed, ts.Played, ts.LastPlayed,
			ts.Cleaned, ts.LastCleaned, ts.Worked, ts.LastWorked,
			ts.Hunger, ts.Boredom, ts.Hygiene, ts.Money,
			ts.CarAccident, ts.WorkAccident,
			ts.SocialSatis, ts.WorkSatis, ts.PersonalSatis, ts.Happiness,
			ts.LastTickAt
		FROM Tama t
		JOIN Tama_stats ts ON t.TamaStatsID = ts.TamaStatId
		WHERE t.DeathDay IS NULL
	`)
	return tamas, err
}

func (t *Ticker) updateStats(ctx context.Context, statID int, stats *engine.TamaStats, now time.Time) error {
	_, err := t.db.ExecContext(ctx, `
		UPDATE Tama_stats SET
			Fed = ?, LastFed = ?, Played = ?, LastPlayed = ?,
			Cleaned = ?, LastCleaned = ?, Worked = ?, LastWorked = ?,
			Hunger = ?, Boredom = ?, Hygiene = ?, Money = ?,
			CarAccident = ?, WorkAccident = ?,
			SocialSatis = ?, WorkSatis = ?, PersonalSatis = ?, Happiness = ?,
			LastTickAt = ?
		WHERE TamaStatId = ?`,
		stats.Fed, stats.LastFed, stats.Played, stats.LastPlayed,
		stats.Cleaned, stats.LastCleaned, stats.Worked, stats.LastWorked,
		int(stats.Hunger), int(stats.Boredom), int(stats.Hygiene), stats.Money,
		stats.CarAccident, stats.WorkAccident,
		stats.SocialSatis, stats.WorkSatis, stats.PersonalSatis, stats.Happiness,
		now,
		statID,
	)
	return err
}

func (t *Ticker) markDead(ctx context.Context, tamaID int, now time.Time) error {
	_, err := t.db.ExecContext(ctx,
		`UPDATE Tama SET DeathDay = ?, CauseOfDeath = 'neglect' WHERE TamaId = ?`,
		now, tamaID,
	)
	return err
}

func (t *Ticker) setSickness(ctx context.Context, tamaID int, sicknessName string) error {
	_, err := t.db.ExecContext(ctx,
		`UPDATE Tama SET Sickness = ? WHERE TamaId = ?`,
		sicknessName, tamaID,
	)
	return err
}

func (t *Ticker) loadSicknesses(ctx context.Context) []engine.DBSickness {
	var rows []DBSicknessRow
	if err := t.db.SelectContext(ctx, &rows, "SELECT * FROM Sickness"); err != nil {
		log.Printf("[ticker] error loading sicknesses: %v", err)
		return nil
	}
	result := make([]engine.DBSickness, len(rows))
	for i, r := range rows {
		desc := ""
		if r.Desc != nil {
			desc = *r.Desc
		}
		result[i] = engine.DBSickness{
			ID:             r.SicknessID,
			Name:           r.Name,
			Desc:           desc,
			Type:           r.Type,
			Severity:       r.Severity,
			ExpirationDays: r.ExpirationDays,
			CureCost:       r.CureCost,
			Bonus:          r.Bonus,
			Malus:          r.Malus,
		}
	}
	return result
}

func (t *Ticker) loadEvents(ctx context.Context) []engine.DBEvent {
	var rows []DBEventRow
	if err := t.db.SelectContext(ctx, &rows, "SELECT * FROM Event"); err != nil {
		log.Printf("[ticker] error loading events: %v", err)
		return nil
	}
	result := make([]engine.DBEvent, len(rows))
	for i, r := range rows {
		desc := ""
		if r.Desc != nil {
			desc = *r.Desc
		}
		var minStage *engine.LifeStage
		if r.MinStage != nil {
			stage := engine.LifeStage(*r.MinStage)
			minStage = &stage
		}
		result[i] = engine.DBEvent{
			ID:       r.EventID,
			Name:     r.Name,
			Desc:     desc,
			Severity: r.Severity,
			Scope:    r.Scope,
			MinStage: minStage,
			Bonus:    r.Bonus,
			Malus:    r.Malus,
		}
	}
	return result
}

func (t *Ticker) loadLifeChoices(ctx context.Context) []engine.DBLifeChoice {
	var rows []DBLifeChoiceRow
	if err := t.db.SelectContext(ctx, &rows, "SELECT * FROM LifeChoices"); err != nil {
		log.Printf("[ticker] error loading life choices: %v", err)
		return nil
	}
	result := make([]engine.DBLifeChoice, len(rows))
	for i, r := range rows {
		desc := ""
		if r.Desc != nil {
			desc = *r.Desc
		}
		result[i] = engine.DBLifeChoice{
			ID:         r.LifeChoicesID,
			Name:       r.Name,
			Desc:       desc,
			Stage:      engine.LifeStage(r.Stage),
			Rarity:     r.Rarity,
			ChoiceType: r.ChoiceType,
			Traits:     r.Traits,
			Bonus:      r.Bonus,
			Malus:      r.Malus,
		}
	}
	return result
}

func (t *Ticker) loadChoicesMade(ctx context.Context, tamaID int) map[int]bool {
	type row struct {
		LifeChoicesID int `db:"LifeChoicesId"`
	}
	var rows []row
	if err := t.db.SelectContext(ctx, &rows,
		"SELECT LifeChoicesId FROM TamaLifeChoiceHistory WHERE TamaId = ?", tamaID); err != nil {
		return make(map[int]bool)
	}
	result := make(map[int]bool, len(rows))
	for _, r := range rows {
		result[r.LifeChoicesID] = true
	}
	return result
}

func (t *Ticker) loadFriendContext(ctx context.Context, userID int) *engine.FriendContext {
	type friendRow struct {
		OtherUserID int  `db:"other_user_id"`
		IsDead      bool `db:"is_dead"`
	}
	var rows []friendRow
	err := t.db.SelectContext(ctx, &rows, `
		SELECT
			CASE WHEN f.SenderID = ? THEN f.ReceiverID ELSE f.SenderID END AS other_user_id,
			CASE WHEN t.DeathDay IS NOT NULL THEN 1 ELSE 0 END AS is_dead
		FROM FriendRequests f
		JOIN Tama t ON t.UserId = CASE WHEN f.SenderID = ? THEN f.ReceiverID ELSE f.SenderID END
		WHERE f.Status = 'accepted'
		  AND (f.SenderID = ? OR f.ReceiverID = ?)
	`, userID, userID, userID, userID)
	if err != nil {
		return nil
	}

	fc := &engine.FriendContext{}
	for _, r := range rows {
		if r.IsDead {
			fc.DeadFriends++
		} else {
			fc.AliveFriends++
		}
	}
	return fc
}

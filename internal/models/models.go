package models

import "time"

type User struct {
	UserID             int        `db:"UserId"`
	Name               string     `db:"Name"`
	LastName           string     `db:"LastName"`
	UserName           string     `db:"UserName"`
	Email              string     `db:"Email"`
	PasswordHash       string     `db:"PasswordHash"`
	ClearanceLevel     int        `db:"ClearanceLevel"`
	Verified           bool       `db:"Verified"`
	ProfilPicture      *string    `db:"ProfilPicture"`
	GamingTime         int        `db:"GamingTime"`
	CreationDate       time.Time  `db:"CreationDate"`
	LastConnectionDate *time.Time `db:"LastConnectionDate"`
}

type Race struct {
	RaceID int     `db:"RaceId"`
	Name   string  `db:"Name"`
	Desc   *string `db:"Desc"`
	Bonus  *string `db:"Bonus"`
	Malus  *string `db:"Malus"`
}

type TamaStat struct {
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
}

type Tama struct {
	TamaID       int        `db:"TamaId"`
	UserID       int        `db:"UserId"`
	TamaStatsID  int        `db:"TamaStatsID"`
	Name         string     `db:"Name"`
	Sexe         *bool      `db:"Sexe"`
	Race         string     `db:"Race"`
	Sickness     *string    `db:"Sickness"`
	Birthday     *time.Time `db:"Birthday"`
	DeathDay     *time.Time `db:"DeathDay"`
	CauseOfDeath *string    `db:"CauseOfDeath"`
	Traits       *string    `db:"Traits"`
}

type Friend struct {
	RequestID     int        `db:"RequestId"`
	SenderID      int        `db:"SenderID"`
	ReceiverID    int        `db:"ReceiverID"`
	Status        string     `db:"Status"`
	DateRequested time.Time  `db:"DateRequested"`
	DateResponded *time.Time `db:"DateResponded"`
}

type ChatMessage struct {
	MessageID  int        `db:"MessageId"`
	SenderID   int        `db:"SenderID"`
	ReceiverID int        `db:"ReceiverID"`
	Body       string     `db:"Body"`
	SentAt     time.Time  `db:"SentAt"`
	ReadAt     *time.Time `db:"ReadAt"`
}

type Sponsor struct {
	SponsorID     int       `db:"SponsorId"`
	SponsoredID   int       `db:"SponsoredId"`
	DateOfSponsor time.Time `db:"DateOfSponsor"`
}

type Sickness struct {
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

type Trait struct {
	TraitID  int     `db:"TraitId"`
	Name     string  `db:"Name"`
	Desc     *string `db:"Desc"`
	Category string  `db:"Category"`
	Bonus    *string `db:"Bonus"`
	Malus    *string `db:"Malus"`
}

type Bonus struct {
	BonusID  int     `db:"BonusId"`
	Name     string  `db:"Name"`
	Desc     *string `db:"Desc"`
	Effet    *string `db:"Effet"`
	Duration *int    `db:"Duration"`
}

type Malus struct {
	MalusID  int     `db:"MalusId"`
	Name     string  `db:"Name"`
	Desc     *string `db:"Desc"`
	Effet    *string `db:"Effet"`
	Duration *int    `db:"Duration"`
}

type Event struct {
	EventID  int     `db:"EventId"`
	Name     string  `db:"Name"`
	Desc     *string `db:"Desc"`
	Severity string  `db:"Severity"`
	Scope    string  `db:"Scope"`
	MinStage *string `db:"MinStage"`
	Bonus    *string `db:"Bonus"`
	Malus    *string `db:"Malus"`
}

type LifeChoice struct {
	LifeChoiceID int     `db:"LifeChoicesId"`
	Name         string  `db:"Name"`
	Desc         *string `db:"Desc"`
	Stage        string  `db:"Stage"`
	Rarity       string  `db:"Rarity"`
	ChoiceType   string  `db:"ChoiceType"`
	Traits       *string `db:"Traits"`
	Bonus        *string `db:"Bonus"`
	Malus        *string `db:"Malus"`
}

type ActiveEvent struct {
	ActiveEventID int        `db:"ActiveEventId"`
	EventID       int        `db:"EventId"`
	TargetUserID  *int       `db:"TargetUserId"`
	StartDate     time.Time  `db:"StartDate"`
	EndDate       *time.Time `db:"EndDate"`
	TriggeredBy   *int       `db:"TriggeredBy"`
	IsGlobal      bool       `db:"IsGlobal"`
}

type AdminNotification struct {
	NotificationID int        `db:"NotificationId"`
	TargetUserID   int        `db:"TargetUserId"`
	Type           string     `db:"Type"`
	Payload        string     `db:"Payload"`
	Message        string     `db:"Message"`
	CreatedAt      time.Time  `db:"CreatedAt"`
	ReadAt         *time.Time `db:"ReadAt"`
}

type LifeChoiceOption struct {
	OptionID      int     `db:"OptionId"`
	LifeChoicesID int     `db:"LifeChoicesId"`
	Label         string  `db:"Label"`
	Desc          *string `db:"Desc"`
	Traits        *string `db:"Traits"`
	Bonus         *string `db:"Bonus"`
	Malus         *string `db:"Malus"`
}

type TamaLifeChoiceHistory struct {
	HistoryID      int       `db:"HistoryId"`
	TamaID         int       `db:"TamaId"`
	LifeChoicesID  int       `db:"LifeChoicesId"`
	ChosenOptionID *int      `db:"ChosenOptionId"`
	Action         string    `db:"Action"`
	CreatedAt      time.Time `db:"CreatedAt"`
}

type StatHistory struct {
	HistoryID     int       `db:"HistoryId"`
	TamaID        int       `db:"TamaId"`
	Hunger        int       `db:"Hunger"`
	Boredom       int       `db:"Boredom"`
	Hygiene       int       `db:"Hygiene"`
	Money         int       `db:"Money"`
	SocialSatis   float64   `db:"SocialSatis"`
	WorkSatis     float64   `db:"WorkSatis"`
	PersonalSatis float64   `db:"PersonalSatis"`
	Happiness     float64   `db:"Happiness"`
	Fed           int       `db:"Fed"`
	Played        int       `db:"Played"`
	Cleaned       int       `db:"Cleaned"`
	Worked        int       `db:"Worked"`
	CarAccident   int       `db:"CarAccident"`
	WorkAccident  int       `db:"WorkAccident"`
	Trigger       string    `db:"Trigger"`
	RecordedAt    time.Time `db:"RecordedAt"`
}

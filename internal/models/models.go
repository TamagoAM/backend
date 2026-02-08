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
	TamaStatID    int     `db:"TamaStatId"`
	Food          int     `db:"Food"`
	Play          int     `db:"Play"`
	Cleaned       int     `db:"Cleaned"`
	CarAccident   int     `db:"CarAccident"`
	WorkAccident  int     `db:"WorkAccident"`
	SocialSatis   float64 `db:"SocialSatis"`
	WorkSatis     float64 `db:"WorkSatis"`
	PersonalSatis float64 `db:"PersonalSatis"`
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
	UserID            int       `db:"UserID"`
	FriendID          int       `db:"FriendID"`
	DateBecameFriends time.Time `db:"DateBecameFriends"`
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
	ExpirationDays *int    `db:"ExpirationDays"`
	Bonus          *string `db:"Bonus"`
	Malus          *string `db:"Malus"`
}

type Trait struct {
	TraitID int     `db:"TraitId"`
	Name    string  `db:"Name"`
	Desc    *string `db:"Desc"`
	Bonus   *string `db:"Bonus"`
	Malus   *string `db:"Malus"`
}

type Bonus struct {
	BonusID int     `db:"BonusId"`
	Name    string  `db:"Name"`
	Desc    *string `db:"Desc"`
	Effet   *string `db:"Effet"`
}

type Malus struct {
	MalusID int     `db:"MalusId"`
	Name    string  `db:"Name"`
	Desc    *string `db:"Desc"`
	Effet   *string `db:"Effet"`
}

type Event struct {
	EventID int     `db:"EventId"`
	Name    string  `db:"Name"`
	Desc    *string `db:"Desc"`
	Bonus   *string `db:"Bonus"`
	Malus   *string `db:"Malus"`
}

type LifeChoice struct {
	LifeChoiceID int     `db:"LifeChoicesId"`
	Name         string  `db:"Name"`
	Desc         *string `db:"Desc"`
	Traits       *string `db:"Traits"`
}

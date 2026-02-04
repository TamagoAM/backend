package models

import "time"

type User struct {
	UserID             int        `db:"UserId"`
	Name               string     `db:"Name"`
	LastName           string     `db:"LastName"`
	UserName           string     `db:"UserName"`
	Email              string     `db:"Email"`
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

package graphql

import (
	"strings"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/jmoiron/sqlx"

	"tamagoam/internal/auth"
	"tamagoam/internal/models"
)

type CreateUserInput struct {
	Name               string
	LastName           string
	UserName           string
	Email              string
	PasswordHash       string
	ClearanceLevel     int
	Verified           bool
	ProfilPicture      *string
	GamingTime         int
	LastConnectionDate *time.Time
}

type CreateRaceInput struct {
	Name  string
	Desc  *string
	Bonus *string
	Malus *string
}

type CreateTamaStatInput struct {
	Fed           int
	LastFed       *time.Time
	Played        int
	LastPlayed    *time.Time
	Cleaned       int
	LastCleaned   *time.Time
	Worked        int
	LastWorked    *time.Time
	Hunger        int
	Boredom       int
	Hygiene       int
	Money         int
	CarAccident   int
	WorkAccident  int
	SocialSatis   float64
	WorkSatis     float64
	PersonalSatis float64
}

type CreateTamaInput struct {
	UserID       int
	TamaStatsID  int
	Name         string
	Sexe         *bool
	Race         string
	Sickness     *string
	Birthday     *time.Time
	DeathDay     *time.Time
	CauseOfDeath *string
	Traits       *string
}

type CreateFriendInput struct {
	UserID            int
	FriendID          int
	DateBecameFriends time.Time
}

type CreateSponsorInput struct {
	SponsorID     int
	SponsoredID   int
	DateOfSponsor time.Time
}

type CreateSicknessInput struct {
	Name           string
	Desc           *string
	Type           string
	Severity       string
	ExpirationDays *int
	CureCost       *int
	Bonus          *string
	Malus          *string
}

type CreateTraitInput struct {
	Name     string
	Desc     *string
	Category string
	Bonus    *string
	Malus    *string
}

type CreateBonusInput struct {
	Name     string
	Desc     *string
	Effet    *string
	Duration *int
}

type CreateMalusInput struct {
	Name     string
	Desc     *string
	Effet    *string
	Duration *int
}

type CreateEventInput struct {
	Name     string
	Desc     *string
	Severity string
	Scope    string
	MinStage *string
	Bonus    *string
	Malus    *string
}

type CreateLifeChoiceInput struct {
	Name       string
	Desc       *string
	Stage      string
	Rarity     string
	ChoiceType string
	Traits     *string
	Bonus      *string
	Malus      *string
}

type CreateActiveEventInput struct {
	EventID      int
	TargetUserID *int
	EndDate      *time.Time
	TriggeredBy  *int
	IsGlobal     bool
}

// sourceAs extracts the concrete value from p.Source whether it was passed
// as a value (T) or a pointer (*T).  This is needed because graphql-go
// passes list-element sources as values but single-record resolver returns
// as pointers.
func sourceAs[T any](src interface{}) (T, bool) {
	if v, ok := src.(T); ok {
		return v, true
	}
	if p, ok := src.(*T); ok && p != nil {
		return *p, true
	}
	var zero T
	return zero, false
}

func NewSchema(db *sqlx.DB) (graphql.Schema, error) {
	store := NewSQLStore(db)

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return u.UserID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return u.Name, nil
				}
				return nil, nil
			}},
			"lastName": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return u.LastName, nil
				}
				return nil, nil
			}},
			"userName": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return u.UserName, nil
				}
				return nil, nil
			}},
			"email": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return u.Email, nil
				}
				return nil, nil
			}},
			"clearanceLevel": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return u.ClearanceLevel, nil
				}
				return nil, nil
			}},
			"verified": &graphql.Field{Type: graphql.Boolean, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return u.Verified, nil
				}
				return nil, nil
			}},
			"profilPicture": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return u.ProfilPicture, nil
				}
				return nil, nil
			}},
			"gamingTime": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return u.GamingTime, nil
				}
				return nil, nil
			}},
			"creationDate": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return formatTimeValue(&u.CreationDate), nil
				}
				return nil, nil
			}},
			"lastConnectionDate": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := sourceAs[models.User](p.Source); ok {
					return formatTimeValue(u.LastConnectionDate), nil
				}
				return nil, nil
			}},
		},
	})

	raceType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Race",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := sourceAs[models.Race](p.Source); ok {
					return r.RaceID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := sourceAs[models.Race](p.Source); ok {
					return r.Name, nil
				}
				return nil, nil
			}},
			"desc": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := sourceAs[models.Race](p.Source); ok {
					return r.Desc, nil
				}
				return nil, nil
			}},
			"bonus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := sourceAs[models.Race](p.Source); ok {
					return r.Bonus, nil
				}
				return nil, nil
			}},
			"malus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := sourceAs[models.Race](p.Source); ok {
					return r.Malus, nil
				}
				return nil, nil
			}},
		},
	})

	tamaStatType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TamaStat",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.TamaStatID, nil
				}
				return nil, nil
			}},
			"fed": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.Fed, nil
				}
				return nil, nil
			}},
			"lastFed": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return formatTimeValue(s.LastFed), nil
				}
				return nil, nil
			}},
			"played": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.Played, nil
				}
				return nil, nil
			}},
			"lastPlayed": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return formatTimeValue(s.LastPlayed), nil
				}
				return nil, nil
			}},
			"cleaned": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.Cleaned, nil
				}
				return nil, nil
			}},
			"lastCleaned": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return formatTimeValue(s.LastCleaned), nil
				}
				return nil, nil
			}},
			"worked": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.Worked, nil
				}
				return nil, nil
			}},
			"lastWorked": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return formatTimeValue(s.LastWorked), nil
				}
				return nil, nil
			}},
			"hunger": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.Hunger, nil
				}
				return nil, nil
			}},
			"boredom": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.Boredom, nil
				}
				return nil, nil
			}},
			"hygiene": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.Hygiene, nil
				}
				return nil, nil
			}},
			"money": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.Money, nil
				}
				return nil, nil
			}},
			"carAccident": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.CarAccident, nil
				}
				return nil, nil
			}},
			"workAccident": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.WorkAccident, nil
				}
				return nil, nil
			}},
			"socialSatis": &graphql.Field{Type: graphql.Float, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.SocialSatis, nil
				}
				return nil, nil
			}},
			"workSatis": &graphql.Field{Type: graphql.Float, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.WorkSatis, nil
				}
				return nil, nil
			}},
			"personalSatis": &graphql.Field{Type: graphql.Float, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.TamaStat](p.Source); ok {
					return s.PersonalSatis, nil
				}
				return nil, nil
			}},
		},
	})

	tamaType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Tama",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return t.TamaID, nil
				}
				return nil, nil
			}},
			"userId": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return t.UserID, nil
				}
				return nil, nil
			}},
			"tamaStatsId": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return t.TamaStatsID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return t.Name, nil
				}
				return nil, nil
			}},
			"sexe": &graphql.Field{Type: graphql.Boolean, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return t.Sexe, nil
				}
				return nil, nil
			}},
			"race": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return t.Race, nil
				}
				return nil, nil
			}},
			"sickness": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return t.Sickness, nil
				}
				return nil, nil
			}},
			"birthday": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return formatDateValue(t.Birthday), nil
				}
				return nil, nil
			}},
			"deathDay": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return formatDateValue(t.DeathDay), nil
				}
				return nil, nil
			}},
			"causeOfDeath": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return t.CauseOfDeath, nil
				}
				return nil, nil
			}},
			"traits": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Tama](p.Source); ok {
					return t.Traits, nil
				}
				return nil, nil
			}},
		},
	})

	friendType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Friend",
		Fields: graphql.Fields{
			"userId": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if f, ok := sourceAs[models.Friend](p.Source); ok {
					return f.UserID, nil
				}
				return nil, nil
			}},
			"friendId": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if f, ok := sourceAs[models.Friend](p.Source); ok {
					return f.FriendID, nil
				}
				return nil, nil
			}},
			"dateBecameFriends": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if f, ok := sourceAs[models.Friend](p.Source); ok {
					return formatDateValue(&f.DateBecameFriends), nil
				}
				return nil, nil
			}},
		},
	})

	sponsorType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Sponsor",
		Fields: graphql.Fields{
			"sponsorId": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sponsor](p.Source); ok {
					return s.SponsorID, nil
				}
				return nil, nil
			}},
			"sponsoredId": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sponsor](p.Source); ok {
					return s.SponsoredID, nil
				}
				return nil, nil
			}},
			"dateOfSponsor": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sponsor](p.Source); ok {
					return formatDateValue(&s.DateOfSponsor), nil
				}
				return nil, nil
			}},
		},
	})

	sicknessType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Sickness",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sickness](p.Source); ok {
					return s.SicknessID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sickness](p.Source); ok {
					return s.Name, nil
				}
				return nil, nil
			}},
			"desc": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sickness](p.Source); ok {
					return s.Desc, nil
				}
				return nil, nil
			}},
			"expirationDays": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sickness](p.Source); ok {
					return s.ExpirationDays, nil
				}
				return nil, nil
			}},
			"type": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sickness](p.Source); ok {
					return s.Type, nil
				}
				return nil, nil
			}},
			"severity": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sickness](p.Source); ok {
					return s.Severity, nil
				}
				return nil, nil
			}},
			"cureCost": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sickness](p.Source); ok {
					return s.CureCost, nil
				}
				return nil, nil
			}},
			"bonus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sickness](p.Source); ok {
					return s.Bonus, nil
				}
				return nil, nil
			}},
			"malus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := sourceAs[models.Sickness](p.Source); ok {
					return s.Malus, nil
				}
				return nil, nil
			}},
		},
	})

	traitType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Trait",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Trait](p.Source); ok {
					return t.TraitID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Trait](p.Source); ok {
					return t.Name, nil
				}
				return nil, nil
			}},
			"desc": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Trait](p.Source); ok {
					return t.Desc, nil
				}
				return nil, nil
			}},
			"category": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Trait](p.Source); ok {
					return t.Category, nil
				}
				return nil, nil
			}},
			"bonus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Trait](p.Source); ok {
					return t.Bonus, nil
				}
				return nil, nil
			}},
			"malus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := sourceAs[models.Trait](p.Source); ok {
					return t.Malus, nil
				}
				return nil, nil
			}},
		},
	})

	bonusType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Bonus",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if b, ok := sourceAs[models.Bonus](p.Source); ok {
					return b.BonusID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if b, ok := sourceAs[models.Bonus](p.Source); ok {
					return b.Name, nil
				}
				return nil, nil
			}},
			"desc": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if b, ok := sourceAs[models.Bonus](p.Source); ok {
					return b.Desc, nil
				}
				return nil, nil
			}},
			"effet": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if b, ok := sourceAs[models.Bonus](p.Source); ok {
					return b.Effet, nil
				}
				return nil, nil
			}},
			"duration": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if b, ok := sourceAs[models.Bonus](p.Source); ok {
					return b.Duration, nil
				}
				return nil, nil
			}},
		},
	})

	malusType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Malus",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if m, ok := sourceAs[models.Malus](p.Source); ok {
					return m.MalusID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if m, ok := sourceAs[models.Malus](p.Source); ok {
					return m.Name, nil
				}
				return nil, nil
			}},
			"desc": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if m, ok := sourceAs[models.Malus](p.Source); ok {
					return m.Desc, nil
				}
				return nil, nil
			}},
			"effet": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if m, ok := sourceAs[models.Malus](p.Source); ok {
					return m.Effet, nil
				}
				return nil, nil
			}},
			"duration": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if m, ok := sourceAs[models.Malus](p.Source); ok {
					return m.Duration, nil
				}
				return nil, nil
			}},
		},
	})

	eventType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Event",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := sourceAs[models.Event](p.Source); ok {
					return e.EventID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := sourceAs[models.Event](p.Source); ok {
					return e.Name, nil
				}
				return nil, nil
			}},
			"desc": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := sourceAs[models.Event](p.Source); ok {
					return e.Desc, nil
				}
				return nil, nil
			}},
			"severity": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := sourceAs[models.Event](p.Source); ok {
					return e.Severity, nil
				}
				return nil, nil
			}},
			"scope": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := sourceAs[models.Event](p.Source); ok {
					return e.Scope, nil
				}
				return nil, nil
			}},
			"minStage": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := sourceAs[models.Event](p.Source); ok {
					return e.MinStage, nil
				}
				return nil, nil
			}},
			"bonus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := sourceAs[models.Event](p.Source); ok {
					return e.Bonus, nil
				}
				return nil, nil
			}},
			"malus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := sourceAs[models.Event](p.Source); ok {
					return e.Malus, nil
				}
				return nil, nil
			}},
		},
	})

	lifeChoiceType := graphql.NewObject(graphql.ObjectConfig{
		Name: "LifeChoice",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := sourceAs[models.LifeChoice](p.Source); ok {
					return l.LifeChoiceID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := sourceAs[models.LifeChoice](p.Source); ok {
					return l.Name, nil
				}
				return nil, nil
			}},
			"desc": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := sourceAs[models.LifeChoice](p.Source); ok {
					return l.Desc, nil
				}
				return nil, nil
			}},
			"stage": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := sourceAs[models.LifeChoice](p.Source); ok {
					return l.Stage, nil
				}
				return nil, nil
			}},
			"rarity": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := sourceAs[models.LifeChoice](p.Source); ok {
					return l.Rarity, nil
				}
				return nil, nil
			}},
			"choiceType": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := sourceAs[models.LifeChoice](p.Source); ok {
					return l.ChoiceType, nil
				}
				return nil, nil
			}},
			"traits": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := sourceAs[models.LifeChoice](p.Source); ok {
					return l.Traits, nil
				}
				return nil, nil
			}},
			"bonus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := sourceAs[models.LifeChoice](p.Source); ok {
					return l.Bonus, nil
				}
				return nil, nil
			}},
			"malus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := sourceAs[models.LifeChoice](p.Source); ok {
					return l.Malus, nil
				}
				return nil, nil
			}},
		},
	})

	activeEventType := graphql.NewObject(graphql.ObjectConfig{
		Name: "ActiveEvent",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if a, ok := sourceAs[models.ActiveEvent](p.Source); ok {
					return a.ActiveEventID, nil
				}
				return nil, nil
			}},
			"eventId": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if a, ok := sourceAs[models.ActiveEvent](p.Source); ok {
					return a.EventID, nil
				}
				return nil, nil
			}},
			"targetUserId": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if a, ok := sourceAs[models.ActiveEvent](p.Source); ok {
					return a.TargetUserID, nil
				}
				return nil, nil
			}},
			"startDate": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if a, ok := sourceAs[models.ActiveEvent](p.Source); ok {
					return formatTimeValue(&a.StartDate), nil
				}
				return nil, nil
			}},
			"endDate": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if a, ok := sourceAs[models.ActiveEvent](p.Source); ok {
					return formatTimeValue(a.EndDate), nil
				}
				return nil, nil
			}},
			"triggeredBy": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if a, ok := sourceAs[models.ActiveEvent](p.Source); ok {
					return a.TriggeredBy, nil
				}
				return nil, nil
			}},
			"isGlobal": &graphql.Field{Type: graphql.Boolean, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if a, ok := sourceAs[models.ActiveEvent](p.Source); ok {
					return a.IsGlobal, nil
				}
				return nil, nil
			}},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListUsers(p.Context)
				},
			},
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.GetUser(p.Context, id)
				},
			},
			"races": &graphql.Field{
				Type: graphql.NewList(raceType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListRaces(p.Context)
				},
			},
			"tamaStats": &graphql.Field{
				Type: graphql.NewList(tamaStatType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListTamaStats(p.Context)
				},
			},
			"tamas": &graphql.Field{
				Type: graphql.NewList(tamaType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListTamas(p.Context)
				},
			},
			"friends": &graphql.Field{
				Type: graphql.NewList(friendType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListFriends(p.Context)
				},
			},
			"sponsors": &graphql.Field{
				Type: graphql.NewList(sponsorType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListSponsors(p.Context)
				},
			},
			"sicknesses": &graphql.Field{
				Type: graphql.NewList(sicknessType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListSickness(p.Context)
				},
			},
			"traits": &graphql.Field{
				Type: graphql.NewList(traitType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListTraits(p.Context)
				},
			},
			"bonuses": &graphql.Field{
				Type: graphql.NewList(bonusType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListBonuses(p.Context)
				},
			},
			"maluses": &graphql.Field{
				Type: graphql.NewList(malusType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListMaluses(p.Context)
				},
			},
			"events": &graphql.Field{
				Type: graphql.NewList(eventType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListEvents(p.Context)
				},
			},
			"lifeChoices": &graphql.Field{
				Type: graphql.NewList(lifeChoiceType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListLifeChoices(p.Context)
				},
			},
			"activeEvents": &graphql.Field{
				Type: graphql.NewList(activeEventType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListActiveEvents(p.Context)
				},
			},
			"activeEventsByUser": &graphql.Field{
				Type: graphql.NewList(activeEventType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Args["userId"].(int)
					return store.ActiveEventsByUser(p.Context, userID)
				},
			},
			"globalActiveEvents": &graphql.Field{
				Type: graphql.NewList(activeEventType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.GlobalActiveEvents(p.Context)
				},
			},

			// ─── User-scoped queries for user monitor ─────────────
			"tamasByUser": &graphql.Field{
				Type: graphql.NewList(tamaType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Args["userId"].(int)
					return store.TamasByUser(p.Context, userID)
				},
			},
			"friendsByUser": &graphql.Field{
				Type: graphql.NewList(friendType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Args["userId"].(int)
					return store.FriendsByUser(p.Context, userID)
				},
			},
			"sponsorsByUser": &graphql.Field{
				Type: graphql.NewList(sponsorType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Args["userId"].(int)
					return store.SponsorsByUser(p.Context, userID)
				},
			},
			"sponsoredByUser": &graphql.Field{
				Type: graphql.NewList(sponsorType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Args["userId"].(int)
					return store.SponsoredByUser(p.Context, userID)
				},
			},
			"tamaStatsByUser": &graphql.Field{
				Type: graphql.NewList(tamaStatType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Args["userId"].(int)
					return store.TamaStatsByUser(p.Context, userID)
				},
			},
		},
	})

	createUserInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateUserInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":               &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"lastName":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"userName":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"email":              &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"password":           &graphql.InputObjectFieldConfig{Type: graphql.String},
			"clearanceLevel":     &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"verified":           &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"profilPicture":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"gamingTime":         &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"lastConnectionDate": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	createRaceInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateRaceInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"desc":  &graphql.InputObjectFieldConfig{Type: graphql.String},
			"bonus": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"malus": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	createTamaStatInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateTamaStatInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"fed":           &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"lastFed":       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"played":        &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"lastPlayed":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"cleaned":       &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"lastCleaned":   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"worked":        &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"lastWorked":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"hunger":        &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"boredom":       &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"hygiene":       &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"money":         &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"carAccident":   &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"workAccident":  &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"socialSatis":   &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"workSatis":     &graphql.InputObjectFieldConfig{Type: graphql.Float},
			"personalSatis": &graphql.InputObjectFieldConfig{Type: graphql.Float},
		},
	})

	createTamaInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateTamaInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"userId":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
			"tamaStatsId":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
			"name":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"sexe":         &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
			"race":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"sickness":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"birthday":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"deathDay":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"causeOfDeath": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"traits":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	createFriendInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateFriendInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"userId":            &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
			"friendId":          &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
			"dateBecameFriends": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	createSponsorInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateSponsorInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"sponsorId":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
			"sponsoredId":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
			"dateOfSponsor": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	createSicknessInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateSicknessInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"desc":           &graphql.InputObjectFieldConfig{Type: graphql.String},
			"type":           &graphql.InputObjectFieldConfig{Type: graphql.String},
			"severity":       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"expirationDays": &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"cureCost":       &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"bonus":          &graphql.InputObjectFieldConfig{Type: graphql.String},
			"malus":          &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	createTraitInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateTraitInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"desc":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"category": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"bonus":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"malus":    &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	createBonusInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateBonusInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"desc":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"effet":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"duration": &graphql.InputObjectFieldConfig{Type: graphql.Int},
		},
	})

	createMalusInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateMalusInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"desc":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"effet":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"duration": &graphql.InputObjectFieldConfig{Type: graphql.Int},
		},
	})

	createEventInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateEventInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"desc":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"severity": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"scope":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"minStage": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"bonus":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"malus":    &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	createLifeChoiceInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateLifeChoiceInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"desc":       &graphql.InputObjectFieldConfig{Type: graphql.String},
			"stage":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"rarity":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"choiceType": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"traits":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"bonus":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"malus":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	createActiveEventInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateActiveEventInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"eventId":      &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
			"targetUserId": &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"endDate":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"triggeredBy":  &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"isGlobal":     &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createUserInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateUserInput{
						Name:     inputMap["name"].(string),
						LastName: inputMap["lastName"].(string),
						UserName: inputMap["userName"].(string),
						Email:    inputMap["email"].(string),
					}
					if v, ok := inputMap["password"]; ok {
						if s, ok := v.(string); ok && s != "" {
							hashed, err := auth.HashPassword(s)
							if err != nil {
								return nil, err
							}
							input.PasswordHash = hashed
						}
					}
					if v, ok := inputMap["clearanceLevel"]; ok {
						if i, ok := v.(int); ok {
							input.ClearanceLevel = i
						}
					}
					if v, ok := inputMap["verified"]; ok {
						if b, ok := v.(bool); ok {
							input.Verified = b
						}
					}
					if v, ok := inputMap["profilPicture"]; ok {
						if s, ok := v.(string); ok {
							input.ProfilPicture = &s
						}
					}
					if v, ok := inputMap["gamingTime"]; ok {
						if i, ok := v.(int); ok {
							input.GamingTime = i
						}
					}
					if v, ok := inputMap["lastConnectionDate"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastConnectionDate = &t
							}
						}
					}
					return store.CreateUser(p.Context, input)
				},
			},
			"updateUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createUserInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateUserInput{
						Name:     inputMap["name"].(string),
						LastName: inputMap["lastName"].(string),
						UserName: inputMap["userName"].(string),
						Email:    inputMap["email"].(string),
					}
					if v, ok := inputMap["password"]; ok {
						if s, ok := v.(string); ok && s != "" {
							hashed, err := auth.HashPassword(s)
							if err != nil {
								return nil, err
							}
							input.PasswordHash = hashed
						}
					}
					// If no password provided on update, preserve existing hash
					if input.PasswordHash == "" {
						existing, err := store.GetUser(p.Context, id)
						if err == nil && existing != nil {
							input.PasswordHash = existing.PasswordHash
						}
					}
					if v, ok := inputMap["clearanceLevel"]; ok {
						if i, ok := v.(int); ok {
							input.ClearanceLevel = i
						}
					}
					if v, ok := inputMap["verified"]; ok {
						if b, ok := v.(bool); ok {
							input.Verified = b
						}
					}
					if v, ok := inputMap["profilPicture"]; ok {
						if s, ok := v.(string); ok {
							input.ProfilPicture = &s
						}
					}
					if v, ok := inputMap["gamingTime"]; ok {
						if i, ok := v.(int); ok {
							input.GamingTime = i
						}
					}
					if v, ok := inputMap["lastConnectionDate"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastConnectionDate = &t
							}
						}
					}
					return store.UpdateUser(p.Context, id, input)
				},
			},
			"deleteUser": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteUser(p.Context, id)
				},
			},
			"createRace": &graphql.Field{
				Type: raceType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createRaceInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateRaceInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.CreateRace(p.Context, input)
				},
			},
			"updateRace": &graphql.Field{
				Type: raceType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createRaceInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateRaceInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.UpdateRace(p.Context, id, input)
				},
			},
			"deleteRace": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteRace(p.Context, id)
				},
			},
			"createTamaStat": &graphql.Field{
				Type: tamaStatType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createTamaStatInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateTamaStatInput{}
					if v, ok := inputMap["fed"]; ok {
						if i, ok := v.(int); ok {
							input.Fed = i
						}
					}
					if v, ok := inputMap["lastFed"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastFed = &t
							}
						}
					}
					if v, ok := inputMap["played"]; ok {
						if i, ok := v.(int); ok {
							input.Played = i
						}
					}
					if v, ok := inputMap["lastPlayed"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastPlayed = &t
							}
						}
					}
					if v, ok := inputMap["cleaned"]; ok {
						if i, ok := v.(int); ok {
							input.Cleaned = i
						}
					}
					if v, ok := inputMap["lastCleaned"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastCleaned = &t
							}
						}
					}
					if v, ok := inputMap["worked"]; ok {
						if i, ok := v.(int); ok {
							input.Worked = i
						}
					}
					if v, ok := inputMap["lastWorked"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastWorked = &t
							}
						}
					}
					if v, ok := inputMap["hunger"]; ok {
						if i, ok := v.(int); ok {
							input.Hunger = i
						}
					}
					if v, ok := inputMap["boredom"]; ok {
						if i, ok := v.(int); ok {
							input.Boredom = i
						}
					}
					if v, ok := inputMap["hygiene"]; ok {
						if i, ok := v.(int); ok {
							input.Hygiene = i
						}
					}
					if v, ok := inputMap["money"]; ok {
						if i, ok := v.(int); ok {
							input.Money = i
						}
					}
					if v, ok := inputMap["carAccident"]; ok {
						if i, ok := v.(int); ok {
							input.CarAccident = i
						}
					}
					if v, ok := inputMap["workAccident"]; ok {
						if i, ok := v.(int); ok {
							input.WorkAccident = i
						}
					}
					if v, ok := inputMap["socialSatis"]; ok {
						if f, ok := v.(float64); ok {
							input.SocialSatis = f
						}
					}
					if v, ok := inputMap["workSatis"]; ok {
						if f, ok := v.(float64); ok {
							input.WorkSatis = f
						}
					}
					if v, ok := inputMap["personalSatis"]; ok {
						if f, ok := v.(float64); ok {
							input.PersonalSatis = f
						}
					}
					return store.CreateTamaStat(p.Context, input)
				},
			},
			"updateTamaStat": &graphql.Field{
				Type: tamaStatType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createTamaStatInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateTamaStatInput{}
					if v, ok := inputMap["fed"]; ok {
						if i, ok := v.(int); ok {
							input.Fed = i
						}
					}
					if v, ok := inputMap["lastFed"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastFed = &t
							}
						}
					}
					if v, ok := inputMap["played"]; ok {
						if i, ok := v.(int); ok {
							input.Played = i
						}
					}
					if v, ok := inputMap["lastPlayed"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastPlayed = &t
							}
						}
					}
					if v, ok := inputMap["cleaned"]; ok {
						if i, ok := v.(int); ok {
							input.Cleaned = i
						}
					}
					if v, ok := inputMap["lastCleaned"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastCleaned = &t
							}
						}
					}
					if v, ok := inputMap["worked"]; ok {
						if i, ok := v.(int); ok {
							input.Worked = i
						}
					}
					if v, ok := inputMap["lastWorked"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastWorked = &t
							}
						}
					}
					if v, ok := inputMap["hunger"]; ok {
						if i, ok := v.(int); ok {
							input.Hunger = i
						}
					}
					if v, ok := inputMap["boredom"]; ok {
						if i, ok := v.(int); ok {
							input.Boredom = i
						}
					}
					if v, ok := inputMap["hygiene"]; ok {
						if i, ok := v.(int); ok {
							input.Hygiene = i
						}
					}
					if v, ok := inputMap["money"]; ok {
						if i, ok := v.(int); ok {
							input.Money = i
						}
					}
					if v, ok := inputMap["carAccident"]; ok {
						if i, ok := v.(int); ok {
							input.CarAccident = i
						}
					}
					if v, ok := inputMap["workAccident"]; ok {
						if i, ok := v.(int); ok {
							input.WorkAccident = i
						}
					}
					if v, ok := inputMap["socialSatis"]; ok {
						if f, ok := v.(float64); ok {
							input.SocialSatis = f
						}
					}
					if v, ok := inputMap["workSatis"]; ok {
						if f, ok := v.(float64); ok {
							input.WorkSatis = f
						}
					}
					if v, ok := inputMap["personalSatis"]; ok {
						if f, ok := v.(float64); ok {
							input.PersonalSatis = f
						}
					}
					return store.UpdateTamaStat(p.Context, id, input)
				},
			},
			"deleteTamaStat": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteTamaStat(p.Context, id)
				},
			},
			"createTama": &graphql.Field{
				Type: tamaType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createTamaInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateTamaInput{
						UserID:      inputMap["userId"].(int),
						TamaStatsID: inputMap["tamaStatsId"].(int),
						Name:        inputMap["name"].(string),
						Race:        inputMap["race"].(string),
					}
					if v, ok := inputMap["sexe"]; ok {
						if b, ok := v.(bool); ok {
							input.Sexe = &b
						}
					}
					if v, ok := inputMap["sickness"]; ok {
						if s, ok := v.(string); ok {
							input.Sickness = &s
						}
					}
					if v, ok := inputMap["birthday"]; ok {
						if s, ok := v.(string); ok {
							if t, err := parseDateString(s); err == nil {
								input.Birthday = t
							}
						}
					}
					if v, ok := inputMap["deathDay"]; ok {
						if s, ok := v.(string); ok {
							if t, err := parseDateString(s); err == nil {
								input.DeathDay = t
							}
						}
					}
					if v, ok := inputMap["causeOfDeath"]; ok {
						if s, ok := v.(string); ok {
							input.CauseOfDeath = &s
						}
					}
					if v, ok := inputMap["traits"]; ok {
						if s, ok := v.(string); ok {
							input.Traits = &s
						}
					}
					return store.CreateTama(p.Context, input)
				},
			},
			"updateTama": &graphql.Field{
				Type: tamaType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createTamaInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateTamaInput{
						UserID:      inputMap["userId"].(int),
						TamaStatsID: inputMap["tamaStatsId"].(int),
						Name:        inputMap["name"].(string),
						Race:        inputMap["race"].(string),
					}
					if v, ok := inputMap["sexe"]; ok {
						if b, ok := v.(bool); ok {
							input.Sexe = &b
						}
					}
					if v, ok := inputMap["sickness"]; ok {
						if s, ok := v.(string); ok {
							input.Sickness = &s
						}
					}
					if v, ok := inputMap["birthday"]; ok {
						if s, ok := v.(string); ok {
							if t, err := parseDateString(s); err == nil {
								input.Birthday = t
							}
						}
					}
					if v, ok := inputMap["deathDay"]; ok {
						if s, ok := v.(string); ok {
							if t, err := parseDateString(s); err == nil {
								input.DeathDay = t
							}
						}
					}
					if v, ok := inputMap["causeOfDeath"]; ok {
						if s, ok := v.(string); ok {
							input.CauseOfDeath = &s
						}
					}
					if v, ok := inputMap["traits"]; ok {
						if s, ok := v.(string); ok {
							input.Traits = &s
						}
					}
					return store.UpdateTama(p.Context, id, input)
				},
			},
			"deleteTama": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteTama(p.Context, id)
				},
			},
			"createFriend": &graphql.Field{
				Type: friendType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createFriendInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					date := time.Now()
					if v, ok := inputMap["dateBecameFriends"]; ok {
						if s, ok := v.(string); ok {
							if t, err := parseDateString(s); err == nil && t != nil {
								date = *t
							}
						}
					}
					input := CreateFriendInput{
						UserID:            inputMap["userId"].(int),
						FriendID:          inputMap["friendId"].(int),
						DateBecameFriends: date,
					}
					return store.CreateFriend(p.Context, input)
				},
			},
			"updateFriend": &graphql.Field{
				Type: friendType,
				Args: graphql.FieldConfigArgument{
					"userId":            &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"friendId":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"dateBecameFriends": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userId := p.Args["userId"].(int)
					friendId := p.Args["friendId"].(int)
					date := time.Now()
					if v, ok := p.Args["dateBecameFriends"]; ok {
						if s, ok := v.(string); ok {
							if t, err := parseDateString(s); err == nil && t != nil {
								date = *t
							}
						}
					}
					return store.UpdateFriend(p.Context, userId, friendId, date)
				},
			},
			"deleteFriend": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"userId":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"friendId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userId := p.Args["userId"].(int)
					friendId := p.Args["friendId"].(int)
					return store.DeleteFriend(p.Context, userId, friendId)
				},
			},
			"createSponsor": &graphql.Field{
				Type: sponsorType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createSponsorInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					date := time.Now()
					if v, ok := inputMap["dateOfSponsor"]; ok {
						if s, ok := v.(string); ok {
							if t, err := parseDateString(s); err == nil && t != nil {
								date = *t
							}
						}
					}
					input := CreateSponsorInput{
						SponsorID:     inputMap["sponsorId"].(int),
						SponsoredID:   inputMap["sponsoredId"].(int),
						DateOfSponsor: date,
					}
					return store.CreateSponsor(p.Context, input)
				},
			},
			"updateSponsor": &graphql.Field{
				Type: sponsorType,
				Args: graphql.FieldConfigArgument{
					"sponsorId":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"sponsoredId":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"dateOfSponsor": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					sponsorId := p.Args["sponsorId"].(int)
					sponsoredId := p.Args["sponsoredId"].(int)
					date := time.Now()
					if v, ok := p.Args["dateOfSponsor"]; ok {
						if s, ok := v.(string); ok {
							if t, err := parseDateString(s); err == nil && t != nil {
								date = *t
							}
						}
					}
					return store.UpdateSponsor(p.Context, sponsorId, sponsoredId, date)
				},
			},
			"deleteSponsor": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"sponsorId":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"sponsoredId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					sponsorId := p.Args["sponsorId"].(int)
					sponsoredId := p.Args["sponsoredId"].(int)
					return store.DeleteSponsor(p.Context, sponsorId, sponsoredId)
				},
			},
			"createSickness": &graphql.Field{
				Type: sicknessType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createSicknessInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateSicknessInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["type"]; ok {
						if s, ok := v.(string); ok {
							input.Type = s
						}
					}
					if input.Type == "" {
						input.Type = "acquired"
					}
					if v, ok := inputMap["severity"]; ok {
						if s, ok := v.(string); ok {
							input.Severity = s
						}
					}
					if input.Severity == "" {
						input.Severity = "mild"
					}
					if v, ok := inputMap["expirationDays"]; ok {
						if i, ok := v.(int); ok {
							input.ExpirationDays = &i
						}
					}
					if v, ok := inputMap["cureCost"]; ok {
						if i, ok := v.(int); ok {
							input.CureCost = &i
						}
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.CreateSickness(p.Context, input)
				},
			},
			"updateSickness": &graphql.Field{
				Type: sicknessType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createSicknessInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateSicknessInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["type"]; ok {
						if s, ok := v.(string); ok {
							input.Type = s
						}
					}
					if input.Type == "" {
						input.Type = "acquired"
					}
					if v, ok := inputMap["severity"]; ok {
						if s, ok := v.(string); ok {
							input.Severity = s
						}
					}
					if input.Severity == "" {
						input.Severity = "mild"
					}
					if v, ok := inputMap["expirationDays"]; ok {
						if i, ok := v.(int); ok {
							input.ExpirationDays = &i
						}
					}
					if v, ok := inputMap["cureCost"]; ok {
						if i, ok := v.(int); ok {
							input.CureCost = &i
						}
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.UpdateSickness(p.Context, id, input)
				},
			},
			"deleteSickness": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteSickness(p.Context, id)
				},
			},
			"createTrait": &graphql.Field{
				Type: traitType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createTraitInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateTraitInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["category"]; ok {
						if s, ok := v.(string); ok {
							input.Category = s
						}
					}
					if input.Category == "" {
						input.Category = "positive"
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.CreateTrait(p.Context, input)
				},
			},
			"updateTrait": &graphql.Field{
				Type: traitType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createTraitInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateTraitInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["category"]; ok {
						if s, ok := v.(string); ok {
							input.Category = s
						}
					}
					if input.Category == "" {
						input.Category = "positive"
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.UpdateTrait(p.Context, id, input)
				},
			},
			"deleteTrait": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteTrait(p.Context, id)
				},
			},
			"createBonus": &graphql.Field{
				Type: bonusType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createBonusInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateBonusInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["effet"]; ok {
						if s, ok := v.(string); ok {
							input.Effet = &s
						}
					}
					if v, ok := inputMap["duration"]; ok {
						if i, ok := v.(int); ok {
							input.Duration = &i
						}
					}
					return store.CreateBonus(p.Context, input)
				},
			},
			"updateBonus": &graphql.Field{
				Type: bonusType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createBonusInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateBonusInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["effet"]; ok {
						if s, ok := v.(string); ok {
							input.Effet = &s
						}
					}
					if v, ok := inputMap["duration"]; ok {
						if i, ok := v.(int); ok {
							input.Duration = &i
						}
					}
					return store.UpdateBonus(p.Context, id, input)
				},
			},
			"deleteBonus": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteBonus(p.Context, id)
				},
			},
			"createMalus": &graphql.Field{
				Type: malusType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createMalusInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateMalusInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["effet"]; ok {
						if s, ok := v.(string); ok {
							input.Effet = &s
						}
					}
					if v, ok := inputMap["duration"]; ok {
						if i, ok := v.(int); ok {
							input.Duration = &i
						}
					}
					return store.CreateMalus(p.Context, input)
				},
			},
			"updateMalus": &graphql.Field{
				Type: malusType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createMalusInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateMalusInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["effet"]; ok {
						if s, ok := v.(string); ok {
							input.Effet = &s
						}
					}
					if v, ok := inputMap["duration"]; ok {
						if i, ok := v.(int); ok {
							input.Duration = &i
						}
					}
					return store.UpdateMalus(p.Context, id, input)
				},
			},
			"deleteMalus": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteMalus(p.Context, id)
				},
			},
			"createEvent": &graphql.Field{
				Type: eventType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createEventInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateEventInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["severity"]; ok {
						if s, ok := v.(string); ok {
							input.Severity = s
						}
					}
					if input.Severity == "" {
						input.Severity = "minor"
					}
					if v, ok := inputMap["scope"]; ok {
						if s, ok := v.(string); ok {
							input.Scope = s
						}
					}
					if input.Scope == "" {
						input.Scope = "individual"
					}
					if v, ok := inputMap["minStage"]; ok {
						if s, ok := v.(string); ok {
							input.MinStage = &s
						}
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.CreateEvent(p.Context, input)
				},
			},
			"updateEvent": &graphql.Field{
				Type: eventType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createEventInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateEventInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["severity"]; ok {
						if s, ok := v.(string); ok {
							input.Severity = s
						}
					}
					if input.Severity == "" {
						input.Severity = "minor"
					}
					if v, ok := inputMap["scope"]; ok {
						if s, ok := v.(string); ok {
							input.Scope = s
						}
					}
					if input.Scope == "" {
						input.Scope = "individual"
					}
					if v, ok := inputMap["minStage"]; ok {
						if s, ok := v.(string); ok {
							input.MinStage = &s
						}
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.UpdateEvent(p.Context, id, input)
				},
			},
			"deleteEvent": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteEvent(p.Context, id)
				},
			},
			"createLifeChoice": &graphql.Field{
				Type: lifeChoiceType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createLifeChoiceInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateLifeChoiceInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["stage"]; ok {
						if s, ok := v.(string); ok {
							input.Stage = s
						}
					}
					if input.Stage == "" {
						input.Stage = "childhood"
					}
					if v, ok := inputMap["rarity"]; ok {
						if s, ok := v.(string); ok {
							input.Rarity = s
						}
					}
					if input.Rarity == "" {
						input.Rarity = "common"
					}
					if v, ok := inputMap["choiceType"]; ok {
						if s, ok := v.(string); ok {
							input.ChoiceType = s
						}
					}
					if input.ChoiceType == "" {
						input.ChoiceType = "pool"
					}
					if v, ok := inputMap["traits"]; ok {
						if s, ok := v.(string); ok {
							input.Traits = &s
						}
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.CreateLifeChoice(p.Context, input)
				},
			},
			"updateLifeChoice": &graphql.Field{
				Type: lifeChoiceType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createLifeChoiceInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateLifeChoiceInput{Name: inputMap["name"].(string)}
					if v, ok := inputMap["desc"]; ok {
						if s, ok := v.(string); ok {
							input.Desc = &s
						}
					}
					if v, ok := inputMap["stage"]; ok {
						if s, ok := v.(string); ok {
							input.Stage = s
						}
					}
					if input.Stage == "" {
						input.Stage = "childhood"
					}
					if v, ok := inputMap["rarity"]; ok {
						if s, ok := v.(string); ok {
							input.Rarity = s
						}
					}
					if input.Rarity == "" {
						input.Rarity = "common"
					}
					if v, ok := inputMap["choiceType"]; ok {
						if s, ok := v.(string); ok {
							input.ChoiceType = s
						}
					}
					if input.ChoiceType == "" {
						input.ChoiceType = "pool"
					}
					if v, ok := inputMap["traits"]; ok {
						if s, ok := v.(string); ok {
							input.Traits = &s
						}
					}
					if v, ok := inputMap["bonus"]; ok {
						if s, ok := v.(string); ok {
							input.Bonus = &s
						}
					}
					if v, ok := inputMap["malus"]; ok {
						if s, ok := v.(string); ok {
							input.Malus = &s
						}
					}
					return store.UpdateLifeChoice(p.Context, id, input)
				},
			},
			"deleteLifeChoice": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteLifeChoice(p.Context, id)
				},
			},
			"createActiveEvent": &graphql.Field{
				Type: activeEventType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createActiveEventInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateActiveEventInput{EventID: inputMap["eventId"].(int)}
					if v, ok := inputMap["targetUserId"]; ok {
						if i, ok := v.(int); ok {
							input.TargetUserID = &i
						}
					}
					if v, ok := inputMap["endDate"]; ok {
						if s, ok := v.(string); ok {
							if t, err := parseDateString(s); err == nil && t != nil {
								input.EndDate = t
							}
						}
					}
					if v, ok := inputMap["triggeredBy"]; ok {
						if i, ok := v.(int); ok {
							input.TriggeredBy = &i
						}
					}
					if v, ok := inputMap["isGlobal"]; ok {
						if b, ok := v.(bool); ok {
							input.IsGlobal = b
						}
					}
					return store.CreateActiveEvent(p.Context, input)
				},
			},
			"deleteActiveEvent": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.DeleteActiveEvent(p.Context, id)
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}

func formatTimeValue(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}

func formatDateValue(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}

func parseDateString(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return &t, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return &t, nil
	} else {
		return nil, err
	}
}

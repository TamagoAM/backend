package graphql

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"tamagoam/internal/models"
)

type Store interface {
	ListUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUserName(ctx context.Context, userName string) (*models.User, error)
	CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error)
	UpdateUser(ctx context.Context, id int, input CreateUserInput) (*models.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
	ListRaces(ctx context.Context) ([]models.Race, error)
	GetRace(ctx context.Context, id int) (*models.Race, error)
	CreateRace(ctx context.Context, input CreateRaceInput) (*models.Race, error)
	UpdateRace(ctx context.Context, id int, input CreateRaceInput) (*models.Race, error)
	DeleteRace(ctx context.Context, id int) (bool, error)
	ListTamaStats(ctx context.Context) ([]models.TamaStat, error)
	GetTamaStat(ctx context.Context, id int) (*models.TamaStat, error)
	CreateTamaStat(ctx context.Context, input CreateTamaStatInput) (*models.TamaStat, error)
	UpdateTamaStat(ctx context.Context, id int, input CreateTamaStatInput) (*models.TamaStat, error)
	DeleteTamaStat(ctx context.Context, id int) (bool, error)
	ListTamas(ctx context.Context) ([]models.Tama, error)
	GetTama(ctx context.Context, id int) (*models.Tama, error)
	CreateTama(ctx context.Context, input CreateTamaInput) (*models.Tama, error)
	UpdateTama(ctx context.Context, id int, input CreateTamaInput) (*models.Tama, error)
	DeleteTama(ctx context.Context, id int) (bool, error)
	ListFriends(ctx context.Context) ([]models.Friend, error)
	GetFriend(ctx context.Context, userID int, friendID int) (*models.Friend, error)
	CreateFriend(ctx context.Context, input CreateFriendInput) (*models.Friend, error)
	UpdateFriend(ctx context.Context, userID int, friendID int, dateBecameFriends time.Time) (*models.Friend, error)
	DeleteFriend(ctx context.Context, userID int, friendID int) (bool, error)
	ListSponsors(ctx context.Context) ([]models.Sponsor, error)
	GetSponsor(ctx context.Context, sponsorID int, sponsoredID int) (*models.Sponsor, error)
	CreateSponsor(ctx context.Context, input CreateSponsorInput) (*models.Sponsor, error)
	UpdateSponsor(ctx context.Context, sponsorID int, sponsoredID int, dateOfSponsor time.Time) (*models.Sponsor, error)
	DeleteSponsor(ctx context.Context, sponsorID int, sponsoredID int) (bool, error)
	ListSickness(ctx context.Context) ([]models.Sickness, error)
	GetSickness(ctx context.Context, id int) (*models.Sickness, error)
	CreateSickness(ctx context.Context, input CreateSicknessInput) (*models.Sickness, error)
	UpdateSickness(ctx context.Context, id int, input CreateSicknessInput) (*models.Sickness, error)
	DeleteSickness(ctx context.Context, id int) (bool, error)
	ListTraits(ctx context.Context) ([]models.Trait, error)
	GetTrait(ctx context.Context, id int) (*models.Trait, error)
	CreateTrait(ctx context.Context, input CreateTraitInput) (*models.Trait, error)
	UpdateTrait(ctx context.Context, id int, input CreateTraitInput) (*models.Trait, error)
	DeleteTrait(ctx context.Context, id int) (bool, error)
	ListBonuses(ctx context.Context) ([]models.Bonus, error)
	GetBonus(ctx context.Context, id int) (*models.Bonus, error)
	CreateBonus(ctx context.Context, input CreateBonusInput) (*models.Bonus, error)
	UpdateBonus(ctx context.Context, id int, input CreateBonusInput) (*models.Bonus, error)
	DeleteBonus(ctx context.Context, id int) (bool, error)
	ListMaluses(ctx context.Context) ([]models.Malus, error)
	GetMalus(ctx context.Context, id int) (*models.Malus, error)
	CreateMalus(ctx context.Context, input CreateMalusInput) (*models.Malus, error)
	UpdateMalus(ctx context.Context, id int, input CreateMalusInput) (*models.Malus, error)
	DeleteMalus(ctx context.Context, id int) (bool, error)
	ListEvents(ctx context.Context) ([]models.Event, error)
	GetEvent(ctx context.Context, id int) (*models.Event, error)
	CreateEvent(ctx context.Context, input CreateEventInput) (*models.Event, error)
	UpdateEvent(ctx context.Context, id int, input CreateEventInput) (*models.Event, error)
	DeleteEvent(ctx context.Context, id int) (bool, error)
	ListLifeChoices(ctx context.Context) ([]models.LifeChoice, error)
	GetLifeChoice(ctx context.Context, id int) (*models.LifeChoice, error)
	CreateLifeChoice(ctx context.Context, input CreateLifeChoiceInput) (*models.LifeChoice, error)
	UpdateLifeChoice(ctx context.Context, id int, input CreateLifeChoiceInput) (*models.LifeChoice, error)
	DeleteLifeChoice(ctx context.Context, id int) (bool, error)

	// User-scoped queries for user monitor
	TamasByUser(ctx context.Context, userID int) ([]models.Tama, error)
	FriendsByUser(ctx context.Context, userID int) ([]models.Friend, error)
	SponsorsByUser(ctx context.Context, userID int) ([]models.Sponsor, error)
	SponsoredByUser(ctx context.Context, userID int) ([]models.Sponsor, error)
	TamaStatsByUser(ctx context.Context, userID int) ([]models.TamaStat, error)

	// Update last connection timestamp on login
	UpdateLastConnection(ctx context.Context, userID int) error
}

type SQLStore struct {
	db *sqlx.DB
}

func NewSQLStore(db *sqlx.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) UpdateLastConnection(ctx context.Context, userID int) error {
	_, err := s.db.ExecContext(ctx, `UPDATE Users SET LastConnectionDate = NOW() WHERE UserId = ?`, userID)
	return err
}

func (s *SQLStore) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := s.db.SelectContext(ctx, &users, "SELECT * FROM Users LIMIT 100")
	return users, err
}

func (s *SQLStore) GetUser(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM Users WHERE UserId = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SQLStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM Users WHERE Email = ?", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SQLStore) GetUserByUserName(ctx context.Context, userName string) (*models.User, error) {
	var user models.User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM Users WHERE UserName = ?", userName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SQLStore) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	res, err := s.db.ExecContext(ctx, `INSERT INTO Users (Name, LastName, UserName, Email, PasswordHash, ClearanceLevel, Verified, ProfilPicture, GamingTime, LastConnectionDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, input.Name, input.LastName, input.UserName, input.Email, input.PasswordHash, input.ClearanceLevel, input.Verified, input.ProfilPicture, input.GamingTime, input.LastConnectionDate)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetUser(ctx, int(id))
}

func (s *SQLStore) UpdateUser(ctx context.Context, id int, input CreateUserInput) (*models.User, error) {
	_, err := s.db.ExecContext(ctx, `UPDATE Users SET Name = ?, LastName = ?, UserName = ?, Email = ?, PasswordHash = ?, ClearanceLevel = ?, Verified = ?, ProfilPicture = ?, GamingTime = ?, LastConnectionDate = ? WHERE UserId = ?`, input.Name, input.LastName, input.UserName, input.Email, input.PasswordHash, input.ClearanceLevel, input.Verified, input.ProfilPicture, input.GamingTime, input.LastConnectionDate, id)
	if err != nil {
		return nil, err
	}
	return s.GetUser(ctx, id)
}

func (s *SQLStore) DeleteUser(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Users WHERE UserId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListRaces(ctx context.Context) ([]models.Race, error) {
	var races []models.Race
	err := s.db.SelectContext(ctx, &races, "SELECT * FROM Race ORDER BY Name")
	return races, err
}

func (s *SQLStore) GetRace(ctx context.Context, id int) (*models.Race, error) {
	var race models.Race
	err := s.db.GetContext(ctx, &race, "SELECT * FROM Race WHERE RaceId = ?", id)
	if err != nil {
		return nil, err
	}
	return &race, nil
}

func (s *SQLStore) CreateRace(ctx context.Context, input CreateRaceInput) (*models.Race, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO Race (Name, `Desc`, Bonus, Malus) VALUES (?, ?, ?, ?)", input.Name, input.Desc, input.Bonus, input.Malus)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetRace(ctx, int(id))
}

func (s *SQLStore) UpdateRace(ctx context.Context, id int, input CreateRaceInput) (*models.Race, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE Race SET Name = ?, `Desc` = ?, Bonus = ?, Malus = ? WHERE RaceId = ?", input.Name, input.Desc, input.Bonus, input.Malus, id)
	if err != nil {
		return nil, err
	}
	return s.GetRace(ctx, id)
}

func (s *SQLStore) DeleteRace(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Race WHERE RaceId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListTamaStats(ctx context.Context) ([]models.TamaStat, error) {
	var stats []models.TamaStat
	err := s.db.SelectContext(ctx, &stats, "SELECT * FROM Tama_stats ORDER BY TamaStatId DESC")
	return stats, err
}

func (s *SQLStore) GetTamaStat(ctx context.Context, id int) (*models.TamaStat, error) {
	var stat models.TamaStat
	err := s.db.GetContext(ctx, &stat, "SELECT * FROM Tama_stats WHERE TamaStatId = ?", id)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

func (s *SQLStore) CreateTamaStat(ctx context.Context, input CreateTamaStatInput) (*models.TamaStat, error) {
	res, err := s.db.ExecContext(ctx, `INSERT INTO Tama_stats (Fed, LastFed, Played, LastPlayed, Cleaned, LastCleaned, Worked, LastWorked, Hunger, Boredom, Hygiene, Money, CarAccident, WorkAccident, SocialSatis, WorkSatis, PersonalSatis) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, input.Fed, input.LastFed, input.Played, input.LastPlayed, input.Cleaned, input.LastCleaned, input.Worked, input.LastWorked, input.Hunger, input.Boredom, input.Hygiene, input.Money, input.CarAccident, input.WorkAccident, input.SocialSatis, input.WorkSatis, input.PersonalSatis)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetTamaStat(ctx, int(id))
}

func (s *SQLStore) UpdateTamaStat(ctx context.Context, id int, input CreateTamaStatInput) (*models.TamaStat, error) {
	_, err := s.db.ExecContext(ctx, `UPDATE Tama_stats SET Fed = ?, LastFed = ?, Played = ?, LastPlayed = ?, Cleaned = ?, LastCleaned = ?, Worked = ?, LastWorked = ?, Hunger = ?, Boredom = ?, Hygiene = ?, Money = ?, CarAccident = ?, WorkAccident = ?, SocialSatis = ?, WorkSatis = ?, PersonalSatis = ? WHERE TamaStatId = ?`, input.Fed, input.LastFed, input.Played, input.LastPlayed, input.Cleaned, input.LastCleaned, input.Worked, input.LastWorked, input.Hunger, input.Boredom, input.Hygiene, input.Money, input.CarAccident, input.WorkAccident, input.SocialSatis, input.WorkSatis, input.PersonalSatis, id)
	if err != nil {
		return nil, err
	}
	return s.GetTamaStat(ctx, id)
}

func (s *SQLStore) DeleteTamaStat(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Tama_stats WHERE TamaStatId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListTamas(ctx context.Context) ([]models.Tama, error) {
	var tamas []models.Tama
	err := s.db.SelectContext(ctx, &tamas, "SELECT * FROM Tama ORDER BY TamaId DESC")
	return tamas, err
}

func (s *SQLStore) GetTama(ctx context.Context, id int) (*models.Tama, error) {
	var tama models.Tama
	err := s.db.GetContext(ctx, &tama, "SELECT * FROM Tama WHERE TamaId = ?", id)
	if err != nil {
		return nil, err
	}
	return &tama, nil
}

func (s *SQLStore) CreateTama(ctx context.Context, input CreateTamaInput) (*models.Tama, error) {
	res, err := s.db.ExecContext(ctx, `INSERT INTO Tama (UserId, TamaStatsID, Name, Sexe, Race, Sickness, Birthday, DeathDay, CauseOfDeath, Traits) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, input.UserID, input.TamaStatsID, input.Name, input.Sexe, input.Race, input.Sickness, input.Birthday, input.DeathDay, input.CauseOfDeath, input.Traits)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetTama(ctx, int(id))
}

func (s *SQLStore) UpdateTama(ctx context.Context, id int, input CreateTamaInput) (*models.Tama, error) {
	_, err := s.db.ExecContext(ctx, `UPDATE Tama SET UserId = ?, TamaStatsID = ?, Name = ?, Sexe = ?, Race = ?, Sickness = ?, Birthday = ?, DeathDay = ?, CauseOfDeath = ?, Traits = ? WHERE TamaId = ?`, input.UserID, input.TamaStatsID, input.Name, input.Sexe, input.Race, input.Sickness, input.Birthday, input.DeathDay, input.CauseOfDeath, input.Traits, id)
	if err != nil {
		return nil, err
	}
	return s.GetTama(ctx, id)
}

func (s *SQLStore) DeleteTama(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Tama WHERE TamaId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListFriends(ctx context.Context) ([]models.Friend, error) {
	var friends []models.Friend
	err := s.db.SelectContext(ctx, &friends, "SELECT * FROM Friends ORDER BY DateBecameFriends DESC")
	return friends, err
}

func (s *SQLStore) GetFriend(ctx context.Context, userID int, friendID int) (*models.Friend, error) {
	var friend models.Friend
	err := s.db.GetContext(ctx, &friend, "SELECT * FROM Friends WHERE UserID = ? AND FriendID = ?", userID, friendID)
	if err != nil {
		return nil, err
	}
	return &friend, nil
}

func (s *SQLStore) CreateFriend(ctx context.Context, input CreateFriendInput) (*models.Friend, error) {
	_, err := s.db.ExecContext(ctx, `INSERT INTO Friends (UserID, FriendID, DateBecameFriends) VALUES (?, ?, ?)`, input.UserID, input.FriendID, input.DateBecameFriends)
	if err != nil {
		return nil, err
	}
	return s.GetFriend(ctx, input.UserID, input.FriendID)
}

func (s *SQLStore) UpdateFriend(ctx context.Context, userID int, friendID int, dateBecameFriends time.Time) (*models.Friend, error) {
	_, err := s.db.ExecContext(ctx, `UPDATE Friends SET DateBecameFriends = ? WHERE UserID = ? AND FriendID = ?`, dateBecameFriends, userID, friendID)
	if err != nil {
		return nil, err
	}
	return s.GetFriend(ctx, userID, friendID)
}

func (s *SQLStore) DeleteFriend(ctx context.Context, userID int, friendID int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Friends WHERE UserID = ? AND FriendID = ?`, userID, friendID)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListSponsors(ctx context.Context) ([]models.Sponsor, error) {
	var sponsors []models.Sponsor
	err := s.db.SelectContext(ctx, &sponsors, "SELECT * FROM Sponsor ORDER BY DateOfSponsor DESC")
	return sponsors, err
}

func (s *SQLStore) GetSponsor(ctx context.Context, sponsorID int, sponsoredID int) (*models.Sponsor, error) {
	var sponsor models.Sponsor
	err := s.db.GetContext(ctx, &sponsor, "SELECT * FROM Sponsor WHERE SponsorId = ? AND SponsoredId = ?", sponsorID, sponsoredID)
	if err != nil {
		return nil, err
	}
	return &sponsor, nil
}

func (s *SQLStore) CreateSponsor(ctx context.Context, input CreateSponsorInput) (*models.Sponsor, error) {
	_, err := s.db.ExecContext(ctx, `INSERT INTO Sponsor (SponsorId, SponsoredId, DateOfSponsor) VALUES (?, ?, ?)`, input.SponsorID, input.SponsoredID, input.DateOfSponsor)
	if err != nil {
		return nil, err
	}
	return s.GetSponsor(ctx, input.SponsorID, input.SponsoredID)
}

func (s *SQLStore) UpdateSponsor(ctx context.Context, sponsorID int, sponsoredID int, dateOfSponsor time.Time) (*models.Sponsor, error) {
	_, err := s.db.ExecContext(ctx, `UPDATE Sponsor SET DateOfSponsor = ? WHERE SponsorId = ? AND SponsoredId = ?`, dateOfSponsor, sponsorID, sponsoredID)
	if err != nil {
		return nil, err
	}
	return s.GetSponsor(ctx, sponsorID, sponsoredID)
}

func (s *SQLStore) DeleteSponsor(ctx context.Context, sponsorID int, sponsoredID int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Sponsor WHERE SponsorId = ? AND SponsoredId = ?`, sponsorID, sponsoredID)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListSickness(ctx context.Context) ([]models.Sickness, error) {
	var sicknesses []models.Sickness
	err := s.db.SelectContext(ctx, &sicknesses, "SELECT * FROM Sickness ORDER BY Name")
	return sicknesses, err
}

func (s *SQLStore) GetSickness(ctx context.Context, id int) (*models.Sickness, error) {
	var sickness models.Sickness
	err := s.db.GetContext(ctx, &sickness, "SELECT * FROM Sickness WHERE SicknessId = ?", id)
	if err != nil {
		return nil, err
	}
	return &sickness, nil
}

func (s *SQLStore) CreateSickness(ctx context.Context, input CreateSicknessInput) (*models.Sickness, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO Sickness (Name, `Desc`, ExpirationDays, Bonus, Malus) VALUES (?, ?, ?, ?, ?)", input.Name, input.Desc, input.ExpirationDays, input.Bonus, input.Malus)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetSickness(ctx, int(id))
}

func (s *SQLStore) UpdateSickness(ctx context.Context, id int, input CreateSicknessInput) (*models.Sickness, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE Sickness SET Name = ?, `Desc` = ?, ExpirationDays = ?, Bonus = ?, Malus = ? WHERE SicknessId = ?", input.Name, input.Desc, input.ExpirationDays, input.Bonus, input.Malus, id)
	if err != nil {
		return nil, err
	}
	return s.GetSickness(ctx, id)
}

func (s *SQLStore) DeleteSickness(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Sickness WHERE SicknessId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListTraits(ctx context.Context) ([]models.Trait, error) {
	var traits []models.Trait
	err := s.db.SelectContext(ctx, &traits, "SELECT * FROM Trait ORDER BY Name")
	return traits, err
}

func (s *SQLStore) GetTrait(ctx context.Context, id int) (*models.Trait, error) {
	var trait models.Trait
	err := s.db.GetContext(ctx, &trait, "SELECT * FROM Trait WHERE TraitId = ?", id)
	if err != nil {
		return nil, err
	}
	return &trait, nil
}

func (s *SQLStore) CreateTrait(ctx context.Context, input CreateTraitInput) (*models.Trait, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO Trait (Name, `Desc`, Bonus, Malus) VALUES (?, ?, ?, ?)", input.Name, input.Desc, input.Bonus, input.Malus)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetTrait(ctx, int(id))
}

func (s *SQLStore) UpdateTrait(ctx context.Context, id int, input CreateTraitInput) (*models.Trait, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE Trait SET Name = ?, `Desc` = ?, Bonus = ?, Malus = ? WHERE TraitId = ?", input.Name, input.Desc, input.Bonus, input.Malus, id)
	if err != nil {
		return nil, err
	}
	return s.GetTrait(ctx, id)
}

func (s *SQLStore) DeleteTrait(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Trait WHERE TraitId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListBonuses(ctx context.Context) ([]models.Bonus, error) {
	var bonuses []models.Bonus
	err := s.db.SelectContext(ctx, &bonuses, "SELECT * FROM Bonus ORDER BY Name")
	return bonuses, err
}

func (s *SQLStore) GetBonus(ctx context.Context, id int) (*models.Bonus, error) {
	var bonus models.Bonus
	err := s.db.GetContext(ctx, &bonus, "SELECT * FROM Bonus WHERE BonusId = ?", id)
	if err != nil {
		return nil, err
	}
	return &bonus, nil
}

func (s *SQLStore) CreateBonus(ctx context.Context, input CreateBonusInput) (*models.Bonus, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO Bonus (Name, `Desc`, Effet) VALUES (?, ?, ?)", input.Name, input.Desc, input.Effet)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetBonus(ctx, int(id))
}

func (s *SQLStore) UpdateBonus(ctx context.Context, id int, input CreateBonusInput) (*models.Bonus, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE Bonus SET Name = ?, `Desc` = ?, Effet = ? WHERE BonusId = ?", input.Name, input.Desc, input.Effet, id)
	if err != nil {
		return nil, err
	}
	return s.GetBonus(ctx, id)
}

func (s *SQLStore) DeleteBonus(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Bonus WHERE BonusId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListMaluses(ctx context.Context) ([]models.Malus, error) {
	var maluses []models.Malus
	err := s.db.SelectContext(ctx, &maluses, "SELECT * FROM Malus ORDER BY Name")
	return maluses, err
}

func (s *SQLStore) GetMalus(ctx context.Context, id int) (*models.Malus, error) {
	var malus models.Malus
	err := s.db.GetContext(ctx, &malus, "SELECT * FROM Malus WHERE MalusId = ?", id)
	if err != nil {
		return nil, err
	}
	return &malus, nil
}

func (s *SQLStore) CreateMalus(ctx context.Context, input CreateMalusInput) (*models.Malus, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO Malus (Name, `Desc`, Effet) VALUES (?, ?, ?)", input.Name, input.Desc, input.Effet)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetMalus(ctx, int(id))
}

func (s *SQLStore) UpdateMalus(ctx context.Context, id int, input CreateMalusInput) (*models.Malus, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE Malus SET Name = ?, `Desc` = ?, Effet = ? WHERE MalusId = ?", input.Name, input.Desc, input.Effet, id)
	if err != nil {
		return nil, err
	}
	return s.GetMalus(ctx, id)
}

func (s *SQLStore) DeleteMalus(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Malus WHERE MalusId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListEvents(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	err := s.db.SelectContext(ctx, &events, "SELECT * FROM Event ORDER BY Name")
	return events, err
}

func (s *SQLStore) GetEvent(ctx context.Context, id int) (*models.Event, error) {
	var event models.Event
	err := s.db.GetContext(ctx, &event, "SELECT * FROM Event WHERE EventId = ?", id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *SQLStore) CreateEvent(ctx context.Context, input CreateEventInput) (*models.Event, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO Event (Name, `Desc`, Bonus, Malus) VALUES (?, ?, ?, ?)", input.Name, input.Desc, input.Bonus, input.Malus)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetEvent(ctx, int(id))
}

func (s *SQLStore) UpdateEvent(ctx context.Context, id int, input CreateEventInput) (*models.Event, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE Event SET Name = ?, `Desc` = ?, Bonus = ?, Malus = ? WHERE EventId = ?", input.Name, input.Desc, input.Bonus, input.Malus, id)
	if err != nil {
		return nil, err
	}
	return s.GetEvent(ctx, id)
}

func (s *SQLStore) DeleteEvent(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Event WHERE EventId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ListLifeChoices(ctx context.Context) ([]models.LifeChoice, error) {
	var choices []models.LifeChoice
	err := s.db.SelectContext(ctx, &choices, "SELECT * FROM LifeChoices ORDER BY Name")
	return choices, err
}

func (s *SQLStore) GetLifeChoice(ctx context.Context, id int) (*models.LifeChoice, error) {
	var choice models.LifeChoice
	err := s.db.GetContext(ctx, &choice, "SELECT * FROM LifeChoices WHERE LifeChoicesId = ?", id)
	if err != nil {
		return nil, err
	}
	return &choice, nil
}

func (s *SQLStore) CreateLifeChoice(ctx context.Context, input CreateLifeChoiceInput) (*models.LifeChoice, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO LifeChoices (Name, `Desc`, Traits) VALUES (?, ?, ?)", input.Name, input.Desc, input.Traits)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetLifeChoice(ctx, int(id))
}

func (s *SQLStore) UpdateLifeChoice(ctx context.Context, id int, input CreateLifeChoiceInput) (*models.LifeChoice, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE LifeChoices SET Name = ?, `Desc` = ?, Traits = ? WHERE LifeChoicesId = ?", input.Name, input.Desc, input.Traits, id)
	if err != nil {
		return nil, err
	}
	return s.GetLifeChoice(ctx, id)
}

func (s *SQLStore) DeleteLifeChoice(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM LifeChoices WHERE LifeChoicesId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

// ─── User-scoped queries for user monitor ────────────────────────

func (s *SQLStore) TamasByUser(ctx context.Context, userID int) ([]models.Tama, error) {
	var tamas []models.Tama
	err := s.db.SelectContext(ctx, &tamas, "SELECT * FROM Tama WHERE UserId = ? ORDER BY TamaId DESC", userID)
	return tamas, err
}

func (s *SQLStore) FriendsByUser(ctx context.Context, userID int) ([]models.Friend, error) {
	var friends []models.Friend
	err := s.db.SelectContext(ctx, &friends, "SELECT * FROM Friends WHERE UserID = ? OR FriendID = ? ORDER BY DateBecameFriends DESC", userID, userID)
	return friends, err
}

func (s *SQLStore) SponsorsByUser(ctx context.Context, userID int) ([]models.Sponsor, error) {
	var sponsors []models.Sponsor
	err := s.db.SelectContext(ctx, &sponsors, "SELECT * FROM Sponsor WHERE SponsorId = ? ORDER BY DateOfSponsor DESC", userID)
	return sponsors, err
}

func (s *SQLStore) SponsoredByUser(ctx context.Context, userID int) ([]models.Sponsor, error) {
	var sponsors []models.Sponsor
	err := s.db.SelectContext(ctx, &sponsors, "SELECT * FROM Sponsor WHERE SponsoredId = ? ORDER BY DateOfSponsor DESC", userID)
	return sponsors, err
}

func (s *SQLStore) TamaStatsByUser(ctx context.Context, userID int) ([]models.TamaStat, error) {
	var stats []models.TamaStat
	err := s.db.SelectContext(ctx, &stats, "SELECT ts.* FROM Tama_stats ts INNER JOIN Tama t ON t.TamaStatsID = ts.TamaStatId WHERE t.UserId = ? ORDER BY ts.TamaStatId DESC", userID)
	return stats, err
}

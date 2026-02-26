package graphql

import (
	"context"
	"fmt"
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
	GetFriendRequest(ctx context.Context, id int) (*models.Friend, error)
	SendFriendRequest(ctx context.Context, senderID int, receiverID int) (*models.Friend, error)
	RespondFriendRequest(ctx context.Context, requestID int, accept bool) (*models.Friend, error)
	DeleteFriend(ctx context.Context, requestID int) (bool, error)
	AcceptedFriendsByUser(ctx context.Context, userID int) ([]models.Friend, error)
	PendingRequestsForUser(ctx context.Context, userID int) ([]models.Friend, error)
	SentRequestsByUser(ctx context.Context, userID int) ([]models.Friend, error)
	AcceptedFriendCount(ctx context.Context, userID int) (int, error)
	SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error)
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
	ListActiveEvents(ctx context.Context) ([]models.ActiveEvent, error)
	GetActiveEvent(ctx context.Context, id int) (*models.ActiveEvent, error)
	CreateActiveEvent(ctx context.Context, input CreateActiveEventInput) (*models.ActiveEvent, error)
	DeleteActiveEvent(ctx context.Context, id int) (bool, error)
	ActiveEventsByUser(ctx context.Context, userID int) ([]models.ActiveEvent, error)
	GlobalActiveEvents(ctx context.Context) ([]models.ActiveEvent, error)

	// User-scoped queries for user monitor
	TamasByUser(ctx context.Context, userID int) ([]models.Tama, error)
	FriendsByUser(ctx context.Context, userID int) ([]models.Friend, error) // alias for AcceptedFriendsByUser
	SponsorsByUser(ctx context.Context, userID int) ([]models.Sponsor, error)
	SponsoredByUser(ctx context.Context, userID int) ([]models.Sponsor, error)
	TamaStatsByUser(ctx context.Context, userID int) ([]models.TamaStat, error)

	// Update last connection timestamp on login
	UpdateLastConnection(ctx context.Context, userID int) error

	// Stat history
	CreateStatHistory(ctx context.Context, input CreateStatHistoryInput) (*models.StatHistory, error)
	StatHistoryByTama(ctx context.Context, tamaID int, since *time.Time) ([]models.StatHistory, error)

	// Night cycle & notifications
	SetLightsOff(ctx context.Context, tamaStatID int, lightsOff bool) error
	UpdateTimezone(ctx context.Context, userID int, timezone string) error
	RegisterPushToken(ctx context.Context, userID int, token string, platform string) error
	UnregisterPushToken(ctx context.Context, userID int, token string) error

	// Store & payment
	ListStoreItems(ctx context.Context) ([]models.StoreItem, error)
	GetStoreItem(ctx context.Context, id int) (*models.StoreItem, error)
	CreatePayment(ctx context.Context, userID, itemID, amount int, currency string) (*models.Payment, error)
	GetPayment(ctx context.Context, id int) (*models.Payment, error)
	PaymentsByUser(ctx context.Context, userID int) ([]models.Payment, error)
	UserInventoryByUser(ctx context.Context, userID int) ([]models.UserInventory, error)

	// Diamond currency
	AddDiamonds(ctx context.Context, userID int, amount int) error
	SpendDiamonds(ctx context.Context, userID int, amount int) error
	GetDiamonds(ctx context.Context, userID int) (int, error)

	// Inventory management
	AddToInventory(ctx context.Context, userID int, itemID int) (*models.UserInventory, error)
	UseInventoryItem(ctx context.Context, userID int, itemID int) error
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
	res, err := s.db.ExecContext(ctx, `INSERT INTO Tama_stats (Fed, LastFed, Played, LastPlayed, Cleaned, LastCleaned, Worked, LastWorked, Hunger, Boredom, Hygiene, Money, CarAccident, WorkAccident, SocialSatis, WorkSatis, PersonalSatis, Happiness) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, input.Fed, input.LastFed, input.Played, input.LastPlayed, input.Cleaned, input.LastCleaned, input.Worked, input.LastWorked, input.Hunger, input.Boredom, input.Hygiene, input.Money, input.CarAccident, input.WorkAccident, input.SocialSatis, input.WorkSatis, input.PersonalSatis, input.Happiness)
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
	_, err := s.db.ExecContext(ctx, `UPDATE Tama_stats SET Fed = ?, LastFed = ?, Played = ?, LastPlayed = ?, Cleaned = ?, LastCleaned = ?, Worked = ?, LastWorked = ?, Hunger = ?, Boredom = ?, Hygiene = ?, Money = ?, CarAccident = ?, WorkAccident = ?, SocialSatis = ?, WorkSatis = ?, PersonalSatis = ?, Happiness = ? WHERE TamaStatId = ?`, input.Fed, input.LastFed, input.Played, input.LastPlayed, input.Cleaned, input.LastCleaned, input.Worked, input.LastWorked, input.Hunger, input.Boredom, input.Hygiene, input.Money, input.CarAccident, input.WorkAccident, input.SocialSatis, input.WorkSatis, input.PersonalSatis, input.Happiness, id)
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
	err := s.db.SelectContext(ctx, &friends, "SELECT * FROM Friends ORDER BY DateRequested DESC")
	return friends, err
}

func (s *SQLStore) GetFriendRequest(ctx context.Context, id int) (*models.Friend, error) {
	var friend models.Friend
	err := s.db.GetContext(ctx, &friend, "SELECT * FROM Friends WHERE RequestId = ?", id)
	if err != nil {
		return nil, err
	}
	return &friend, nil
}

func (s *SQLStore) SendFriendRequest(ctx context.Context, senderID int, receiverID int) (*models.Friend, error) {
	res, err := s.db.ExecContext(ctx, `INSERT INTO Friends (SenderID, ReceiverID, Status) VALUES (?, ?, 'pending')`, senderID, receiverID)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetFriendRequest(ctx, int(id))
}

func (s *SQLStore) RespondFriendRequest(ctx context.Context, requestID int, accept bool) (*models.Friend, error) {
	status := "declined"
	if accept {
		status = "accepted"
	}
	_, err := s.db.ExecContext(ctx, `UPDATE Friends SET Status = ?, DateResponded = NOW() WHERE RequestId = ? AND Status = 'pending'`, status, requestID)
	if err != nil {
		return nil, err
	}
	return s.GetFriendRequest(ctx, requestID)
}

func (s *SQLStore) DeleteFriend(ctx context.Context, requestID int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM Friends WHERE RequestId = ?`, requestID)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) AcceptedFriendsByUser(ctx context.Context, userID int) ([]models.Friend, error) {
	var friends []models.Friend
	err := s.db.SelectContext(ctx, &friends, "SELECT * FROM Friends WHERE Status = 'accepted' AND (SenderID = ? OR ReceiverID = ?) ORDER BY DateResponded DESC", userID, userID)
	return friends, err
}

func (s *SQLStore) PendingRequestsForUser(ctx context.Context, userID int) ([]models.Friend, error) {
	var friends []models.Friend
	err := s.db.SelectContext(ctx, &friends, "SELECT * FROM Friends WHERE Status = 'pending' AND ReceiverID = ? ORDER BY DateRequested DESC", userID)
	return friends, err
}

func (s *SQLStore) SentRequestsByUser(ctx context.Context, userID int) ([]models.Friend, error) {
	var friends []models.Friend
	err := s.db.SelectContext(ctx, &friends, "SELECT * FROM Friends WHERE Status = 'pending' AND SenderID = ? ORDER BY DateRequested DESC", userID)
	return friends, err
}

func (s *SQLStore) AcceptedFriendCount(ctx context.Context, userID int) (int, error) {
	var count int
	err := s.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM Friends WHERE Status = 'accepted' AND (SenderID = ? OR ReceiverID = ?)", userID, userID)
	return count, err
}

func (s *SQLStore) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	var users []models.User
	pattern := "%" + query + "%"
	err := s.db.SelectContext(ctx, &users, "SELECT * FROM Users WHERE UserName LIKE ? OR Name LIKE ? OR LastName LIKE ? ORDER BY UserName LIMIT ?", pattern, pattern, pattern, limit)
	return users, err
}

func (s *SQLStore) FriendsByUser(ctx context.Context, userID int) ([]models.Friend, error) {
	return s.AcceptedFriendsByUser(ctx, userID)
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
	res, err := s.db.ExecContext(ctx, "INSERT INTO Sickness (Name, `Desc`, Type, Severity, ExpirationDays, CureCost, Bonus, Malus) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", input.Name, input.Desc, input.Type, input.Severity, input.ExpirationDays, input.CureCost, input.Bonus, input.Malus)
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
	_, err := s.db.ExecContext(ctx, "UPDATE Sickness SET Name = ?, `Desc` = ?, Type = ?, Severity = ?, ExpirationDays = ?, CureCost = ?, Bonus = ?, Malus = ? WHERE SicknessId = ?", input.Name, input.Desc, input.Type, input.Severity, input.ExpirationDays, input.CureCost, input.Bonus, input.Malus, id)
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
	res, err := s.db.ExecContext(ctx, "INSERT INTO Trait (Name, `Desc`, Category, Bonus, Malus) VALUES (?, ?, ?, ?, ?)", input.Name, input.Desc, input.Category, input.Bonus, input.Malus)
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
	_, err := s.db.ExecContext(ctx, "UPDATE Trait SET Name = ?, `Desc` = ?, Category = ?, Bonus = ?, Malus = ? WHERE TraitId = ?", input.Name, input.Desc, input.Category, input.Bonus, input.Malus, id)
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
	res, err := s.db.ExecContext(ctx, "INSERT INTO Bonus (Name, `Desc`, Effet, Duration) VALUES (?, ?, ?, ?)", input.Name, input.Desc, input.Effet, input.Duration)
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
	_, err := s.db.ExecContext(ctx, "UPDATE Bonus SET Name = ?, `Desc` = ?, Effet = ?, Duration = ? WHERE BonusId = ?", input.Name, input.Desc, input.Effet, input.Duration, id)
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
	res, err := s.db.ExecContext(ctx, "INSERT INTO Malus (Name, `Desc`, Effet, Duration) VALUES (?, ?, ?, ?)", input.Name, input.Desc, input.Effet, input.Duration)
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
	_, err := s.db.ExecContext(ctx, "UPDATE Malus SET Name = ?, `Desc` = ?, Effet = ?, Duration = ? WHERE MalusId = ?", input.Name, input.Desc, input.Effet, input.Duration, id)
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
	res, err := s.db.ExecContext(ctx, "INSERT INTO Event (Name, `Desc`, Severity, Scope, MinStage, Bonus, Malus) VALUES (?, ?, ?, ?, ?, ?, ?)", input.Name, input.Desc, input.Severity, input.Scope, input.MinStage, input.Bonus, input.Malus)
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
	_, err := s.db.ExecContext(ctx, "UPDATE Event SET Name = ?, `Desc` = ?, Severity = ?, Scope = ?, MinStage = ?, Bonus = ?, Malus = ? WHERE EventId = ?", input.Name, input.Desc, input.Severity, input.Scope, input.MinStage, input.Bonus, input.Malus, id)
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
	res, err := s.db.ExecContext(ctx, "INSERT INTO LifeChoices (Name, `Desc`, Stage, Rarity, ChoiceType, Traits, Bonus, Malus) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", input.Name, input.Desc, input.Stage, input.Rarity, input.ChoiceType, input.Traits, input.Bonus, input.Malus)
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
	_, err := s.db.ExecContext(ctx, "UPDATE LifeChoices SET Name = ?, `Desc` = ?, Stage = ?, Rarity = ?, ChoiceType = ?, Traits = ?, Bonus = ?, Malus = ? WHERE LifeChoicesId = ?", input.Name, input.Desc, input.Stage, input.Rarity, input.ChoiceType, input.Traits, input.Bonus, input.Malus, id)
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

func (s *SQLStore) FriendsOldByUser(ctx context.Context, userID int) ([]models.Friend, error) {
	// Deprecated: kept for compilation; use AcceptedFriendsByUser
	return s.AcceptedFriendsByUser(ctx, userID)
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

// ─── ActiveEvent CRUD ─────────────────────────────────

func (s *SQLStore) ListActiveEvents(ctx context.Context) ([]models.ActiveEvent, error) {
	var events []models.ActiveEvent
	err := s.db.SelectContext(ctx, &events, "SELECT * FROM ActiveEvent ORDER BY StartDate DESC")
	return events, err
}

func (s *SQLStore) GetActiveEvent(ctx context.Context, id int) (*models.ActiveEvent, error) {
	var event models.ActiveEvent
	err := s.db.GetContext(ctx, &event, "SELECT * FROM ActiveEvent WHERE ActiveEventId = ?", id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *SQLStore) CreateActiveEvent(ctx context.Context, input CreateActiveEventInput) (*models.ActiveEvent, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO ActiveEvent (EventId, TargetUserId, EndDate, TriggeredBy, IsGlobal) VALUES (?, ?, ?, ?, ?)", input.EventID, input.TargetUserID, input.EndDate, input.TriggeredBy, input.IsGlobal)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetActiveEvent(ctx, int(id))
}

func (s *SQLStore) DeleteActiveEvent(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM ActiveEvent WHERE ActiveEventId = ?`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *SQLStore) ActiveEventsByUser(ctx context.Context, userID int) ([]models.ActiveEvent, error) {
	var events []models.ActiveEvent
	err := s.db.SelectContext(ctx, &events, "SELECT * FROM ActiveEvent WHERE TargetUserId = ? OR IsGlobal = TRUE ORDER BY StartDate DESC", userID)
	return events, err
}

func (s *SQLStore) GlobalActiveEvents(ctx context.Context) ([]models.ActiveEvent, error) {
	var events []models.ActiveEvent
	err := s.db.SelectContext(ctx, &events, "SELECT * FROM ActiveEvent WHERE IsGlobal = TRUE ORDER BY StartDate DESC")
	return events, err
}

// ─── LifeChoiceOption CRUD ─────────────────────────────

func (s *SQLStore) ListOptionsByChoice(ctx context.Context, lifeChoicesID int) ([]models.LifeChoiceOption, error) {
	var opts []models.LifeChoiceOption
	err := s.db.SelectContext(ctx, &opts, "SELECT * FROM LifeChoiceOption WHERE LifeChoicesId = ? ORDER BY OptionId", lifeChoicesID)
	return opts, err
}

func (s *SQLStore) ListAllOptions(ctx context.Context) ([]models.LifeChoiceOption, error) {
	var opts []models.LifeChoiceOption
	err := s.db.SelectContext(ctx, &opts, "SELECT * FROM LifeChoiceOption ORDER BY LifeChoicesId, OptionId")
	return opts, err
}

func (s *SQLStore) GetOption(ctx context.Context, id int) (*models.LifeChoiceOption, error) {
	var opt models.LifeChoiceOption
	err := s.db.GetContext(ctx, &opt, "SELECT * FROM LifeChoiceOption WHERE OptionId = ?", id)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}

func (s *SQLStore) CreateOption(ctx context.Context, input CreateLifeChoiceOptionInput) (*models.LifeChoiceOption, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO LifeChoiceOption (LifeChoicesId, Label, `Desc`, Traits, Bonus, Malus) VALUES (?, ?, ?, ?, ?, ?)", input.LifeChoicesID, input.Label, input.Desc, input.Traits, input.Bonus, input.Malus)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetOption(ctx, int(id))
}

func (s *SQLStore) UpdateOption(ctx context.Context, id int, input CreateLifeChoiceOptionInput) (*models.LifeChoiceOption, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE LifeChoiceOption SET LifeChoicesId = ?, Label = ?, `Desc` = ?, Traits = ?, Bonus = ?, Malus = ? WHERE OptionId = ?", input.LifeChoicesID, input.Label, input.Desc, input.Traits, input.Bonus, input.Malus, id)
	if err != nil {
		return nil, err
	}
	return s.GetOption(ctx, id)
}

func (s *SQLStore) DeleteOption(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, "DELETE FROM LifeChoiceOption WHERE OptionId = ?", id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

// ─── TamaLifeChoiceHistory CRUD ────────────────────────

func (s *SQLStore) CreateHistory(ctx context.Context, input CreateLifeChoiceHistoryInput) (*models.TamaLifeChoiceHistory, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO TamaLifeChoiceHistory (TamaId, LifeChoicesId, ChosenOptionId, Action) VALUES (?, ?, ?, ?)", input.TamaID, input.LifeChoicesID, input.ChosenOptionID, input.Action)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	var h models.TamaLifeChoiceHistory
	err = s.db.GetContext(ctx, &h, "SELECT * FROM TamaLifeChoiceHistory WHERE HistoryId = ?", id)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (s *SQLStore) ListHistoryByTama(ctx context.Context, tamaID int) ([]models.TamaLifeChoiceHistory, error) {
	var hist []models.TamaLifeChoiceHistory
	err := s.db.SelectContext(ctx, &hist, "SELECT * FROM TamaLifeChoiceHistory WHERE TamaId = ? ORDER BY CreatedAt DESC", tamaID)
	return hist, err
}

// ─── StatHistory CRUD ──────────────────────────────

func (s *SQLStore) CreateStatHistory(ctx context.Context, input CreateStatHistoryInput) (*models.StatHistory, error) {
	res, err := s.db.ExecContext(ctx, `INSERT INTO StatHistory (TamaId, Hunger, Boredom, Hygiene, Money, SocialSatis, WorkSatis, PersonalSatis, Happiness, Fed, Played, Cleaned, Worked, CarAccident, WorkAccident, `+"`Trigger`"+`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.TamaID, input.Hunger, input.Boredom, input.Hygiene, input.Money,
		input.SocialSatis, input.WorkSatis, input.PersonalSatis, input.Happiness,
		input.Fed, input.Played, input.Cleaned, input.Worked,
		input.CarAccident, input.WorkAccident, input.Trigger)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	var h models.StatHistory
	err = s.db.GetContext(ctx, &h, "SELECT * FROM StatHistory WHERE HistoryId = ?", id)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (s *SQLStore) StatHistoryByTama(ctx context.Context, tamaID int, since *time.Time) ([]models.StatHistory, error) {
	var hist []models.StatHistory
	if since != nil {
		err := s.db.SelectContext(ctx, &hist, "SELECT * FROM StatHistory WHERE TamaId = ? AND RecordedAt >= ? ORDER BY RecordedAt ASC", tamaID, since)
		return hist, err
	}
	err := s.db.SelectContext(ctx, &hist, "SELECT * FROM StatHistory WHERE TamaId = ? ORDER BY RecordedAt ASC", tamaID)
	return hist, err
}

// ─── Night cycle & Notifications ──────────────────

func (s *SQLStore) SetLightsOff(ctx context.Context, tamaStatID int, lightsOff bool) error {
	if lightsOff {
		_, err := s.db.ExecContext(ctx, "UPDATE Tama_stats SET LightsOff = TRUE, LightsOffAt = NOW() WHERE TamaStatId = ?", tamaStatID)
		return err
	}
	_, err := s.db.ExecContext(ctx, "UPDATE Tama_stats SET LightsOff = FALSE, LightsOffAt = NULL WHERE TamaStatId = ?", tamaStatID)
	return err
}

func (s *SQLStore) UpdateTimezone(ctx context.Context, userID int, timezone string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE Users SET Timezone = ? WHERE UserId = ?", timezone, userID)
	return err
}

func (s *SQLStore) RegisterPushToken(ctx context.Context, userID int, token string, platform string) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO PushToken (UserId, Token, Platform) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE Platform = VALUES(Platform), CreatedAt = NOW()",
		userID, token, platform)
	return err
}

func (s *SQLStore) UnregisterPushToken(ctx context.Context, userID int, token string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM PushToken WHERE UserId = ? AND Token = ?", userID, token)
	return err
}

// ─── Store & Payment methods ─────────────────────

func (s *SQLStore) ListStoreItems(ctx context.Context) ([]models.StoreItem, error) {
	var items []models.StoreItem
	err := s.db.SelectContext(ctx, &items, "SELECT * FROM StoreItem WHERE Active = TRUE ORDER BY Category, Price")
	return items, err
}

func (s *SQLStore) GetStoreItem(ctx context.Context, id int) (*models.StoreItem, error) {
	var item models.StoreItem
	err := s.db.GetContext(ctx, &item, "SELECT * FROM StoreItem WHERE ItemId = ?", id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *SQLStore) CreatePayment(ctx context.Context, userID, itemID, amount int, currency string) (*models.Payment, error) {
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO Payment (UserId, ItemId, Amount, Currency, Status) VALUES (?, ?, ?, ?, 'pending')",
		userID, itemID, amount, currency)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	var payment models.Payment
	err = s.db.GetContext(ctx, &payment, "SELECT * FROM Payment WHERE PaymentId = ?", id)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *SQLStore) GetPayment(ctx context.Context, id int) (*models.Payment, error) {
	var payment models.Payment
	err := s.db.GetContext(ctx, &payment, "SELECT * FROM Payment WHERE PaymentId = ?", id)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *SQLStore) PaymentsByUser(ctx context.Context, userID int) ([]models.Payment, error) {
	var payments []models.Payment
	err := s.db.SelectContext(ctx, &payments, "SELECT * FROM Payment WHERE UserId = ? ORDER BY CreatedAt DESC", userID)
	return payments, err
}

func (s *SQLStore) UserInventoryByUser(ctx context.Context, userID int) ([]models.UserInventory, error) {
	var items []models.UserInventory
	err := s.db.SelectContext(ctx, &items, "SELECT * FROM UserInventory WHERE UserId = ?", userID)
	return items, err
}

// ─── Diamond currency methods ────────────────────

func (s *SQLStore) AddDiamonds(ctx context.Context, userID int, amount int) error {
	_, err := s.db.ExecContext(ctx, "UPDATE Users SET Diamonds = Diamonds + ? WHERE UserId = ?", amount, userID)
	return err
}

func (s *SQLStore) SpendDiamonds(ctx context.Context, userID int, amount int) error {
	res, err := s.db.ExecContext(ctx, "UPDATE Users SET Diamonds = Diamonds - ? WHERE UserId = ? AND Diamonds >= ?", amount, userID, amount)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("not enough diamonds")
	}
	return nil
}

func (s *SQLStore) GetDiamonds(ctx context.Context, userID int) (int, error) {
	var diamonds int
	err := s.db.GetContext(ctx, &diamonds, "SELECT Diamonds FROM Users WHERE UserId = ?", userID)
	return diamonds, err
}

// ─── Inventory management ────────────────────────

func (s *SQLStore) AddToInventory(ctx context.Context, userID int, itemID int) (*models.UserInventory, error) {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO UserInventory (UserId, ItemId, Quantity) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE Quantity = Quantity + 1",
		userID, itemID)
	if err != nil {
		return nil, err
	}
	var inv models.UserInventory
	err = s.db.GetContext(ctx, &inv, "SELECT * FROM UserInventory WHERE UserId = ? AND ItemId = ?", userID, itemID)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (s *SQLStore) UseInventoryItem(ctx context.Context, userID int, itemID int) error {
	// Decrement quantity; delete if it reaches 0
	res, err := s.db.ExecContext(ctx,
		"UPDATE UserInventory SET Quantity = Quantity - 1 WHERE UserId = ? AND ItemId = ? AND Quantity > 0", userID, itemID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not in inventory or already used")
	}
	// Clean up zero-quantity rows
	_, _ = s.db.ExecContext(ctx, "DELETE FROM UserInventory WHERE UserId = ? AND ItemId = ? AND Quantity <= 0", userID, itemID)
	return nil
}

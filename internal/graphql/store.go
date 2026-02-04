package graphql

import (
	"context"

	"github.com/jmoiron/sqlx"

	"tamagoam/internal/models"
)

type Store interface {
	ListUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error)
	ListRaces(ctx context.Context) ([]models.Race, error)
}

type SQLStore struct {
	db *sqlx.DB
}

func NewSQLStore(db *sqlx.DB) *SQLStore {
	return &SQLStore{db: db}
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

func (s *SQLStore) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	res, err := s.db.ExecContext(ctx, `INSERT INTO Users (Name, LastName, UserName, Email, ProfilPicture, GamingTime, LastConnectionDate) VALUES (?, ?, ?, ?, ?, ?, ?)`, input.Name, input.LastName, input.UserName, input.Email, input.ProfilPicture, input.GamingTime, input.LastConnectionDate)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetUser(ctx, int(id))
}

func (s *SQLStore) ListRaces(ctx context.Context) ([]models.Race, error) {
	var races []models.Race
	err := s.db.SelectContext(ctx, &races, "SELECT * FROM Race ORDER BY Name")
	return races, err
}

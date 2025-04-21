package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/models"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrPVZCityForbidden = errors.New("forbidden city")
	ErrSessionOpen      = errors.New("session already open")
	ErrNoOpenSession    = errors.New("no open session")
)

type Service struct {
	db  *sql.DB
	sql sq.StatementBuilderType
}

func New(db *sql.DB) *Service {
	return &Service{db: db, sql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (s *Service) CreateUser(ctx context.Context, email, pass, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.sql.Insert("users").Columns("email", "password", "role").Values(email, string(hash), role).RunWith(s.db).Exec()
	return err
}

func (s *Service) Authenticate(ctx context.Context, email, pass string) (*models.User, error) {
	var u models.User
	err := s.sql.Select("id", "email", "password", "role").From("users").Where(sq.Eq{"email": email}).RunWith(s.db).QueryRow().Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)) != nil {
		return nil, ErrUnauthorized
	}
	return &u, nil
}

func (s *Service) CreatePVZ(ctx context.Context, city string) (*models.PVZ, error) {
	allowed := map[string]bool{"Москва": true, "Санкт-Петербург": true, "Казань": true}
	var u models.PVZ
	if !allowed[city] {
		return nil, ErrPVZCityForbidden
	}
	err := s.sql.Insert("pvz").Columns("city").Values(city).Suffix("RETURNING id, city, registered_at").RunWith(s.db).QueryRow().Scan(&u.ID, &u.City, &u.RegisteredAt)
	if err != nil {
		return nil, err
	}
	return &models.PVZ{ID: u.ID, City: u.City, RegisteredAt: u.RegisteredAt}, nil
}

func (s *Service) ListPVZ(ctx context.Context, from, to time.Time, limit, offset uint64) ([]models.PVZ, error) {
	rows, err := s.sql.Select("id", "city", "registered_at").From("pvz").RunWith(s.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.PVZ
	for rows.Next() {
		var p models.PVZ
		if err := rows.Scan(&p.ID, &p.City, &p.RegisteredAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (s *Service) StartSession(ctx context.Context, pvzID int64) (*models.Session, error) {
	var count int
	err := s.sql.Select("count(*)").From("sessions").Where(sq.Eq{"pvz_id": pvzID, "status": "in_progress"}).RunWith(s.db).QueryRow().Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrSessionOpen
	}
	var sess models.Session
	err = s.sql.Insert("sessions").Columns("pvz_id", "status").Values(pvzID, "in_progress").Suffix("RETURNING id, pvz_id, started_at, status").RunWith(s.db).QueryRow().Scan(&sess.ID, &sess.PVZID, &sess.StartedAt, &sess.Status)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *Service) AddItem(ctx context.Context, pvzID, sessionID int64, itemType string) (*models.Item, error) {
	var status string
	var sid int64
	err := s.sql.Select("id", "status").From("sessions").Where(sq.Eq{"pvz_id": pvzID, "id": sessionID}).RunWith(s.db).QueryRow().Scan(&sid, &status)
	if err != nil {
		return nil, ErrNoOpenSession
	}
	if status != "in_progress" {
		return nil, ErrSessionOpen
	}
	var item models.Item
	err = s.sql.Insert("items").Columns("session_id", "type").Values(sessionID, itemType).Suffix("RETURNING id, session_id, type, added_at").RunWith(s.db).QueryRow().Scan(&item.ID, &item.SessionID, &item.Type, &item.AddedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Service) DeleteLastItem(ctx context.Context, pvzID, sessionID int64) error {
	var status string
	err := s.sql.Select("status").From("sessions").Where(sq.Eq{"pvz_id": pvzID, "id": sessionID}).RunWith(s.db).QueryRow().Scan(&status)
	if err != nil || status != "in_progress" {
		return ErrNoOpenSession
	}
	var id int64
	err = s.sql.Select("id").From("items").Where(sq.Eq{"session_id": sessionID}).OrderBy("id DESC").Limit(1).RunWith(s.db).QueryRow().Scan(&id)
	if err != nil {
		return err
	}
	_, err = s.sql.Delete("items").Where(sq.Eq{"id": id}).RunWith(s.db).Exec()
	return err
}

func (s *Service) CloseSession(ctx context.Context, pvzID, sessionID int64) error {
	res, err := s.sql.Update("sessions").Set("status", "close").Where(sq.Eq{"pvz_id": pvzID, "id": sessionID, "status": "in_progress"}).RunWith(s.db).Exec()
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNoOpenSession
	}
	return nil
}

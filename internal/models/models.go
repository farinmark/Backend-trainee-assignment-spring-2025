package models

type User struct {
	ID       int64
	Email    string
	Password string
	Role     string
}

type PVZ struct {
	ID           int64  `db:"id" json:"id"`
	City         string `db:"city" json:"city"`
	RegisteredAt string `db:"registered_at" json:"registered_at"`
}

type Session struct {
	ID        int64  `db:"id" json:"id"`
	PVZID     int64  `db:"pvz_id" json:"pvz_id"`
	StartedAt string `db:"started_at" json:"started_at"`
	Status    string `db:"status" json:"status"`
}

type Item struct {
	ID        int64  `db:"id" json:"id"`
	SessionID int64  `db:"session_id" json:"session_id"`
	Type      string `db:"type" json:"type"`
	AddedAt   string `db:"added_at" json:"added_at"`
}

package adapter

import (
	"context"
	"database/sql"
	"driver-service/internal/domain"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Import driver
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Simple migration on startup
	query := `
	CREATE TABLE IF NOT EXISTS drivers (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		phone TEXT,
		license_number TEXT,
		status TEXT,
		suspend_reason TEXT,
		created_at DATETIME
	);
	`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

func (r *SQLiteDB) Create(ctx context.Context, d *domain.Driver) error {
	query := `
		INSERT INTO drivers (id, name, phone, license_number, status, suspend_reason, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		d.ID, d.Name, d.Phone, d.LicenseNumber, d.Status, d.SuspendReason, d.CreatedAt,
	)
	return err
}

func (r *SQLiteDB) Get(ctx context.Context, id string) (*domain.Driver, error) {
	query := `SELECT id, name, phone, license_number, status, suspend_reason, created_at FROM drivers WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var d domain.Driver
	var reason sql.NullString // Handle nullable field

	err := row.Scan(&d.ID, &d.Name, &d.Phone, &d.LicenseNumber, &d.Status, &reason, &d.CreatedAt)
	if err != nil {
		return nil, err
	}

	if reason.Valid {
		d.SuspendReason = &reason.String
	}
	return &d, nil
}

func (r *SQLiteDB) List(ctx context.Context, offset, limit int) ([]domain.Driver, int, error) {
	// 1. Get Data
	query := `
		SELECT id, name, phone, license_number, status, suspend_reason, created_at 
		FROM drivers 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var drivers []domain.Driver
	for rows.Next() {
		var d domain.Driver
		var reason sql.NullString
		if err := rows.Scan(&d.ID, &d.Name, &d.Phone, &d.LicenseNumber, &d.Status, &reason, &d.CreatedAt); err != nil {
			return nil, 0, err
		}
		if reason.Valid {
			d.SuspendReason = &reason.String
		}
		drivers = append(drivers, d)
	}

	// 2. Get Total Count (for pagination meta)
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM drivers").Scan(&total); err != nil {
		return nil, 0, err
	}

	return drivers, total, nil
}

func (r *SQLiteDB) Suspend(ctx context.Context, id string, reason string) error {
	query := `UPDATE drivers SET status = ?, suspend_reason = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, domain.StatusSuspended, reason, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("driver not found")
	}
	return nil
}

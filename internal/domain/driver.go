package domain

import (
	"context"
	"time"
)

type DriverStatus string

const (
	StatusActive    DriverStatus = "active"
	StatusSuspended DriverStatus = "suspended"
)

type Driver struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Phone         string       `json:"phone"`
	LicenseNumber string       `json:"license_number"`
	Status        DriverStatus `json:"status"`
	SuspendReason *string      `json:"suspend_reason,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

type ListResponse struct {
	Data []Driver   `json:"data"`
	Meta Pagination `json:"meta"`
}

// Repository defines database operations
type DriverRepository interface {
	Create(ctx context.Context, d *Driver) error
	List(ctx context.Context, offset, limit int) ([]Driver, int, error)
	Suspend(ctx context.Context, id string, reason string) error
	Get(ctx context.Context, id string) (*Driver, error)
}

// Cacher defines caching strategy
type DriverCache interface {
	GetList(ctx context.Context, page, limit int) (*ListResponse, bool)
	SetList(ctx context.Context, page, limit int, data *ListResponse)
	InvalidateList(ctx context.Context) // Clear list cache on updates
}

// EventQueue defines async messaging
type EventQueue interface {
	Publish(ctx context.Context, topic string, payload interface{}) error
}

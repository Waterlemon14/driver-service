package service

import (
	"context"
	"driver-service/internal/domain"
	"fmt"
	"github.com/google/uuid" // You might need to add this or use simple math/rand for ID
	"math"
	"time"
)

type DriverService struct {
	repo  domain.DriverRepository
	cache domain.DriverCache
	queue domain.EventQueue
}

func NewDriverService(r domain.DriverRepository, c domain.DriverCache, q domain.EventQueue) *DriverService {
	return &DriverService{repo: r, cache: c, queue: q}
}

func (s *DriverService) CreateDriver(ctx context.Context, req domain.Driver) (*domain.Driver, error) {
	// 1. Validation
	if req.Name == "" || req.LicenseNumber == "" {
		return nil, fmt.Errorf("name and license_number are required")
	}

	// 2. New Driver
	driver := &domain.Driver{
		ID:            uuid.New().String(),
		Name:          req.Name,
		Phone:         req.Phone,
		LicenseNumber: req.LicenseNumber,
		Status:        domain.StatusActive,
		CreatedAt:     time.Now(),
	}

	// 3. Save
	if err := s.repo.Create(ctx, driver); err != nil {
		return nil, err
	}

	// 4. Invalidate Cache (New data available)
	s.cache.InvalidateList(ctx)

	// 5. Publish Event (Async)
	// Ignore error for now to not block response,
	// Prefer to log error to trace
	_ = s.queue.Publish(ctx, "driver.created", map[string]string{"id": driver.ID})

	return driver, nil
}

func (s *DriverService) ListDrivers(ctx context.Context, page, limit int) (*domain.ListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// 1. Check Cache
	if cached, hit := s.cache.GetList(ctx, page, limit); hit {
		return cached, nil
	}

	// 2. Query DB
	offset := (page - 1) * limit
	drivers, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	// 3. Response
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	resp := &domain.ListResponse{
		Data: drivers,
		Meta: domain.Pagination{
			Page:       page,
			Limit:      limit,
			TotalCount: total,
			TotalPages: totalPages,
		},
	}

	// 4. Set Cache
	s.cache.SetList(ctx, page, limit, resp)

	return resp, nil
}

func (s *DriverService) SuspendDriver(ctx context.Context, id string, reason string) error {
	if reason == "" {
		return fmt.Errorf("reason is required")
	}

	// 1. Update DB
	if err := s.repo.Suspend(ctx, id, reason); err != nil {
		return err
	}

	// 2. Invalidate Cache (Data changed)
	s.cache.InvalidateList(ctx)

	// 3. Publish Event
	_ = s.queue.Publish(ctx, "driver.suspended", map[string]string{"id": id, "reason": reason})

	return nil
}

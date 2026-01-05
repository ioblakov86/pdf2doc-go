package jobs

import (
	"context"
	"sync"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusDone      Status = "done"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

type Job struct {
	ID        string
	Status    Status
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Result    []byte
	cancel    context.CancelFunc
}

type Manager struct {
	mu    sync.RWMutex
	jobs  map[string]*Job
	queue chan *Job
	sem   chan struct{}
}

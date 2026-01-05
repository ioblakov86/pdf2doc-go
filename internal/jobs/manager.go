package jobs

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func NewManager(maxConcurrent int) *Manager {
	m := &Manager{
		jobs:  make(map[string]*Job),
		queue: make(chan *Job, 100),
		sem:   make(chan struct{}, maxConcurrent),
	}

	for i := 0; i < maxConcurrent; i++ {
		go m.worker()
	}

	return m
}

func (m *Manager) worker() {
	for job := range m.queue {
		m.sem <- struct{}{} // ждём свободный слот

		go func(j *Job) {
			defer func() { <-m.sem }()

			m.run(job)

		}(job)
	}
}

func (m *Manager) run(job *Job) {
	ctx, cancel := context.WithCancel(context.Background())
	job.cancel = cancel

	m.setStatus(job.ID, StatusRunning)

	result, err := Convert(ctx, job.ID) // наша бизнес логика

	if err != nil {
		if err == context.Canceled {
			m.setStatus(job.ID, StatusCancelled)
		} else {
			m.setError(job.ID, err)
		}
		return
	}

	m.mu.Lock()
	job.Status = StatusDone
	job.Result = result
	job.UpdatedAt = time.Now()
	m.mu.Unlock()
}

func (m *Manager) Submit() string {
	job := &Job{
		ID:        uuid.NewString(),
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}

	m.mu.Lock()
	m.jobs[job.ID] = job
	m.mu.Unlock()

	m.queue <- job
	return job.ID
}

func (m *Manager) Get(id string) (*Job, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	job, ok := m.jobs[id]
	return job, ok
}

func (m *Manager) Cancel(id string) {
	m.mu.RLock()
	job, ok := m.jobs[id]
	m.mu.RUnlock()
	if ok && job.cancel != nil {
		job.cancel()
	}
}

func (m *Manager) setStatus(id string, status Status) {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, ok := m.jobs[id]
	if !ok {
		return
	}
	job.Status = status
}

func (m *Manager) setError(id string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, ok := m.jobs[id]
	if !ok {
		return
	}
	job.Status = StatusFailed
	job.Error = err.Error()
}

func Convert(ctx context.Context, id string) ([]byte, error) {
	select {
	case <-time.After(5 * time.Second):
		return []byte("fake docx content"), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

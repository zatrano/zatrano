package queue

import (
	"fmt"
	"sync"
)

// FakeQueue captures dispatched jobs in memory for testing.
type FakeQueue struct {
	mu   sync.RWMutex
	jobs []*DispatchedJob
}

// DispatchedJob represents a captured job dispatch.
type DispatchedJob struct {
	Name string
	Data interface{}
}

// NewFakeQueue creates a new fake queue.
func NewFakeQueue() *FakeQueue {
	return &FakeQueue{
		jobs: make([]*DispatchedJob, 0),
	}
}

// Dispatch captures the job instead of dispatching it.
func (f *FakeQueue) Dispatch(name string, data interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	job := &DispatchedJob{
		Name: name,
		Data: data,
	}
	f.jobs = append(f.jobs, job)
	return nil
}

// GetDispatchedJobs returns all captured jobs.
func (f *FakeQueue) GetDispatchedJobs() []*DispatchedJob {
	f.mu.RLock()
	defer f.mu.RUnlock()

	jobs := make([]*DispatchedJob, len(f.jobs))
	copy(jobs, f.jobs)
	return jobs
}

// ClearJobs clears all captured jobs.
func (f *FakeQueue) ClearJobs() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.jobs = nil
}

// AssertJobDispatched asserts that a job with the given name was dispatched.
func (f *FakeQueue) AssertJobDispatched(name string) {
	jobs := f.GetDispatchedJobs()
	for _, job := range jobs {
		if job.Name == name {
			return
		}
	}
	panic(fmt.Sprintf("expected job %q not dispatched", name))
}

// AssertJobNotDispatched asserts that a job with the given name was not dispatched.
func (f *FakeQueue) AssertJobNotDispatched(name string) {
	jobs := f.GetDispatchedJobs()
	for _, job := range jobs {
		if job.Name == name {
			panic(fmt.Sprintf("expected job %q not to be dispatched", name))
		}
	}
}

// AssertJobCount asserts the number of jobs dispatched.
func (f *FakeQueue) AssertJobCount(count int) {
	jobs := f.GetDispatchedJobs()
	if len(jobs) != count {
		panic(fmt.Sprintf("expected %d jobs, got %d", count, len(jobs)))
	}
}

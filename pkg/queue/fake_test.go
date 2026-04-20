package queue

import (
	"testing"
)

func TestFakeQueue_Dispatch(t *testing.T) {
	fake := NewFakeQueue()

	err := fake.Dispatch("test_job", map[string]interface{}{"key": "value"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	jobs := fake.GetDispatchedJobs()
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}

	if jobs[0].Name != "test_job" {
		t.Errorf("expected job name test_job, got %s", jobs[0].Name)
	}
}

func TestFakeQueue_AssertJobDispatched(t *testing.T) {
	fake := NewFakeQueue()
	fake.Dispatch("welcome_email", nil)

	// Should not panic
	fake.AssertJobDispatched("welcome_email")
}

func TestFakeQueue_AssertJobNotDispatched(t *testing.T) {
	fake := NewFakeQueue()

	// Should not panic
	fake.AssertJobNotDispatched("non_existent_job")
}

package schedule

import (
	"context"
	"testing"
	"time"
)

func TestBuilderRegisterAndList(t *testing.T) {
	Reset()
	defer Reset()

	job := func(ctx context.Context) error {
		return nil
	}

	task, err := Call(job).
		Name("daily_report").
		Description("Daily summary report").
		Daily().
		At("08:00").
		WithoutOverlap().
		Register()
	if err != nil {
		t.Fatalf("expected register to succeed, got %v", err)
	}

	if task.Spec != "0 8 * * *" {
		t.Fatalf("expected daily at 08:00 spec, got %q", task.Spec)
	}

	tasks := Tasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 registered task, got %d", len(tasks))
	}

	if tasks[0].Name != "daily_report" {
		t.Fatalf("expected task name daily_report, got %q", tasks[0].Name)
	}
}

func TestBuilderInvalidSpec(t *testing.T) {
	Reset()
	defer Reset()

	_, err := Call(func(context.Context) error { return nil }).
		Name("bad_spec").
		WithSpec("invalid spec").
		Register()
	if err == nil {
		t.Fatal("expected invalid spec error")
	}
}

func TestBuilderWithoutOverlapDefaults(t *testing.T) {
	Reset()
	defer Reset()

	task, err := Call(func(context.Context) error { return nil }).
		Name("lock_test").
		EveryMinute().
		WithoutOverlap().
		Register()
	if err != nil {
		t.Fatalf("expected register to succeed, got %v", err)
	}
	if task.LockTTL != 5*time.Minute {
		t.Fatalf("expected default lock ttl 5m, got %v", task.LockTTL)
	}
}

package schedule

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

var (
	registry          = &taskRegistry{tasks: make(map[string]*Task)}
	releaseLockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`)
)

type Task struct {
	Name        string
	Description string
	Spec        string
	Job         func(context.Context) error
	SkipOverlap bool
	LockTTL     time.Duration
	createdAt   time.Time
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

type taskRegistry struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func (r *taskRegistry) add(task *Task) error {
	if task == nil {
		return fmt.Errorf("task is nil")
	}
	if strings.TrimSpace(task.Name) == "" {
		return fmt.Errorf("task name is required")
	}
	if strings.TrimSpace(task.Spec) == "" {
		return fmt.Errorf("cron spec is required")
	}
	if task.Job == nil {
		return fmt.Errorf("task job is required")
	}
	if task.LockTTL <= 0 {
		task.LockTTL = 5 * time.Minute
	}
	if _, err := cron.ParseStandard(task.Spec); err != nil {
		return fmt.Errorf("invalid cron spec %q: %w", task.Spec, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.Name]; exists {
		return fmt.Errorf("task %q already registered", task.Name)
	}
	task.createdAt = time.Now()
	r.tasks[task.Name] = task
	return nil
}

func (r *taskRegistry) list() []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var tasks []*Task
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (r *taskRegistry) clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks = make(map[string]*Task)
}

func Call(job func(context.Context) error) *Builder {
	return &Builder{task: &Task{Job: job, LockTTL: 5 * time.Minute}}
}

type Builder struct {
	task *Task
	err  error
}

func (b *Builder) Name(name string) *Builder {
	b.task.Name = strings.TrimSpace(name)
	return b
}

func (b *Builder) Description(description string) *Builder {
	b.task.Description = description
	return b
}

func (b *Builder) EveryMinute() *Builder {
	b.task.Spec = "* * * * *"
	return b
}

func (b *Builder) Hourly() *Builder {
	b.task.Spec = "0 * * * *"
	return b
}

func (b *Builder) Daily() *Builder {
	b.task.Spec = "0 0 * * *"
	return b
}

func (b *Builder) Weekly() *Builder {
	b.task.Spec = "0 0 * * 0"
	return b
}

func (b *Builder) Monthly() *Builder {
	b.task.Spec = "0 0 1 * *"
	return b
}

func (b *Builder) At(value string) *Builder {
	if b.err != nil {
		return b
	}
	t, err := time.Parse("15:04", value)
	if err != nil {
		b.err = fmt.Errorf("invalid time format %q: %w", value, err)
		return b
	}

	parts := strings.Fields(b.task.Spec)
	if len(parts) != 5 {
		b.err = fmt.Errorf("cannot set time on invalid cron spec %q", b.task.Spec)
		return b
	}

	parts[0] = "0"
	parts[1] = strconv.Itoa(t.Hour())
	b.task.Spec = strings.Join(parts, " ")
	return b
}

func (b *Builder) WithSpec(spec string) *Builder {
	b.task.Spec = strings.TrimSpace(spec)
	return b
}

func (b *Builder) WithoutOverlap() *Builder {
	b.task.SkipOverlap = true
	return b
}

func (b *Builder) WithLockTTL(ttl time.Duration) *Builder {
	b.task.LockTTL = ttl
	return b
}

func (b *Builder) Register() (*Task, error) {
	if b.err != nil {
		return nil, b.err
	}
	if err := registry.add(b.task); err != nil {
		return nil, err
	}
	return b.task, nil
}

func Register(task *Task) error {
	return registry.add(task)
}

func Tasks() []*Task {
	return registry.list()
}

func Reset() {
	registry.clear()
}

type Scheduler struct {
	cron       *cron.Cron
	redis      *redis.Client
	registered map[string]cron.EntryID
}

func New(redisClient *redis.Client) *Scheduler {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	return &Scheduler{
		cron:       cron.New(cron.WithParser(parser), cron.WithChain()),
		redis:      redisClient,
		registered: make(map[string]cron.EntryID),
	}
}

func (s *Scheduler) Register(task *Task) error {
	if task == nil {
		return fmt.Errorf("task is nil")
	}
	if _, err := cron.ParseStandard(task.Spec); err != nil {
		return fmt.Errorf("invalid cron spec %q: %w", task.Spec, err)
	}
	job := func() {
		ctx := context.Background()
		if task.SkipOverlap {
			if s.redis == nil {
				log.Printf("warning: skipping overlap prevention for task %q because redis is not configured", task.Name)
			} else {
				release, ok, err := s.acquireLock(ctx, task)
				if err != nil {
					log.Printf("failed to acquire lock for task %q: %v", task.Name, err)
					return
				}
				if !ok {
					log.Printf("task %q skipped because another instance is running", task.Name)
					return
				}
				defer func() {
					if err := release(ctx); err != nil {
						log.Printf("failed to release lock for task %q: %v", task.Name, err)
					}
				}()
			}
		}

		if err := task.Job(ctx); err != nil {
			log.Printf("task %q failed: %v", task.Name, err)
		}
	}

	id, err := s.cron.AddFunc(task.Spec, job)
	if err != nil {
		return fmt.Errorf("register task %q: %w", task.Name, err)
	}
	s.registered[task.Name] = id
	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() context.Context {
	return s.cron.Stop()
}

func (s *Scheduler) Run(ctx context.Context) error {
	s.Start()
	<-ctx.Done()
	s.Stop()
	return nil
}

func (s *Scheduler) acquireLock(ctx context.Context, task *Task) (func(context.Context) error, bool, error) {
	key := fmt.Sprintf("zatrano:schedule:lock:%s", task.Name)
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	ok, err := s.redis.SetNX(ctx, key, value, task.LockTTL).Result()
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	return func(ctx context.Context) error {
		return releaseLockScript.Run(ctx, s.redis, []string{key}, value).Err()
	}, true, nil
}

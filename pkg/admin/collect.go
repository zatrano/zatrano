package admin

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/zatrano/zatrano/pkg/core"
	"github.com/zatrano/zatrano/pkg/queue"
)

// QueueDepth holds ready and delayed counts for one named queue (Redis layout).
type QueueDepth struct {
	Name    string
	Ready   int64
	Delayed int64
}

// MetricsSnapshot is passed to admin/metrics templates.
type MetricsSnapshot struct {
	Queues       []QueueDepth
	FailedJobs   int64
	FailedJobsOK bool
	RedisHits    int64
	RedisMisses  int64
	RedisInfoOK  bool
	CacheHitRate string
}

func collectMetrics(ctx context.Context, a *core.App) MetricsSnapshot {
	out := MetricsSnapshot{CacheHitRate: "—"}
	if a == nil {
		return out
	}
	names := a.Config.Admin.QueueNames
	if len(names) == 0 {
		names = []string{"default"}
	}
	if a.Redis != nil {
		for _, q := range names {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}
			rd, err := a.Redis.LLen(ctx, queue.RedisReadyListKey(q)).Result()
			if err != nil {
				rd = -1
			}
			dd, err := a.Redis.ZCard(ctx, queue.RedisDelayedZSetKey(q)).Result()
			if err != nil {
				dd = -1
			}
			out.Queues = append(out.Queues, QueueDepth{Name: q, Ready: rd, Delayed: dd})
		}
		info, err := a.Redis.Info(ctx, "stats").Result()
		if err == nil {
			out.RedisInfoOK = true
			for _, line := range strings.Split(info, "\r\n") {
				if strings.HasPrefix(line, "keyspace_hits:") {
					out.RedisHits, _ = strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, "keyspace_hits:")), 10, 64)
				}
				if strings.HasPrefix(line, "keyspace_misses:") {
					out.RedisMisses, _ = strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, "keyspace_misses:")), 10, 64)
				}
			}
			h, m := out.RedisHits, out.RedisMisses
			if h+m > 0 {
				out.CacheHitRate = strconv.FormatFloat(float64(h)/float64(h+m)*100, 'f', 1, 64) + "% (Redis keyspace)"
			}
		}
	}
	if a.DB != nil {
		var n int64
		tx := a.DB.WithContext(ctx).Raw("SELECT COUNT(*) FROM zatrano_failed_jobs").Scan(&n)
		if tx.Error == nil {
			out.FailedJobs = n
			out.FailedJobsOK = true
		}
	}
	return out
}

// tailLogFile returns up to maxLines from the end of a text file.
func tailLogFile(path string, maxLines int) ([]string, error) {
	if path == "" || maxLines <= 0 {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
		if len(lines) > maxLines*4 {
			lines = lines[len(lines)-maxLines*2:]
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}
	return lines, nil
}

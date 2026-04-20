package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/zatrano/zatrano/pkg/schedule"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Scheduled task commands",
}

var scheduleRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run registered scheduled tasks",
	Long: `Starts the scheduler and runs all registered tasks based on their cron expressions.

Tasks are registered via pkg/schedule. Use Redis when overlap prevention is enabled.
Examples:
  zatrano schedule run
  zatrano schedule run --env dev --config-dir config`,
	RunE: runScheduleRun,
}

var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered scheduled tasks",
	RunE:  runScheduleList,
}

func init() {
	for _, cmd := range []*cobra.Command{scheduleRunCmd, scheduleListCmd} {
		cmd.Flags().String("env", "", "environment profile")
		cmd.Flags().String("config-dir", "config", "config directory")
		cmd.Flags().Bool("no-dotenv", false, "skip .env loading")
	}

	scheduleCmd.AddCommand(scheduleRunCmd, scheduleListCmd)
	rootCmd.AddCommand(scheduleCmd)
}

func runScheduleRun(cmd *cobra.Command, _ []string) error {
	cfg, err := loadQueueConfig(cmd)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	tasks := schedule.Tasks()
	if len(tasks) == 0 {
		fmt.Println("No scheduled tasks registered.")
		return nil
	}

	var client *redis.Client
	for _, task := range tasks {
		if task.SkipOverlap {
			client, err = openRedisForQueue(cfg)
			if err != nil {
				return err
			}
			break
		}
	}

	if client != nil {
		defer func() { _ = client.Close() }()
	}

	scheduler := schedule.New(client)
	for _, task := range tasks {
		if err := scheduler.Register(task); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\nShutting down scheduler...")
		cancel()
	}()

	scheduler.Start()
	fmt.Printf("Scheduler started with %d task(s). Press Ctrl+C to stop.\n", len(tasks))
	<-ctx.Done()
	return nil
}

func runScheduleList(cmd *cobra.Command, _ []string) error {
	tasks := schedule.Tasks()
	if len(tasks) == 0 {
		fmt.Println("No scheduled tasks registered.")
		return nil
	}

	fmt.Printf("%-24s %-18s %-10s %s\n", "NAME", "SCHEDULE", "OVERLAP", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 80))
	for _, task := range tasks {
		overlap := "no"
		if task.SkipOverlap {
			overlap = "yes"
		}
		fmt.Printf("%-24s %-18s %-10s %s\n", task.Name, task.Spec, overlap, task.Description)
	}
	return nil
}

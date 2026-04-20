package zatrano

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/zatrano/zatrano/pkg/config"
	"github.com/zatrano/zatrano/pkg/core"
	"github.com/zatrano/zatrano/pkg/meta"
	"github.com/zatrano/zatrano/pkg/server"
)

// StartOptions configures the embedded HTTP server (same behavior as `zatrano serve`).
type StartOptions struct {
	Env       string
	ConfigDir string
	Addr      string
	NoDotenv  bool
	// RegisterRoutes mounts app-specific modules (e.g. internal/routes.Register). Optional.
	RegisterRoutes func(a *core.App, app *fiber.App)
	// ShutdownHooks run after SIGINT/SIGTERM or a graceful upgrade handoff, before Fiber
	// begins draining. They share the same deadline as http.shutdown_timeout.
	ShutdownHooks []func(context.Context) error
}

// Start boots the ZATRANO HTTP server. Intended for generated apps: `zatrano.Run()`.
func Start(opts StartOptions) error {
	cfg, err := config.Load(config.LoadOptions{
		Env:       opts.Env,
		ConfigDir: opts.ConfigDir,
		DotEnv:    !opts.NoDotenv,
	})
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	if opts.Addr != "" {
		cfg.HTTPAddr = opts.Addr
	}

	app, err := core.Bootstrap(cfg)
	if err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}
	defer func() {
		if cerr := app.Close(); cerr != nil {
			app.Log.Warn("shutdown resources", zap.Error(cerr))
		}
	}()

	fiberApp := core.NewFiber(app)
	server.Mount(app, fiberApp, server.MountOptions{RegisterRoutes: opts.RegisterRoutes})
	app.Fiber = fiberApp

	shutdownTimeout := cfg.HTTP.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = 15 * time.Second
	}

	app.Log.Info("zatrano starting",
		zap.String("version", meta.Version),
		zap.String("env", cfg.Env),
		zap.String("addr", cfg.HTTPAddr),
		zap.Bool("graceful_restart", cfg.HTTP.GracefulRestart),
	)

	errCh := make(chan error, 1)
	ready := make(chan struct{})
	var readyOnce sync.Once
	markReady := func() {
		readyOnce.Do(func() { close(ready) })
	}

	listenCfg := fiber.ListenConfig{
		DisableStartupMessage: true,
		BeforeServeFunc: func(*fiber.App) error {
			markReady()
			return nil
		},
	}

	var upg *tableflip.Upgrader
	if cfg.HTTP.GracefulRestart {
		upg, err = tableflip.New(tableflip.Options{
			PIDFile: strings.TrimSpace(cfg.HTTP.GracefulRestartPIDFile),
		})
		if err != nil {
			return fmt.Errorf("graceful_restart (tableflip): %w", err)
		}
		defer upg.Stop()

		goGracefulUSR2Upgrade(app.Log, upg)

		ln, lerr := upg.Listen("tcp", strings.TrimSpace(cfg.HTTPAddr))
		if lerr != nil {
			return fmt.Errorf("listen: %w", lerr)
		}

		listenCfg.BeforeServeFunc = func(*fiber.App) error {
			markReady()
			return upg.Ready()
		}

		go func() {
			errCh <- fiberApp.Listener(ln, listenCfg)
		}()
	} else {
		go func() {
			errCh <- fiberApp.Listen(strings.TrimSpace(cfg.HTTPAddr), listenCfg)
		}()
	}

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
		return fmt.Errorf("server exited before becoming ready")
	case <-ready:
		app.Log.Info("listening", zap.String("url", localBaseURL(cfg.HTTPAddr)))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if upg != nil {
		select {
		case <-ctx.Done():
		case <-upg.Exit():
			app.Log.Info("graceful restart handoff: draining this process")
		}
	} else {
		<-ctx.Done()
	}
	stop()

	app.Log.Info("shutdown signal received, draining connections...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	for i, h := range opts.ShutdownHooks {
		if h == nil {
			continue
		}
		if herr := h(shutdownCtx); herr != nil {
			app.Log.Warn("shutdown hook failed", zap.Int("index", i), zap.Error(herr))
		}
	}

	if err := fiberApp.ShutdownWithContext(shutdownCtx); err != nil {
		return fmt.Errorf("fiber shutdown: %w", err)
	}

	if err := <-errCh; err != nil {
		return fmt.Errorf("server: %w", err)
	}
	app.Log.Info("server stopped cleanly")
	return nil
}

// Run reads environment variables and config from the working directory (with .env) and starts the server.
func Run() error {
	return Start(StartOptions{})
}

func localBaseURL(addr string) string {
	if len(addr) > 0 && addr[0] == ':' {
		return "http://127.0.0.1" + addr
	}
	return "http://" + addr
}

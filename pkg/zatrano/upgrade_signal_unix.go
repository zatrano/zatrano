//go:build !windows

package zatrano

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudflare/tableflip"
	"go.uber.org/zap"
)

func goGracefulUSR2Upgrade(log *zap.Logger, upg *tableflip.Upgrader) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR2)
	go func() {
		for range ch {
			if err := upg.Upgrade(); err != nil {
				log.Warn("graceful restart (SIGUSR2) failed", zap.Error(err))
			} else {
				log.Info("graceful restart: new process spawned (SIGUSR2); old process will drain")
			}
		}
	}()
}

//go:build windows

package zatrano

import (
	"github.com/cloudflare/tableflip"
	"go.uber.org/zap"
)

func goGracefulUSR2Upgrade(_ *zap.Logger, _ *tableflip.Upgrader) {}

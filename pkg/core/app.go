package core

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/zatrano/zatrano/pkg/audit"
	"github.com/zatrano/zatrano/pkg/auth"
	"github.com/zatrano/zatrano/pkg/broadcast"
	"github.com/zatrano/zatrano/pkg/cache"
	"github.com/zatrano/zatrano/pkg/config"
	"github.com/zatrano/zatrano/pkg/events"
	"github.com/zatrano/zatrano/pkg/features"
	"github.com/zatrano/zatrano/pkg/i18n"
	"github.com/zatrano/zatrano/pkg/mail"
	"github.com/zatrano/zatrano/pkg/notifications"
	"github.com/zatrano/zatrano/pkg/queue"
	"github.com/zatrano/zatrano/pkg/search"
	"github.com/zatrano/zatrano/pkg/view"
)

// App is the root application container.
type App struct {
	Config *config.Config
	Log    *zap.Logger
	Fiber  *fiber.App

	DB    *gorm.DB
	Redis *redis.Client

	// SessionStore is set when Redis-backed sessions are enabled (use session.FromContext in handlers).
	SessionStore *session.Store

	// I18n is loaded when config i18n.enabled (JSON catalogs under locales_dir).
	I18n *i18n.Bundle

	// Gate is the resource-based authorization registry (define/check abilities).
	Gate *auth.Gate

	// RBAC is the role-based access control manager (role → permission mapping, DB-backed).
	RBAC *auth.RBACManager

	// Cache is the application cache manager (memory or Redis driver).
	Cache *cache.Manager

	// Queue is the background job queue manager (Redis-backed).
	Queue *queue.Manager

	// Mail is the email sending manager (SMTP, log drivers).
	Mail *mail.Manager

	// Notifications routes outbound notifications (mail channel uses Mail).
	Notifications *notifications.Manager

	// Events is the event dispatcher (pub/sub bus).
	Events *events.Dispatcher

	// Broadcast is the optional WebSocket / SSE hub (nil when disabled in config).
	Broadcast *broadcast.Hub

	// Audit is the audit writer (always non-nil; NopWriter when audit.enabled is false).
	Audit audit.Writer

	// Search is optional Meilisearch / Typesense client (nil when search.enabled is false).
	Search *search.Client

	// Features is the feature-flag registry (always non-nil; no-op when features.enabled is false).
	Features *features.Registry

	// View is the server-rendered template engine with layout inheritance,
	// component partials, flash messages, old-input, and asset versioning.
	View *view.Renderer
}

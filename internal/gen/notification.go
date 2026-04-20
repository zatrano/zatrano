package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Notification generates a notification stub.
func Notification(moduleRoot, baseDir, rawName string, dryRun bool) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid notification name %q (use letters, digits, _ or -)", rawName)
	}
	pascal := snakeToPascal(name)
	outDir := filepath.Join(moduleRoot, baseDir)
	path := filepath.Join(outDir, name+".go")
	if dryRun {
		return []string{path}, nil
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, []byte(tmplNotification(pascal)), 0o644); err != nil {
		return nil, err
	}
	return []string{path}, nil
}

func tmplNotification(pascal string) string {
	return fmt.Sprintf(`package notifications

import (
	"github.com/zatrano/zatrano/pkg/notifications"
)

// %[1]sNotification is a custom notification.
type %[1]sNotification struct {
	*notifications.BaseNotification
}

// New%[1]sNotification creates a new %[1]s notification.
func New%[1]sNotification(recipient string) *%[1]sNotification {
	return &%[1]sNotification{
		BaseNotification: notifications.NewNotification(
			"TODO: Add subject",
			"TODO: Add body",
			recipient,
		),
	}
}

// TODO: Add custom methods for this notification type.
`, pascal)
}

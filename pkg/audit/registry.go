package audit

import (
	"reflect"
	"sync"
)

var (
	subjectMu sync.RWMutex
	subjects  = map[reflect.Type]string{}
)

// RegisterSubject registers a model type for automatic activity logging.
// Call once per type at application startup, e.g. audit.RegisterSubject[Product]("products").
func RegisterSubject[T any](subjectType string) {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	subjectMu.Lock()
	defer subjectMu.Unlock()
	subjects[t] = subjectType
}

func subjectForModelType(rt reflect.Type) (string, bool) {
	if rt == nil {
		return "", false
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	subjectMu.RLock()
	defer subjectMu.RUnlock()
	s, ok := subjects[rt]
	return s, ok
}

package form

import (
	"html/template"
	"net/url"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Builder, bir HTML formu oluşturmak için gereken tüm durumu ve metodları içerir.
type Builder struct {
	ctx         *fiber.Ctx
	model       interface{}
	oldInput    url.Values
	errors      map[string]string
	csrfToken   string
	csrfField   string
	action      string
	method      string
	isMultipart bool
}

// Config, bir Builder'ı yapılandırmak için opsiyonel ayarları içerir.
type Config struct {
	Action    string
	Method    string
	Model     interface{}
	Multipart bool
}

// New, bir Fiber context'i ve opsiyonel konfigürasyon alarak yeni bir form builder oluşturur.
// CSRF, old input ve hataları otomatik olarak context ve session'dan yönetir.
func New(c *fiber.Ctx, config Config) *Builder {
	// Session'dan flash edilmiş verileri al
	oldInputData, errorsData := getFlashData(c)
	
	return &Builder{
		ctx:         c,
		action:      config.Action,
		method:      config.Method,
		model:       config.Model,
		isMultipart: config.Multipart,
		errors:      errorsData,
		oldInput:    oldInputData,
		csrfToken:   c.Locals("csrf", "").(string), // CSRF middleware'inden al
		csrfField:   "_csrf",
	}
}

// FlashAndRedirect, validasyon hatalarını ve eski girdileri session'a kaydeder
// ve kullanıcıyı belirtilen yola yönlendirir.
func FlashAndRedirect(c *fiber.Ctx, errors map[string]string, path string) error {
	store, ok := c.Locals("session_store").(*session.Store)
	if !ok {
		// Session yoksa, hata ver ama panikleme
		return c.Redirect(path)
	}
	sess, _ := store.Get(c)
	formValues, _ := c.FormValues()

	sess.Set(flashKeyErrors, errors)
	sess.Set(flashKeyOldInput, formValues)
	sess.Save()
	
	flash.Error(c, "Please check the form for errors.") // flash paketi varsayımı
	return c.Redirect(path)
}

// getFlashData, session'dan hataları ve eski girdileri okur ve siler.
func getFlashData(c *fiber.Ctx) (url.Values, map[string]string) {
	store, ok := c.Locals("session_store").(*session.Store)
	if !ok {
		return make(url.Values), make(map[string]string)
	}
	
	sess, err := store.Get(c)
	if err != nil {
		return make(url.Values), make(map[string]string)
	}
	
	var oldInput url.Values
	var errors map[string]string
	
	if old, ok := sess.Get(flashKeyOldInput).(url.Values); ok {
		oldInput = old
		sess.Delete(flashKeyOldInput)
	}
	
	if errs, ok := sess.Get(flashKeyErrors).(map[string]string); ok {
		errors = errs
		sess.Delete(flashKeyErrors)
	}
	
	sess.Save()
	
	return oldInput, errors
}

const (
	flashKeyErrors   = "zatrano_flash_errors"
	flashKeyOldInput = "zatrano_flash_old_input"
)
package view

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

const (
	// Session içinde flash verileri saklamak için kullanılacak anahtarlar
	flashKeyErrors    = "zatrano_flash_errors"
	flashKeyOldInput  = "zatrano_flash_old_input"
	flashKeyMessages  = "zatrano_flash_messages"
)

// Message, tek bir flash mesajının yapısını tutar.
type Message struct {
	Type    string `json:"type"` // "success", "error", "warning", "info"
	Title   string `json:"title"`
	Message string `json:"message"`
}

// NewEngine, ZATRANO'nun v2 uyumlu view motorunu oluşturur.
func NewEngine() *html.Engine {
	engine := html.NewFileSystem(http.Dir("./views"), ".html")
	
	// Geliştirme ortamında şablonları her istekte yeniden yükle.
	// Bu, sunucuyu yeniden başlatmadan .html dosyalarındaki değişiklikleri görmenizi sağlar.
	// Production'da false yapılabilir.
	engine.Reload(true) 
	
	return engine
}

// Render, view render işlemini kolaylaştıran merkezi bir cephedir (facade).
// Handler'dan gelen verileri, context'ten aldığı global verilerle birleştirir
// ve son halini şablona gönderir.
func Render(c *fiber.Ctx, templateName string, data fiber.Map, layout ...string) error {
	// 1. Her istek için global verileri hazırla.
	finalData := collectGlobalData(c)

	// 2. Handler'dan gelen özel verileri global verilere ekle.
	// Handler'dan gelen bir anahtar, global'deki bir anahtarı ezer.
	for key, value := range data {
		finalData[key] = value
	}

	// 3. Layout'u ayarla.
	if len(layout) > 0 {
		finalData["layout"] = layout[0]
	} else if _, ok := finalData["layout"]; !ok {
		finalData["layout"] = "layouts/main" // Varsayılan layout
	}
	
	// 4. Şablonu render et.
	return c.Render(templateName, finalData)
}

// collectGlobalData, her view'a gönderilmesi gereken verileri merkezi olarak toplar.
// CSRF token'ı, flash mesajlar, validasyon hataları, eski form girdileri ve
// giriş yapmış kullanıcı gibi verileri yönetir.
func collectGlobalData(c *fiber.Ctx) fiber.Map {
	// CSRF token'ını al
	csrfToken, _ := c.Locals("csrf").(string)

	// Session'ı al
	store, _ := c.Locals("session_store").(*session.Store)
	sess, err := store.Get(c)
	if err != nil {
		// Session alınamazsa, boş verilerle devam et.
		log.Printf("View Engine: Could not get session - %v", err)
		return fiber.Map{
			"CSRFToken": csrfToken,
		}
	}
	
	// Flash verilerini session'dan oku
	errors := getFlashData(sess, flashKeyErrors)
	oldInput := getFlashData(sess, flashKeyOldInput)
	messages := getFlashData(sess, flashKeyMessages)
	
	// Okunan flash verileri için session'ı kaydet (çünkü getFlashData onları sildi)
	if err := sess.Save(); err != nil {
		log.Printf("View Engine: Could not save session after clearing flash data - %v", err)
	}
	
	// Mesajları JSON'a çevirerek SweetAlert için hazırla
	flashMessagesJson, _ := json.Marshal(messages)

	// TODO: Auth (Giriş yapmış kullanıcı) verisini ekle
	// authUser := c.Locals("user")

	return fiber.Map{
		"CSRFToken":         csrfToken,
		"Errors":            errors,

		// OldInput'u `url.Values` yerine `map[string]string` olarak saklamak daha kolay olabilir,
		// bu, Flash paketinin sorumluluğundadır. Şimdilik bu şekilde bırakıyoruz.
		"OldInput":          oldInput, 

		// SweetAlert için JSON formatında, güvenli bir şekilde
		"FlashMessagesJson": template.JS(flashMessagesJson),

		// "AuthUser": authUser,
	}
}

// getFlashData, session'dan belirli bir anahtardaki veriyi okur ve
// okuduktan hemen sonra o anahtarı session'dan SİLER.
func getFlashData(sess *session.Session, key string) interface{} {
	data := sess.Get(key)
	if data != nil {
		sess.Delete(key)
	}
	return data
}
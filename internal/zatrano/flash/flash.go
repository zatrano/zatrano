package flash

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const sessionKey = "zatrano_flash"

// Message, tek bir flash mesajının yapısını tutar.
type Message struct {
	Type    string `json:"type"` // "success", "error", "warning", "info"
	Title   string `json:"title"`
	Message string `json:"message"`
}

// getStore, context'ten session store'u güvenli bir şekilde alır.
func getStore(c *fiber.Ctx) *session.Store {
	store, ok := c.Locals("session_store").(*session.Store)
	if !ok {
		// Bu durum, session middleware'i çalışmıyorsa oluşur.
		// Paniklemek yerine boş dönebiliriz, ama geliştirme aşamasında loglamak iyi olur.
		return nil
	}
	return store
}

// addMessage, session'a yeni bir flash mesajı ekler.
func addMessage(c *fiber.Ctx, msgType, title, message string) {
	store := getStore(c)
	if store == nil {
		return
	}
	sess, _ := store.Get(c)
	
	// Session'daki mevcut mesajları al
	var messages []Message
	if raw := sess.Get(sessionKey); raw != nil {
		json.Unmarshal(raw.([]byte), &messages)
	}
	
	// Yeni mesajı ekle
	messages = append(messages, Message{Type: msgType, Title: title, Message: message})
	
	// JSON'a çevirip session'a geri yaz
	rawBytes, _ := json.Marshal(messages)
	sess.Set(sessionKey, rawBytes)
	sess.Save()
}

// Success, bir başarı mesajı flash'ler.
func Success(c *fiber.Ctx, message string) {
	addMessage(c, "success", "Success!", message)
}

// Error, bir hata mesajı flash'ler.
func Error(c *fiber.Ctx, message string) {
	addMessage(c, "error", "Error!", message)
}

// Info, bir bilgi mesajı flash'ler.
func Info(c *fiber.Ctx, message string) {
	addMessage(c, "info", "Info", message)
}

// Warning, bir uyarı mesajı flash'ler.
func Warning(c *fiber.Ctx, message string) {
	addMessage(c, "warning", "Warning", message)
}

// GetMessages, session'daki tüm flash mesajları alır ve session'ı TEMİZLER.
// Bu fonksiyon, view.Render tarafından çağrılmalıdır.
func GetMessages(c *fiber.Ctx) []Message {
	store := getStore(c)
	if store == nil {
		return nil
	}
	sess, _ := store.Get(c)
	
	raw := sess.Get(sessionKey)
	if raw == nil {
		return nil
	}

	// Mesajları okuduktan sonra session'dan sil.
	sess.Delete(sessionKey)
	sess.Save()
	
	var messages []Message
	json.Unmarshal(raw.([]byte), &messages)
	
	return messages
}
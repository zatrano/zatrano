package handlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/zatrano/zatrano/app/services"
	"github.com/zatrano/zatrano/internal/zatrano/flash"
	"github.com/zatrano/zatrano/internal/zatrano/form"
	"github.com/zatrano/zatrano/internal/zatrano/query" // Merkezi Query Parser paketimiz
	"github.com/zatrano/zatrano/internal/zatrano/view"
)

// BaseHandler, SSR için standart ve gelişmiş CRUD işlemlerini yönetir.
type BaseHandler[T any, CForm any, UForm any] struct {
	// Servis arayüzü artık jenerik form tiplerini de biliyor.
	Service      services.IBaseService[T, CForm, UForm]
	ResourceName string
}

func NewBaseHandler[T any, CForm any, UForm any](service services.IBaseService[T, CForm, UForm], resourceName string) *BaseHandler[T, CForm, UForm] {
	return &BaseHandler[T, CForm, UForm]{Service: service, ResourceName: resourceName}
}

// ===================================================================
// CRUD METODLARI
// ===================================================================

// Index, merkezi query parser'ı kullanarak gelişmiş listeleme yapar.
func (h *BaseHandler[T, CForm, UForm]) Index(c *fiber.Ctx) error {
	// 1. Gelen HTTP isteğinin query parametrelerini merkezi parser ile analiz et.
	// Bu, ?filter[...], ?sort=..., ?page[...]=... gibi tüm formatları anlar.
	parser := query.NewParser()
	q := parser.Parse(c)

	// 2. Oluşturulan Query nesnesini doğrudan servis katmanına gönder.
	// Handler artık filtreleme veya sıralama mantığını bilmiyor.
	records, total, err := h.Service.Find(q)
	if err != nil {
		flash.Error(c, "Could not fetch records: "+err.Error())
		records = []T{}
	}

	// 3. Paginasyon verilerini view'a göndermek için hazırla.
	pagination := fiber.Map{
		"Total":      total,
		"Page":       q.Pagination.Page,
		"Limit":      q.Pagination.Limit,
		"TotalPages": (total + int64(q.Pagination.Limit) - 1) / int64(q.Pagination.Limit),
		"HasPrev":    q.Pagination.Page > 1,
		"HasNext":    int64(q.Pagination.Page*q.Pagination.Limit) < total,
	}

	// 4. View'ı render et.
	return view.Render(c, fmt.Sprintf("pages/%s/index", h.ResourceName), fiber.Map{
		"Title":      fmt.Sprintf("%s List", strings.Title(h.ResourceName)),
		"Records":    records,
		"Pagination": pagination,
		"Filters":    c.Queries(), // Arama formunu tekrar doldurmak için ham query'yi gönder.
	})
}

// Create, oluşturma formunu gösterir.
func (h *BaseHandler[T, CForm, UForm]) Create(c *fiber.Ctx) error {
	formBuilder := form.New(c, form.Config{
		Action: fmt.Sprintf("/%s", h.ResourceName),
		Method: "POST",
	})
	return view.Render(c, fmt.Sprintf("pages/%s/create", h.ResourceName), fiber.Map{
		"Title": "Create " + strings.Title(h.ResourceName),
		"Form":  formBuilder,
	})
}

// Store, yeni bir kaydı doğrular ve oluşturur.
func (h *BaseHandler[T, CForm, UForm]) Store(c *fiber.Ctx) error {
	var formModel CForm
	if err := c.BodyParser(&formModel); err != nil {
		flash.Error(c, "Invalid form data.")
		return c.RedirectBack(c.Path())
	}

	errors, err := form.Validate(&formModel)
	if err != nil {
		return form.FlashAndRedirect(c, errors, c.Path())
	}

	if _, err := h.Service.CreateFromForm(&formModel); err != nil {
		flash.Error(c, "Could not create record: "+err.Error())
		return form.FlashAndRedirect(c, nil, c.Path())
	}

	flash.Success(c, strings.Title(h.ResourceName)+" created successfully.")
	return c.Redirect("/" + h.ResourceName)
}

// Show, tekil kayıt görüntüleme sayfasını render eder.
func (h *BaseHandler[T, CForm, UForm]) Show(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	record, err := h.Service.GetByID(uint(id))
	if err != nil {
		return fiber.ErrNotFound
	}
	return view.Render(c, fmt.Sprintf("pages/%s/show", h.ResourceName), fiber.Map{
		"Title":  "View " + strings.Title(h.ResourceName),
		"Record": record,
	})
}

// Edit, güncelleme formunu gösterir.
func (h *BaseHandler[T, CForm, UForm]) Edit(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	record, err := h.Service.GetByID(uint(id))
	if err != nil {
		return fiber.ErrNotFound
	}

	formBuilder := form.New(c, form.Config{
		Action: fmt.Sprintf("/%s/%d", h.ResourceName, id),
		Method: "PUT",
		Model:  record,
	})

	return view.Render(c, fmt.Sprintf("pages/%s/edit", h.ResourceName), fiber.Map{
		"Title":  "Edit " + strings.Title(h.ResourceName),
		"Form":   formBuilder,
		"Record": record,
	})
}

// Update, bir kaydı doğrular ve günceller.
func (h *BaseHandler[T, CForm, UForm]) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	var formModel UForm
	if err := c.BodyParser(&formModel); err != nil {
		flash.Error(c, "Invalid form data.")
		return c.RedirectBack(c.Path())
	}

	errors, err := form.Validate(&formModel)
	if err != nil {
		return form.FlashAndRedirect(c, errors, c.Path())
	}

	if _, err := h.Service.UpdateFromForm(uint(id), &formModel); err != nil {
		flash.Error(c, "Could not update record: "+err.Error())
		return form.FlashAndRedirect(c, nil, c.Path())
	}

	flash.Success(c, strings.Title(h.ResourceName)+" updated successfully.")
	return c.Redirect(fmt.Sprintf("/%s/%d", h.ResourceName, id))
}

// Destroy, bir kaydı siler.
func (h *BaseHandler[T, CForm, UForm]) Destroy(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		flash.Error(c, "Could not delete record: "+err.Error())
		return c.RedirectBack(fmt.Sprintf("/%s", h.ResourceName))
	}

	flash.Success(c, strings.Title(h.ResourceName)+" deleted successfully.")
	return c.Redirect("/" + h.ResourceName)
}

// ===================================================================
// HELPER METODLARI ARTIK GEREKLİ DEĞİL
// ===================================================================

// func (h *BaseHandler[T, CForm, UForm]) parseFilters(c *fiber.Ctx) map[string]interface{} { ... } // SİLİNDİ
// func (h *BaseHandler[T, CForm, UForm]) parseOrderBy(c *fiber.Ctx, defaultSort string) string { ... } // SİLİNDİ

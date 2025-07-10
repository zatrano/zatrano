package handlers

import (
	"fmt"
	"strings"
	
	"github.com/gofiber/fiber/v2"
	"github.com/zatrano/zatrano/app/services"
	"github.com/zatrano/zatrano/internal/zatrano/flash"
	"github.com/zatrano/zatrano/internal/zatrano/form" // DAHİLİ form paketimiz
	"github.com/zatrano/zatrano/internal/zatrano/view"
)

// BaseHandler, SSR için standart CRUD işlemlerini yönetir.
type BaseHandler[T any, CForm any, UForm any] struct {
	// Servis arayüzü artık form tiplerini de biliyor.
	Service      services.IBaseService[T, CForm, UForm]
	ResourceName string // "users", "products" vb. (küçük harfle)
}

func NewBaseHandler[T any, CForm any, UForm any](service services.IBaseService[T, CForm, UForm], resourceName string) *BaseHandler[T, CForm, UForm] {
	return &BaseHandler[T, CForm, UForm]{Service: service, ResourceName: resourceName}
}

// Index, listeleme sayfasını render eder.
func (h *BaseHandler[T, CForm, UForm]) Index(c *fiber.Ctx) error {
	records, _, _ := h.Service.GetAll(1, 100)
	return view.Render(c, fmt.Sprintf("pages/%s/index", h.ResourceName), fiber.Map{
		"Title":   fmt.Sprintf("%s List", strings.Title(h.ResourceName)),
		"Records": records,
	})
}

// Create, oluşturma formunu gösterir.
func (h *BaseHandler[T, CForm, UForm]) Create(c *fiber.Ctx) error {
	// Artık CSRF token'ı veya session ile uğraşmıyoruz.
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
		return c.RedirectBack(fmt.Sprintf("/%s/create", h.ResourceName))
	}

	// Dahili form paketimizin Validate fonksiyonunu kullanıyoruz.
	errors, err := form.Validate(&formModel)
	if err != nil {
		// Hata varsa, merkezi FlashAndRedirect helper'ını kullanıyoruz.
		return form.FlashAndRedirect(c, errors, fmt.Sprintf("/%s/create", h.ResourceName))
	}

	// Handler, veri dönüşümü yapmıyor.
	// Doğrulanmış form struct'ını doğrudan servisin doğru metoduna gönderiyor.
	if _, err := h.Service.CreateFromForm(&formModel); err != nil {
		flash.Error(c, "Could not create record: "+err.Error())
		return form.FlashAndRedirect(c, nil, fmt.Sprintf("/%s/create", h.ResourceName))
	}

	flash.Success(c, strings.Title(h.ResourceName)+" created successfully.")
	return c.Redirect("/" + h.ResourceName)
}

// Show, tekil kayıt görüntüleme sayfasını render eder.
func (h *BaseHandler[T, CForm, UForm]) Show(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
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
	id, _ := c.ParamsInt("id")
	record, err := h.Service.GetByID(uint(id))
	if err != nil {
		return fiber.ErrNotFound
	}

	// `form.New` artık modeli alıp `oldInput` mantığını kendi içinde yönetiyor.
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
	id, _ := c.ParamsInt("id")
	var formModel UForm
	if err := c.BodyParser(&formModel); err != nil {
		flash.Error(c, "Invalid form data.")
		return c.RedirectBack(fmt.Sprintf("/%s/%d/edit", h.ResourceName, id))
	}

	errors, err := form.Validate(&formModel)
	if err != nil {
		return form.FlashAndRedirect(c, errors, fmt.Sprintf("/%s/%d/edit", h.ResourceName, id))
	}
	
	// Handler, doğrulanmış formu doğrudan servise gönderiyor.
	if _, err := h.Service.UpdateFromForm(uint(id), &formModel); err != nil {
		flash.Error(c, "Could not update record: "+err.Error())
		return form.FlashAndRedirect(c, nil, fmt.Sprintf("/%s/%d/edit", h.ResourceName, id))
	}

	flash.Success(c, strings.Title(h.ResourceName)+" updated successfully.")
	return c.Redirect(fmt.Sprintf("/%s/%d", h.ResourceName, id))
}

// Destroy, bir kaydı siler.
func (h *BaseHandler[T, CForm, UForm]) Destroy(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := h.Service.Delete(uint(id)); err != nil {
		flash.Error(c, "Could not delete record: "+err.Error())
		return c.RedirectBack(fmt.Sprintf("/%s", h.ResourceName))
	}
	flash.Success(c, strings.Title(h.ResourceName)+" deleted successfully.")
	return c.Redirect("/" + h.ResourceName)
}
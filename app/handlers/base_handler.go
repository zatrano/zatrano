package handlers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	builder "github.com/zatrano/form-builder"
	"github.com/zatrano/zatrano/app/services"
	"github.com/zatrano/zatrano/internal/zatrano/flash"
	"github.com/zatrano/zatrano/internal/zatrano/view"
)

// BaseHandler, SSR için standart CRUD işlemlerini yönetir.
type BaseHandler[T any, CForm any, UForm any] struct {
	Service      services.IBaseService[T]
	ResourceName string // "users", "products" vb. (küçük harfle)
}

func NewBaseHandler[T any, CForm any, UForm any](service services.IBaseService[T], resourceName string) *BaseHandler[T, CForm, UForm] {
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
	form := builder.New(builder.Config{
		Action:    fmt.Sprintf("/%s", h.ResourceName),
		Method:    "POST",
		CSRFToken: c.Locals("csrf").(string),
	})
	return view.Render(c, fmt.Sprintf("pages/%s/create", h.ResourceName), fiber.Map{
		"Title": "Create " + strings.Title(h.ResourceName),
		"Form":  form,
	})
}

// Store, yeni bir kaydı doğrular ve oluşturur.
func (h *BaseHandler[T, CForm, UForm]) Store(c *fiber.Ctx) error {
	var formModel CForm
	if err := c.BodyParser(&formModel); err != nil {
		flash.Error(c, "Invalid form data.")
		return c.RedirectBack(fmt.Sprintf("/%s/create", h.ResourceName))
	}

	errors, err := builder.Validate(&formModel)
	if err != nil {
		store := c.Locals("session_store").(*session.Store)
		sess, _ := store.Get(c)
		formValues, _ := c.FormValues()
		sess.Set("errors", errors)
		sess.Set("old_input", builder.ParseToOldInput(formValues))
		sess.Save()
		flash.Error(c, "Please check the form for errors.")
		return c.Redirect(fmt.Sprintf("/%s/create", h.ResourceName))
	}

	var dbModel T
	// TODO: formModel'i dbModel'e map'leyen bir helper fonksiyonu gerekecek.
	// Örneğin: mappers.Map(&formModel, &dbModel)
	// Şimdilik basit bir varsayım yapalım:
	reflect.ValueOf(&dbModel).Elem().FieldByName("Name").Set(reflect.ValueOf(formModel).Elem().FieldByName("Name"))
	reflect.ValueOf(&dbModel).Elem().FieldByName("Email").Set(reflect.ValueOf(formModel).Elem().FieldByName("Email"))


	if err := h.Service.Create(&dbModel); err != nil {
		flash.Error(c, "Could not create record: "+err.Error())
		return c.Redirect(fmt.Sprintf("/%s/create", h.ResourceName))
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

	form := builder.New(builder.Config{
		Action:    fmt.Sprintf("/%s/%d", h.ResourceName, id),
		Method:    "PUT",
		CSRFToken: c.Locals("csrf").(string),
		Model:     record,
	})

	return view.Render(c, fmt.Sprintf("pages/%s/edit", h.ResourceName), fiber.Map{
		"Title":  "Edit " + strings.Title(h.ResourceName),
		"Form":   form,
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

	errors, err := builder.Validate(&formModel)
	if err != nil {
		store := c.Locals("session_store").(*session.Store)
		sess, _ := store.Get(c)
		formValues, _ := c.FormValues()
		sess.Set("errors", errors)
		sess.Set("old_input", builder.ParseToOldInput(formValues))
		sess.Save()
		flash.Error(c, "Please check the form for errors.")
		return c.Redirect(fmt.Sprintf("/%s/%d/edit", h.ResourceName, id))
	}
	
	var dbModel T
	// TODO: formModel -> dbModel mapping
	reflect.ValueOf(&dbModel).Elem().FieldByName("Name").Set(reflect.ValueOf(formModel).Elem().FieldByName("Name"))
	reflect.ValueOf(&dbModel).Elem().FieldByName("Email").Set(reflect.ValueOf(formModel).Elem().FieldByName("Email"))


	if _, err := h.Service.Update(uint(id), &dbModel); err != nil {
		flash.Error(c, "Could not update record: "+err.Error())
		return c.Redirect(fmt.Sprintf("/%s/%d/edit", h.ResourceName, id))
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
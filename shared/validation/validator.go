package validation

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/Manuelda10/pcshop-backend/shared/http"

	"github.com/go-playground/validator/v10"
)

// MessageFunc genera el mensaje de error a partir del FieldError del validator.
type MessageFunc func(e validator.FieldError) string

// Rule representa una validación custom a nivel de campo.
type Rule struct {
	Tag     string
	Fn      validator.Func
	Message MessageFunc
}

// StructRule representa una validación cross-field a nivel de struct.
type StructRule struct {
	Type   any
	Fn     validator.StructLevelFunc
	Errors map[string]string // tag → mensaje
}

// Validator encapsula el engine de go-playground/validator con reglas y mensajes registrados.
type Validator struct {
	engine   *validator.Validate
	messages map[string]MessageFunc
	defaults map[string]MessageFunc
	mu       sync.RWMutex
}

// New crea un Validator limpio con los defaults en español.
func New() *Validator {
	v := &Validator{
		engine:   validator.New(),
		messages: make(map[string]MessageFunc),
		defaults: defaultMessages(),
	}

	// Usar json tags como nombre de campo
	v.engine.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})

	return v
}

// RegisterRule registra una validación de campo con su mensaje.
func (v *Validator) RegisterRule(r Rule) error {
	if err := v.engine.RegisterValidation(r.Tag, r.Fn); err != nil {
		return fmt.Errorf("registering rule %q: %w", r.Tag, err)
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	if r.Message != nil {
		v.messages[r.Tag] = r.Message
	}
	return nil
}

// RegisterStructRule registra una validación cross-field con sus mensajes.
func (v *Validator) RegisterStructRule(r StructRule) {
	v.engine.RegisterStructValidation(r.Fn, r.Type)

	v.mu.Lock()
	defer v.mu.Unlock()

	for tag, msg := range r.Errors {
		captured := msg
		v.messages[tag] = func(_ validator.FieldError) string {
			return captured
		}
	}
}

// OverrideMessage permite sobreescribir mensajes default o custom después de la creación.
func (v *Validator) OverrideMessage(tag string, fn MessageFunc) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.messages[tag] = fn
}

// Validate valida un struct y retorna los errores como []FieldError.
// Retorna nil si no hay errores.
func (v *Validator) Validate(req any) []http.FieldError {
	err := v.engine.Struct(req)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []http.FieldError{{Field: "_", Message: "Error de validación desconocido"}}
	}

	errs := make([]http.FieldError, 0, len(validationErrors))
	for _, e := range validationErrors {
		errs = append(errs, http.FieldError{
			Field:   e.Field(),
			Message: v.resolveMessage(e),
		})
	}
	return errs
}

func (v *Validator) resolveMessage(e validator.FieldError) string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	// 1. Custom registrado
	if fn, ok := v.messages[e.Tag()]; ok {
		return fn(e)
	}

	// 2. Default conocido
	if fn, ok := v.defaults[e.Tag()]; ok {
		return fn(e)
	}

	// 3. Fallback genérico
	return "Valor inválido"
}

func defaultMessages() map[string]MessageFunc {
	return map[string]MessageFunc{
		"required": func(_ validator.FieldError) string {
			return "Este campo es obligatorio"
		},
		"email": func(_ validator.FieldError) string {
			return "El email no tiene un formato válido"
		},
		"min": func(e validator.FieldError) string {
			return fmt.Sprintf("Mínimo %s caracteres", e.Param())
		},
		"max": func(e validator.FieldError) string {
			return fmt.Sprintf("Máximo %s caracteres", e.Param())
		},
		"len": func(e validator.FieldError) string {
			return fmt.Sprintf("Debe tener exactamente %s caracteres", e.Param())
		},
		"oneof": func(e validator.FieldError) string {
			return fmt.Sprintf("Debe ser uno de: %s", e.Param())
		},
		"url": func(_ validator.FieldError) string {
			return "La URL no tiene un formato válido"
		},
		"uuid": func(_ validator.FieldError) string {
			return "Debe ser un UUID válido"
		},
		"numeric": func(_ validator.FieldError) string {
			return "Debe ser un valor numérico"
		},
		"alpha": func(_ validator.FieldError) string {
			return "Solo se permiten letras"
		},
		"alphanum": func(_ validator.FieldError) string {
			return "Solo se permiten letras y números"
		},
		"gte": func(e validator.FieldError) string {
			return fmt.Sprintf("Debe ser mayor o igual a %s", e.Param())
		},
		"lte": func(e validator.FieldError) string {
			return fmt.Sprintf("Debe ser menor o igual a %s", e.Param())
		},
		"eqfield": func(e validator.FieldError) string {
			return fmt.Sprintf("Debe ser igual al campo %s", e.Param())
		},
	}
}

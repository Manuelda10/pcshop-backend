package rules

import (
	"unicode"

	"github.com/Manuelda10/pcshop-backend/shared/validation"
	"github.com/go-playground/validator/v10"
)

// StrongPassword valida: 8+ chars, 1 mayúscula, 1 minúscula, 1 dígito, 1 especial.
func StrongPassword() validation.Rule {
	return validation.Rule{
		Tag: "strong_password",
		Fn: func(fl validator.FieldLevel) bool {
			pw := fl.Field().String()
			if len(pw) < 8 {
				return false
			}
			var upper, lower, digit, special bool
			for _, c := range pw {
				switch {
				case unicode.IsUpper(c):
					upper = true
				case unicode.IsLower(c):
					lower = true
				case unicode.IsDigit(c):
					digit = true
				case unicode.IsPunct(c) || unicode.IsSymbol(c):
					special = true
				}
			}
			return upper && lower && digit && special
		},
		Message: func(_ validator.FieldError) string {
			return "Debe tener al menos 8 caracteres, una mayúscula, una minúscula, un número y un carácter especial"
		},
	}
}

// PhoneNumberPE valida un teléfono peruano: exactamente 9 dígitos.
func PhoneNumberPE() validation.Rule {
	return validation.Rule{
		Tag: "phone_pe",
		Fn: func(fl validator.FieldLevel) bool {
			s := fl.Field().String()
			if len(s) != 9 {
				return false
			}
			for _, c := range s {
				if !unicode.IsDigit(c) {
					return false
				}
			}
			return true
		},
		Message: func(_ validator.FieldError) string {
			return "Debe ser un número de teléfono válido (9 dígitos)"
		},
	}
}

// DocumentNumber valida: solo dígitos, entre 8 y 12 caracteres.
func DocumentNumber() validation.Rule {
	return validation.Rule{
		Tag: "document_number",
		Fn: func(fl validator.FieldLevel) bool {
			s := fl.Field().String()
			if len(s) < 8 || len(s) > 12 {
				return false
			}
			for _, c := range s {
				if !unicode.IsDigit(c) {
					return false
				}
			}
			return true
		},
		Message: func(_ validator.FieldError) string {
			return "Debe contener entre 8 y 12 dígitos numéricos"
		},
	}
}

// DocumentType valida que sea uno de los tipos permitidos.
func DocumentType() validation.Rule {
	allowed := map[string]bool{
		"DNI": true,
		"CE":  true,
		"RUC": true,
		"PP":  true,
	}
	return validation.Rule{
		Tag: "document_type",
		Fn: func(fl validator.FieldLevel) bool {
			return allowed[fl.Field().String()]
		},
		Message: func(_ validator.FieldError) string {
			return "Debe ser uno de: DNI, CE, RUC, PP"
		},
	}
}

// PrivacyPolicyAccepted valida que el bool sea true.
func PrivacyPolicyAccepted() validation.Rule {
	return validation.Rule{
		Tag: "privacy_accepted",
		Fn: func(fl validator.FieldLevel) bool {
			return fl.Field().Bool()
		},
		Message: func(_ validator.FieldError) string {
			return "Debes aceptar la política de privacidad"
		},
	}
}

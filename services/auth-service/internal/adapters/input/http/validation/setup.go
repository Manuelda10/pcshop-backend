package validation

import (
	"auth-service/internal/adapters/input/http/validation/rules"

	"github.com/Manuelda10/pcshop-backend/shared/validation"
)

func NewAuthValidator() *validation.Validator {
	return validation.Build().
		WithRules(
			rules.StrongPassword(),
			rules.PhoneNumberPE(),
			rules.DocumentNumber(),
			rules.DocumentType(),
			rules.PrivacyPolicyAccepted(),
		).
		MustBuild()
}

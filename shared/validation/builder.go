package validation

// Builder permite configurar un Validator de forma fluida.
//
// Uso:
//
//	v := validation.Build().
//	    WithRules(myRules...).
//	    WithStructRules(myStructRules...).
//	    MustBuild()
type Builder struct {
	rules       []Rule
	structRules []StructRule
	overrides   map[string]MessageFunc
	err         error
}

// Build crea un nuevo Builder.
func Build() *Builder {
	return &Builder{
		overrides: make(map[string]MessageFunc),
	}
}

// WithRules agrega reglas de campo.
func (b *Builder) WithRules(rules ...Rule) *Builder {
	b.rules = append(b.rules, rules...)
	return b
}

// WithStructRules agrega reglas cross-field.
func (b *Builder) WithStructRules(rules ...StructRule) *Builder {
	b.structRules = append(b.structRules, rules...)
	return b
}

// WithMessage sobreescribe el mensaje de un tag específico.
func (b *Builder) WithMessage(tag string, fn MessageFunc) *Builder {
	b.overrides[tag] = fn
	return b
}

// Build construye el Validator. Retorna error si alguna regla falla al registrarse.
func (b *Builder) Build() (*Validator, error) {
	v := New()

	for _, r := range b.rules {
		if err := v.RegisterRule(r); err != nil {
			return nil, err
		}
	}

	for _, r := range b.structRules {
		v.RegisterStructRule(r)
	}

	for tag, fn := range b.overrides {
		v.OverrideMessage(tag, fn)
	}

	return v, nil
}

// MustBuild es como Build pero hace panic si hay error.
// Ideal para init/setup donde un fallo es irrecuperable.
func (b *Builder) MustBuild() *Validator {
	v, err := b.Build()
	if err != nil {
		panic("validation: " + err.Error())
	}
	return v
}

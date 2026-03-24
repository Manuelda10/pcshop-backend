package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"notification-service/internal/core/domain"
)

// subjectMap asocia cada tipo de evento con su subject de email.
var subjectMap = map[domain.NotificationType]string{
	domain.NotificationUserRegistered:       "Verifica tu email",
	domain.NotificationPasswordResetRequest: "Restablecer tu contraseña",
}

// contentMap asocia cada tipo de evento con su template de contenido.
var contentMap = map[domain.NotificationType]string{
	domain.NotificationUserRegistered:       VerifyEmailContent,
	domain.NotificationPasswordResetRequest: ResetPasswordContent,
}

// Renderer implementa output.TemplateRenderer usando html/template de Go.
type Renderer struct {
	appName  string
	layout   *template.Template
	contents map[domain.NotificationType]*template.Template
}

// NewRenderer compila todos los templates al inicializar.
// Si algún template tiene error de sintaxis, falla aquí (fail fast)
// en vez de fallar en runtime cuando llega un evento.
func NewRenderer(appName string) (*Renderer, error) {
	layout, err := template.New("layout").Parse(BaseLayout)
	if err != nil {
		return nil, fmt.Errorf("compiling base layout: %w", err)
	}

	contents := make(map[domain.NotificationType]*template.Template)
	for eventType, content := range contentMap {
		tmpl, err := template.New(string(eventType)).Parse(content)
		if err != nil {
			return nil, fmt.Errorf("compiling template %s: %w", eventType, err)
		}
		contents[eventType] = tmpl
	}

	return &Renderer{
		appName:  appName,
		layout:   layout,
		contents: contents,
	}, nil
}

func (r *Renderer) Render(eventType domain.NotificationType, data map[string]string) (string, string, error) {
	subject, ok := subjectMap[eventType]
	if !ok {
		return "", "", fmt.Errorf("no subject for event type: %s", eventType)
	}

	contentTmpl, ok := r.contents[eventType]
	if !ok {
		return "", "", fmt.Errorf("no template for event type: %s", eventType)
	}

	// 1. Renderizar el contenido específico con los datos del payload
	contentData := toTemplateData(data)
	var contentBuf bytes.Buffer
	if err := contentTmpl.Execute(&contentBuf, contentData); err != nil {
		return "", "", fmt.Errorf("render content: %w", err)
	}

	// 2. Inyectar el contenido renderizado en el layout base
	layoutData := map[string]interface{}{
		"Subject": subject,
		"AppName": r.appName,
		"Year":    time.Now().Year(),
		"Content": template.HTML(contentBuf.String()),
	}

	var htmlBuf bytes.Buffer
	if err := r.layout.Execute(&htmlBuf, layoutData); err != nil {
		return "", "", fmt.Errorf("render layout: %w", err)
	}

	return subject, htmlBuf.String(), nil
}

// toTemplateData mapea las keys del payload (camelCase del JSON)
// a PascalCase para los templates de Go.
func toTemplateData(data map[string]string) map[string]string {
	return map[string]string{
		"FirstName": data["firstName"],
		"Email":     data["email"],
		"OTPCode":   data["otpCode"],
	}
}

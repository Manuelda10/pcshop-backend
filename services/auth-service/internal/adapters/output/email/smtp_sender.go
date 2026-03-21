package email

import (
	"auth-service/config"
	"auth-service/internal/core/ports/output"
	"context"
	"fmt"
	"net/smtp"

	"github.com/Manuelda10/pcshop-backend/shared/logger"
)

type smtpSender struct {
	cfg config.EmailConfig
}

// NewSMTPSender crea una nueva instancia del email sender.
func NewSMTPSender(cfg config.EmailConfig) output.EmailSender {
	return &smtpSender{cfg: cfg}
}

func (s *smtpSender) SendVerificationEmail(ctx context.Context, toEmail, toName, code string) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("email"),
		logger.Operation("send_verification_email"),
		logger.Email(toEmail),
	)

	log.Debug("sending verification email")

	subject := "Verifica tu cuenta en PC Store"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
  <h2>Hola %s 👋</h2>
  <p>Gracias por registrarte en <strong>PC Store</strong>.</p>
  <p>Tu código de verificación es:</p>
  <div style="background: #f4f4f4; padding: 20px; text-align: center; border-radius: 8px; margin: 20px 0;">
    <span style="font-size: 36px; font-weight: bold; letter-spacing: 8px; color: #333;">%s</span>
  </div>
  <p style="color: #666;">Este código expira en <strong>24 horas</strong>.</p>
  <p style="color: #999; font-size: 12px;">Si no creaste esta cuenta, ignora este email.</p>
</body>
</html>`, toName, code)

	if err := s.send(toEmail, toName, subject, body); err != nil {
		log.Error("failed to send verification email", logger.Err(err))
		return err
	}

	log.Info("verification email sent")
	return nil
}

func (s *smtpSender) SendPasswordResetEmail(ctx context.Context, toEmail, toName, code string) error {
	log := logger.FromContext(ctx).With(
		logger.Layer("email"),
		logger.Operation("send_password_reset_email"),
		logger.Email(toEmail),
	)

	log.Debug("sending password reset email")

	subject := "Restablece tu contraseña en PC Store"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
  <h2>Hola %s</h2>
  <p>Recibimos una solicitud para restablecer tu contraseña.</p>
  <p>Tu código es:</p>
  <div style="background: #f4f4f4; padding: 20px; text-align: center; border-radius: 8px; margin: 20px 0;">
    <span style="font-size: 36px; font-weight: bold; letter-spacing: 8px; color: #333;">%s</span>
  </div>
  <p style="color: #666;">Este código expira en <strong>1 hora</strong>.</p>
  <p style="color: #999; font-size: 12px;">Si no solicitaste esto, ignora este email.</p>
</body>
</html>`, toName, code)

	if err := s.send(toEmail, toName, subject, body); err != nil {
		log.Error("failed to send password reset email", logger.Err(err))
		return err
	}

	log.Info("password reset email sent")
	return nil
}

// send es el helper base que construye y envía el mensaje SMTP.
func (s *smtpSender) send(toEmail, toName, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	from := fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.FromAddress)
	to := fmt.Sprintf("%s <%s>", toName, toEmail)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, to, subject, htmlBody,
	)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.FromAddress, []string{toEmail}, []byte(msg))
}

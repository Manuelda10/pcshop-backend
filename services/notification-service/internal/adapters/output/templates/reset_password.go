package templates

// ResetPasswordContent es el contenido del email de reset de contraseña con OTP.
const ResetPasswordContent = `
<h1 style="margin: 0 0 8px; font-size: 22px; font-weight: 700; color: #18181b; line-height: 28px;">
  Restablecer contraseña
</h1>
<p style="margin: 0 0 24px; font-size: 15px; line-height: 24px; color: #52525b;">
  Hola {{.FirstName}}, recibimos una solicitud para restablecer la contraseña de tu cuenta. Usa el siguiente código:
</p>

<!-- OTP Code -->
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="margin: 0 0 24px;">
  <tr>
    <td align="center">
      <div style="display: inline-block; padding: 16px 40px; background-color: #f4f4f5; border-radius: 8px; border: 1px solid #e4e4e7;">
        <span style="font-size: 32px; font-weight: 700; letter-spacing: 8px; color: #18181b; font-family: 'Courier New', monospace;">{{.OTPCode}}</span>
      </div>
    </td>
  </tr>
</table>

<!-- Warning box -->
<table role="presentation" width="100%" cellpadding="0" cellspacing="0">
  <tr>
    <td style="background-color: #fef2f2; border-radius: 8px; padding: 16px 20px; border-left: 3px solid #ef4444;">
      <p style="margin: 0; font-size: 13px; line-height: 20px; color: #991b1b;">
        Este código expira en <strong>10 minutos</strong>. Si no solicitaste este cambio, te recomendamos cambiar tu contraseña inmediatamente por seguridad.
      </p>
    </td>
  </tr>
</table>`

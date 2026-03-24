package templates

// BaseLayout es el layout compartido por todos los emails.
// Usa table-based layout con CSS inline porque los clientes de email
// (Gmail, Outlook, Apple Mail) ignoran <style> blocks.
// Max-width 600px es el estándar de la industria para emails.
const BaseLayout = `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>{{.Subject}}</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f4f4f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; -webkit-font-smoothing: antialiased;">

  <!-- Outer wrapper -->
  <table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="background-color: #f4f4f5;">
    <tr>
      <td align="center" style="padding: 40px 16px;">

        <!-- Logo -->
        <table role="presentation" width="600" cellpadding="0" cellspacing="0" style="max-width: 600px; width: 100%;">
          <tr>
            <td align="center" style="padding-bottom: 24px;">
              <span style="font-size: 24px; font-weight: 700; color: #18181b; letter-spacing: -0.5px;">{{.AppName}}</span>
            </td>
          </tr>
        </table>

        <!-- Main card -->
        <table role="presentation" width="600" cellpadding="0" cellspacing="0" style="max-width: 600px; width: 100%; background-color: #ffffff; border-radius: 12px; overflow: hidden; border: 1px solid #e4e4e7;">

          <!-- Accent bar -->
          <tr>
            <td style="height: 4px; background-color: #18181b;"></td>
          </tr>

          <!-- Content -->
          <tr>
            <td style="padding: 40px 40px 32px;">
              {{.Content}}
            </td>
          </tr>

          <!-- Divider -->
          <tr>
            <td style="padding: 0 40px;">
              <div style="height: 1px; background-color: #e4e4e7;"></div>
            </td>
          </tr>

          <!-- Footer inside card -->
          <tr>
            <td style="padding: 24px 40px 32px;">
              <p style="margin: 0; font-size: 13px; line-height: 20px; color: #a1a1aa;">
                Este email fue enviado por {{.AppName}}. Si no solicitaste esta acción, puedes ignorar este mensaje de forma segura.
              </p>
            </td>
          </tr>
        </table>

        <!-- Footer outside card -->
        <table role="presentation" width="600" cellpadding="0" cellspacing="0" style="max-width: 600px; width: 100%;">
          <tr>
            <td align="center" style="padding: 24px 0;">
              <p style="margin: 0; font-size: 12px; color: #a1a1aa;">
                © {{.Year}} {{.AppName}}. Todos los derechos reservados.
              </p>
            </td>
          </tr>
        </table>

      </td>
    </tr>
  </table>

</body>
</html>`

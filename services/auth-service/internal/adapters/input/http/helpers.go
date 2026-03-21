package http

import (
	"context"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

// googleUserInfo representa el perfil básico retornado por la API de Google.
type googleUserInfo struct {
	Sub           string `json:"sub"` // ID único de Google
	Email         string `json:"email"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
}

// getGoogleUserInfo intercambia el oauth2 token por el perfil del usuario.
func getGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*googleUserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo googleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// generateOAuthState genera un state aleatorio para prevenir ataques CSRF en OAuth2.
func generateOAuthState() string {
	b := make([]byte, 16)
	_, _ = crand.Read(b)
	return fmt.Sprintf("%x", b)
}

// extractUserIDFromToken extrae el user_id del JWT sin verificar la firma.
// Solo para uso interno después de un login exitoso donde ya confiamos en el token.
// La verificación real ocurre en el middleware JWT.
func extractUserIDFromToken(tokenString string) string {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return ""
	}

	claims := jwt.MapClaims{}
	parser := jwt.NewParser()
	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return ""
	}

	userID, _ := claims["user_id"].(string)
	return userID
}

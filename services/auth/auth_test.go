package auth

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAuth_GenerateToken(t *testing.T) {
	vp := viper.New()
	vp.Set("app.auth.secret", "secret")
	auth := NewAuth(vp, zap.NewExample())
	token, err := auth.GenerateToken("123-456-789-0", "auth.test@gmail.com", map[string]interface{}{
		"group": "admin",
	})
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
	}

	_, err = auth.ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken() error = %v", err)
	}

	claims, err := auth.GetClaims(token)
	if err != nil {
		t.Errorf("GetClaims() error = %v", err)
	}

	if claims.Email != "auth.test@gmail.com" {
		t.Errorf("we expect auth.test@gmail.com got %v", claims.Email)
	}
	if int64(claims.Exp.(float64)) > time.Now().Add(time.Hour*25).Unix() ||
		int64(claims.Exp.(float64)) < time.Now().Add(time.Hour*23).Unix() {
		t.Errorf("we expect expiration ok go %v", claims.Exp)
	}
	if claims.Sub != "123-456-789-0" {
		t.Errorf("we expect 123-456-789-0 got %v", claims.Sub)
	}
	if claims.Data["group"] != "admin" {
		t.Errorf("we expect admin got %v", claims.Data)
	}

}

func TestAuth_WithDosNotExpire(t *testing.T) {
	vp := viper.New()
	vp.Set("app.auth.secret", "secret")
	auth := NewAuth(vp, zap.NewExample())
	token, err := auth.GenerateToken("123-456-789-0", "auth.test@gmail.com", map[string]interface{}{
		"group": "admin",
	}, WithDoesNotExpireOption())
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
	}

	_, err = auth.ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken() error = %v", err)
	}

	claims, err := auth.GetClaims(token)
	if err != nil {
		t.Errorf("GetClaims() error = %v", err)
	}

	if claims.Exp != nil {
		t.Errorf("we expect nil got %v", claims.Exp)
	}

}

func TestAuth_SetExpire(t *testing.T) {
	vp := viper.New()
	vp.Set("app.auth.secret", "secret")
	auth := NewAuth(vp, zap.NewExample())
	token, err := auth.GenerateToken("123-456-789-0", "auth.test@gmail.com", map[string]interface{}{
		"group": "admin",
	}, WithExpirationOption(time.Now().Add(time.Hour*1)))
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
	}

	_, err = auth.ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken() error = %v", err)
	}

	claims, err := auth.GetClaims(token)
	if err != nil {
		t.Errorf("GetClaims() error = %v", err)
	}

	if int64(claims.Exp.(float64)) > time.Now().Add(time.Hour*2).Unix() ||
		int64(claims.Exp.(float64)) < time.Now().Add(time.Hour*1).Unix() {
		t.Errorf("we expect expiration ok go %v", claims.Exp)
	}
}

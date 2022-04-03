package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret string = "secret"

const (
	Issuer    = "KoR"
	SignAlg   = "HS256"
	Subject   = "authorization"
	ValidTerm = 30 * 24 * time.Hour
)

type AdditionalClaims struct {
	UserID uint64
}

type CustomClaims struct {
	*jwt.RegisteredClaims
	*AdditionalClaims
}

type JWTToken struct {
	ID    string
	Token string
}

type JWTRawToken struct {
	ID string
	*AdditionalClaims
}

func Init() {
	if secretEnv := os.Getenv("JWT_SECRET"); secretEnv != "" {
		jwtSecret = secretEnv
	}
}

func Generate(uuid string, audience string, userID uint64) (*JWTToken, error) {
	t, now := jwt.New(jwt.GetSigningMethod(SignAlg)), time.Now()
	t.Claims = &CustomClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ValidTerm)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    Issuer,
			Subject:   Subject,
			ID:        uuid,
			Audience:  jwt.ClaimStrings{audience},
		},
		AdditionalClaims: &AdditionalClaims{
			UserID: userID,
		},
	}
	token, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}
	return &JWTToken{ID: uuid, Token: token}, nil
}

func Verify(tokenString string) (*JWTRawToken, error) {
	token, err := parse(tokenString)
	if err != nil {
		return nil, err
	}
	if token.Method.Alg() != SignAlg {
		return nil, errors.New("Invalid Algorithm")
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("Invaid Claims")
	}
	now := time.Now()
	if claims.NotBefore.After(now) {
		return nil, errors.New("Not Started Token")
	}
	if claims.ExpiresAt.Before(now) {
		return nil, errors.New("Expired Token")
	}
	return &JWTRawToken{
		ID: claims.ID,
		AdditionalClaims: claims.AdditionalClaims,
	}, nil
}

func parse(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
}

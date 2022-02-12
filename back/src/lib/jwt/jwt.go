package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/SongCastle/KoR/lib/random"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret string = "secret"

const (
	Issuer    = "KoR"
	SignAlg   = "HS256"
	Subject   = "authorization"
	UUIDLen   = 32
	ValidTerm = 30 * 24 * time.Hour
)

type CustomClaims struct {
	*jwt.RegisteredClaims
	UserID uint64
}

func Init() {
	if secretEnv := os.Getenv("JWT_SECRET"); secretEnv != "" {
		jwtSecret = secretEnv
	}
}

func Generate(userID uint64, audience string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod(SignAlg))
	now := time.Now()
	t.Claims = &CustomClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ValidTerm)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    Issuer,
			Subject:   Subject,
			ID:        random.Generate(UUIDLen),
			Audience:  jwt.ClaimStrings{audience},
		},
		UserID: userID,
	}
	return t.SignedString([]byte(jwtSecret))
}

func Validate(tokenString string, userID uint64) (bool, error) {
	token, err := parse(tokenString)
	if err != nil {
		return false, err
	}
	if token.Method.Alg() != SignAlg {
		return false, errors.New("Invalid Algorithm")
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return false, errors.New("Invaid Claims")
	}
	now := time.Now()
	if claims.NotBefore.After(now) {
		return false, errors.New("Not Started Token")
	}
	if claims.ExpiresAt.Before(now) {
		return false, errors.New("Expired Token")
	}
	return claims.UserID == userID, nil
}

func parse(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &CustomClaims{},  func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
}

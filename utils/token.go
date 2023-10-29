package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(payload interface{}, privateKey string, ttl time.Duration) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("Could not decode private key: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("Could not parse private key: %w", err)
	}

	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"sub": payload,
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Could not create token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("Could not decode public key: %w", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("Could not parse public key: %w", err)
	}
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	fmt.Println("claims", claims)
	if err != nil {
		return nil, fmt.Errorf("Could not parse token: %w", err)
	}
	return claims["sub"], nil
}

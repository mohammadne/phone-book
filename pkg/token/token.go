package token

import (
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token interface {
	CreateTokenString(data any) (string, error)
	ExtractTokenData(tokenString string, data any) error
}

type token struct {
	privateEd25519Key crypto.PrivateKey
	publicEd25519Key  crypto.PublicKey
	expiration        time.Duration
}

func New(cfg *Config) (Token, error) {
	token := &token{}
	var err error

	privatePemKey := []byte(cfg.PrivatePem)
	token.privateEd25519Key, err = jwt.ParseEdPrivateKeyFromPEM(privatePemKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Ed25519 private key: %v", err)
	}

	publicPemKey := []byte(cfg.PublicPem)
	token.publicEd25519Key, err = jwt.ParseEdPublicKeyFromPEM(publicPemKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Ed25519 public key: %v", err)
	}

	token.expiration = cfg.Expiration

	return token, nil
}

type Payload struct {
	Data []byte `json:"data"`
	jwt.RegisteredClaims
}

func (token *token) CreateTokenString(data any) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		errStr := fmt.Sprintf("error marshal data: %v", err)
		return "", errors.New(errStr)
	}

	expierdAt := jwt.NewNumericDate(time.Now().Add(token.expiration))
	registeredClaim := jwt.RegisteredClaims{ExpiresAt: expierdAt}
	payload := &Payload{dataBytes, registeredClaim}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, payload)
	return jwtToken.SignedString(token.privateEd25519Key)
}

const (
	inValidToken        = "invalid token"
	errorMappingPayload = "error mapping the payload"
	errorUnmarshalData  = "error unmarshaling the data"
)

func (token *token) ExtractTokenData(tokenString string, data any) error {
	checkSigningMethod := func(jwtToken *jwt.Token) (any, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("wrong signing method: %v", jwtToken.Header["alg"])
		}
		return token.publicEd25519Key, nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Payload{}, checkSigningMethod)
	if err != nil {
		errStr := fmt.Sprintf("error: %v, token: %s", err, tokenString)
		return errors.New(errStr)
	}

	if !jwtToken.Valid {
		errStr := fmt.Sprintf("%s, token: %v", inValidToken, jwtToken)
		return errors.New(errStr)
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		errStr := fmt.Sprintf("%s: %s, token: %v", inValidToken, errorMappingPayload, jwtToken)
		return errors.New(errStr)
	}

	if err := json.Unmarshal([]byte(payload.Data), data); err != nil {
		errStr := fmt.Sprintf("%s: %s, data: %s", inValidToken, errorUnmarshalData, payload.Data)
		return errors.New(errStr)
	}

	return nil
}

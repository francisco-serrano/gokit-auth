package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

const key = "abc123"

func CreateToken(sessionID string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(sessionID))

	// to base64
	signedMac := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signedMac + "|" + sessionID
}

func ParseToken(token string) (string, error) {
	aux := strings.SplitN(token, "|", 2)
	if len(aux) != 2 {
		return "", fmt.Errorf("error while splitting token")
	}

	encodedSignature := aux[0]
	sessionID := aux[1]

	signature, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		return "", fmt.Errorf("error while decoding signature: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(sessionID))

	if !hmac.Equal(signature, mac.Sum(nil)) {
		return "", fmt.Errorf("not the same signed session and session")
	}

	return sessionID, nil
}
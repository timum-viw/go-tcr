package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"strings"
	"errors"
    "os"
    "fmt"
)

func getHmacKey() []byte {
	keyString := os.Getenv("HMAC_KEY")
    return []byte(keyString)
}

// generate returns a string "<base64_integer>.<hmac_signature>"
func Generate(value uint32) string {
    b := make([]byte, 4)
    binary.BigEndian.PutUint32(b, value)
    encodedValue := base64.StdEncoding.EncodeToString(b)

    key := getHmacKey()
    h := hmac.New(sha256.New, key)
    h.Write(b)
    signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

    return fmt.Sprintf("%s.%s", encodedValue, signature)
}

func Validate(hash string) (uint32, error) {
    parts := strings.Split(hash, ".")
	if len(parts) != 2 {
		return 0, errors.New("invalid input format")
	}

	encodedInt, encodedHMAC := parts[0], parts[1]

	// Decode the integer from base64
	decodedIntBytes, err := base64.StdEncoding.DecodeString(encodedInt)
	if err != nil {
		return 0, errors.New("invalid base64 encoding for integer")
	}

	intValue := binary.BigEndian.Uint32(decodedIntBytes)

    secretKey := getHmacKey()
	// Compute HMAC of the integer string
	h := hmac.New(sha256.New, secretKey)
	h.Write(decodedIntBytes)
	computedHMAC := h.Sum(nil)
	decodedHMAC, err := base64.StdEncoding.DecodeString(encodedHMAC)
	if err != nil {
		return 0, err
	}

	// Compare the computed HMAC with the provided one
	if !hmac.Equal(computedHMAC, decodedHMAC) {
		return 0, errors.New("HMAC validation failed")
	}

	return intValue, nil
}
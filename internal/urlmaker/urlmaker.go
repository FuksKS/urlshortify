package urlmaker

import (
	"crypto/sha256"
	"encoding/hex"
)

func MakeShortUrl(input string) string {
	inputBytes := []byte(input)

	// Вычисление хэша с использованием SHA-256
	hash := sha256.Sum256(inputBytes)

	// Преобразование хэша в строку в шестнадцатеричном формате
	hashString := hex.EncodeToString(hash[:])

	return hashString[:8]
}

package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"math/big"
)

const urlSize = 7
const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	lettersCount = big.NewInt(int64(len(base62Chars)))
)

func generateShortUrl() ([]byte, error) {
	result := make([]byte, urlSize)

	for i := 0; i < urlSize; i++ {
		n, err := rand.Int(rand.Reader, lettersCount)
		if err != nil {
			return nil, err
		}
		result[i] = base62Chars[n.Int64()]
	}
	return result, nil
}

func InsertUrl(url []byte) ([]byte, error) {
	shortURL, err := generateShortUrl()
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(url)
	var longURL []byte
	err = getSetStatement.QueryRow(shortURL, hash[:], url).Scan(&longURL)
	if err != nil {
		return nil, err
	}
	return longURL, nil
}

func GetLongUrl(key []byte) ([]byte, error) {
	var longURL []byte

	err := findStatement.QueryRow(key).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return longURL, nil
}

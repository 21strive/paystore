package pin

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/21strive/redifu"
	"golang.org/x/crypto/argon2"
	"paystore/balance"
	"strings"
)

type Pin struct {
	*redifu.Record
	PIN         string
	Salt        string
	BalanceUUID string
}

func (p *Pin) SetPIN(pin string) error {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(pin), salt, argonTime, memory, threads, keyLen)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	p.PIN = fmt.Sprintf("$argon2id$%s$%s", encodedSalt, encodedHash)
	return nil
}

func (p *Pin) SetBalance(account balance.Balance) {
	p.BalanceUUID = account.GetUUID()
}

func (p *Pin) VerifiyPin(inputPIN string) (bool, error) {
	parts := strings.Split(p.PIN, "$")
	if len(parts) != 4 {
		return false, InvalidHashFormat
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}

	inputHash := argon2.IDKey([]byte(inputPIN), salt, argonTime, memory, threads, keyLen)
	return subtle.ConstantTimeCompare(hash, inputHash) == 1, nil
}

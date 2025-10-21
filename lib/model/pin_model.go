package model

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/21strive/redifu"
	"golang.org/x/crypto/argon2"
	"paystore/lib/def"
	def2 "paystore/lib/def"
	"strings"
)

type Pin struct {
	*redifu.Record
	PIN         string
	Salt        string
	BalanceUUID string
}

func (p *Pin) SetPIN(pin string) error {
	salt := make([]byte, def.SaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(pin), salt, def.ArgonTime, def.Memory, def.Threads, def.KeyLen)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	p.PIN = fmt.Sprintf("$argon2id$%s$%s", encodedSalt, encodedHash)
	return nil
}

func (p *Pin) SetBalance(account Balance) {
	p.BalanceUUID = account.GetUUID()
}

func (p *Pin) VerifiyPin(inputPIN string) (bool, error) {
	parts := strings.Split(p.PIN, "$")
	if len(parts) != 4 {
		return false, def2.InvalidHashFormat
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}

	inputHash := argon2.IDKey([]byte(inputPIN), salt, def.ArgonTime, def.Memory, def.Threads, def.KeyLen)
	return subtle.ConstantTimeCompare(hash, inputHash) == 1, nil
}

package pin

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/21strive/redifu"
	"golang.org/x/crypto/argon2"
	"paystore/balance"
)

type Pin struct {
	*redifu.Record
	PIN         string
	BalanceUUID string
}

func (p *Pin) SetPIN(pin string) error {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(pin), salt, time, memory, threads, keyLen)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	p.PIN = fmt.Sprintf("$argon2id$%s$%s", encodedSalt, encodedHash)
	return nil
}

func (p *Pin) SetBalance(account balance.Balance) {
	p.BalanceUUID = account.GetUUID()
}

func (p *Pin) VerifiyPin(input string) {

}

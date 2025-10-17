package pin

import "errors"

const (
	argonTime = 1
	memory    = 64 * 1024 // 64 MB
	threads   = 4
	keyLen    = 32
	saltLen   = 16
)

var InvalidHashFormat = errors.New("invalid hash format")

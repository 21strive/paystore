package pin

import "errors"

const (
	ArgonTime = 1
	Memory    = 64 * 1024 // 64 MB
	Threads   = 4
	KeyLen    = 32
	SaltLen   = 16
)

var InvalidHashFormat = errors.New("invalid hash format")

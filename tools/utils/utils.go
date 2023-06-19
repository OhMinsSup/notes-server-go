package utils

import (
	"crypto/rand"
	"encoding/base32"
)

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769").WithPadding(base32.NoPadding)

// NewRandomString returns a random string of the given length.
// The resulting entropy will be (5 * length) bits.
func NewRandomString(length int) string {
	data := make([]byte, 1+(length*5/8))
	rand.Read(data)
	return encoding.EncodeToString(data)[:length]
}

func NewBool(b bool) *bool       { return &b }

func NewInt(n int) *int          { return &n }

func NewInt64(n int64) *int64    { return &n }

func NewString(s string) *string { return &s }
package hashing

import "crypto/sha1"

// Sha1 will generate a sha1 hash from b
func Sha1(b []byte) []byte {
	h := sha1.New()
	h.Write(b)
	return h.Sum(nil)
}

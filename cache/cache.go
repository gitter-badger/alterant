package cache

import (
	"crypto/sha1"
	"encoding/base64"
)

func SHAFromBytes(b []byte) string {
	h := sha1.New()
	h.Write(b)
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return sha
}

func SHAFromString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return sha
}

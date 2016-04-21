package hasher

import (
	"crypto/sha1"
	"encoding/base64"
)

func SHA1FromBytes(b []byte) string {
	h := sha1.New()
	h.Write(b)
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return sha
}

func SHA1FromString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return sha
}

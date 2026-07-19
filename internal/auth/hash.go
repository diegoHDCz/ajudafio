package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken returns a deterministic, non-reversible fingerprint of a token
// suitable for storage/lookup (unlike bcrypt, this is fast by design and
// is not meant to resist offline brute-force of a *low-entropy* secret —
// refresh tokens are high-entropy JWTs, so a fast hash is the right tool).
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

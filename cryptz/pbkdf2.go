package cryptz

import (
	"crypto/hmac"
	"hash"
)

// PBKDF2Key derives a cryptographic key from a password and salt using the PBKDF2 algorithm.
//
// Parameters:
//   password - The input password as a byte slice.
//   salt     - A unique salt as a byte slice. Use a cryptographically secure random value.
//   iter     - The number of iterations. A higher value increases computational cost and security.
//              The minimum recommended value is 10,000; values below this may be vulnerable to brute-force attacks.
//   keyLen   - The desired length of the derived key in bytes.
//   h        - A constructor for the underlying hash function (e.g., sha256.New).
//
// Returns:
//   A byte slice containing the derived key of length keyLen.
//
// Security notes:
//   - Always use a unique, random salt for each password.
//   - Choose a sufficiently high iteration count (iter) to slow down brute-force attacks.
//   - Refer to current security guidelines for recommended parameters.
func PBKDF2Key(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

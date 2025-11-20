package cryptz

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	_ "crypto/sha256"
	"encoding/binary"
	"errors"
	"hash"
)

type PBKDF2KeyDeriver struct {
	// Iter is pbkdf2 iteration count, it is recommended to be at least 10_000
	Iter uint32
	// Hash is the hash algorithm used in pbkdf2, default is sha256
	// If the Hash is not available(invalid or not linked), sha256 will be used
	Hash crypto.Hash
}

var (
	pbkdf2ID               = [IDLen]byte{'P', 'B', 'K', 'D', 'F', '2', 0, 1}
	ErrInvalidPBKDF2Header = errors.New("invalid PBKDF2 header")
)

const (
	defHash         = crypto.SHA256
	pbkdf2HeaderLen = IDLen + 5 // ID(8)|Hash(1)|Iter(4)
)

func (P PBKDF2KeyDeriver) ID() [8]byte {
	return pbkdf2ID
}

func (P PBKDF2KeyDeriver) Key(password, salt []byte, keyLen int) []byte {
	h := P.Hash
	if !h.Available() {
		h = defHash
	}
	return PBKDF2Key(password, salt, int(P.Iter), keyLen, h.New)
}

func (P PBKDF2KeyDeriver) Header() []byte {
	// ID(8)|Hash(1)|Iter(4)
	header := [pbkdf2HeaderLen]byte{}
	copy(header[:], pbkdf2ID[:])
	h := P.Hash
	if !h.Available() {
		h = defHash
	}
	header[8] = uint8(h)
	binary.BigEndian.PutUint32(header[IDLen+1:IDLen+5], P.Iter)

	return header[:]
}

func (P PBKDF2KeyDeriver) HeaderLen() int {
	// ID(8)|Hash(1)|Iter(4)
	return pbkdf2HeaderLen
}

func (P PBKDF2KeyDeriver) Restore(deriverHeader []byte) (KeyDeriver, error) {
	if len(deriverHeader) < pbkdf2HeaderLen {
		return nil, ErrInvalidPBKDF2Header
	}

	if !bytes.Equal(deriverHeader[:IDLen], pbkdf2ID[:]) {
		return nil, ErrInvalidPBKDF2Header
	}

	return PBKDF2KeyDeriver{
		Iter: binary.BigEndian.Uint32(deriverHeader[IDLen+1 : IDLen+5]),
		Hash: crypto.Hash(deriverHeader[IDLen]),
	}, nil
}

// PBKDF2Key derives a cryptographic key from a password and salt using the PBKDF2 algorithm.
//
// Parameters:
//
//	password - The input password as a byte slice.
//	salt     - A unique salt as a byte slice. Use a cryptographically secure random value.
//	iter     - The number of iterations. A higher value increases computational cost and security.
//	           The minimum recommended value is 10,000; values below this may be vulnerable to brute-force attacks.
//	keyLen   - The desired length of the derived key in bytes.
//	h        - A constructor for the underlying hash function (e.g., sha256.New).
//
// Returns:
//
//	A byte slice containing the derived key of length keyLen.
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

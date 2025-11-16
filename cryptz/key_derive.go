package cryptz

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	IDLen = 8
)

type KeyDeriver interface {
	// ID returns the unique identifier of the KeyDeriver
	// It must be 8 bytes long
	ID() [IDLen]byte
	// Key derives a key from the given password and salt
	Key(password, salt []byte, keyLen int) []byte
	// Header returns the header bytes representing the KeyDeriver and configuration
	// It must contain the ID as the first 8 bytes
	Header() []byte
	// HeaderLen returns the length of the header
	HeaderLen() int
	// Restore restores the KeyDeriver from the given header bytes
	Restore(deriverHeader []byte) (KeyDeriver, error)
}

type PBKDF2KeyDeriver struct {
	// Iter is pbkdf2 iteration count, it is recommended to be at least 10_000
	Iter uint32
	// Hash is the hash algorithm used in pbkdf2
	Hash HashAlgo
}

var (
	pbkdf2ID               = [IDLen]byte{'P', 'B', 'K', 'D', 'F', '2', 0, 1}
	ErrInvalidPBKDF2Header = errors.New("invalid PBKDF2 header")
)

const (
	pbkdf2HeaderLen = IDLen + 5 // ID(8)|Hash(1)|Iter(4)
)

func (P PBKDF2KeyDeriver) ID() [8]byte {
	return pbkdf2ID
}

func (P PBKDF2KeyDeriver) Key(password, salt []byte, keyLen int) []byte {
	return PBKDF2Key(password, salt, int(P.Iter), keyLen, P.Hash.Factory())
}

func (P PBKDF2KeyDeriver) Header() []byte {
	// ID(8)|Hash(1)|Iter(4)
	header := [pbkdf2HeaderLen]byte{}
	copy(header[:], pbkdf2ID[:])
	header[8] = uint8(P.Hash)
	binary.BigEndian.PutUint32(header[9:13], P.Iter)

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
		Hash: HashAlgo(deriverHeader[IDLen]),
	}, nil
}

package cryptz

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

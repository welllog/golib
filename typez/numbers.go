package typez

type Int64 interface {
	~int64 | ~uint64
}

type Signed interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type Unsigned interface {
	~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uintptr
}

type Float interface {
	~float64 | ~float32
}

type Integer interface {
	Signed | Unsigned
}

type Number interface {
	Integer | Float
}

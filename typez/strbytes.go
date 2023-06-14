package typez

type StrOrBytes interface {
	~string | ~[]byte
}

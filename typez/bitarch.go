package typez

const (
	Is64bitArch = ^uint(0)>>63 == 1
	Is32bitArch = ^uint(0)>>63 == 0
	WordBits    = 32 << (^uint(0) >> 63)
)

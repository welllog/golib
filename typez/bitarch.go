package typez

import "unsafe"

const (
	Is64bitArch = ^uint(0)>>63 == 1
	Is32bitArch = ^uint(0)>>63 == 0
	WordBits    = 32 << (^uint(0) >> 63) // go official translation only support 32/64 bit
	WordBytes   = unsafe.Sizeof(uintptr(0))
)

package buffer

// Constants for buffer sizes in bytes, kilobytes, and megabytes.
// These constants are defined using bitwise left shift operations to represent
// power-of-two sizes, commonly used in computing for buffer sizes.
const (
	SIZE_1B    = 1 << 0  // 1 byte
	SIZE_2B    = 1 << 1  // 2 bytes
	SIZE_4B    = 1 << 2  // 4 bytes
	SIZE_8B    = 1 << 3  // 8 bytes
	SIZE_16B   = 1 << 4  // 16 bytes
	SIZE_32B   = 1 << 5  // 32 bytes
	SIZE_64B   = 1 << 6  // 64 bytes
	SIZE_128B  = 1 << 7  // 128 bytes
	SIZE_256B  = 1 << 8  // 256 bytes
	SIZE_512B  = 1 << 9  // 512 bytes
	SIZE_1KB   = 1 << 10 // 1 kilobyte (1024 bytes)
	SIZE_2KB   = 1 << 11 // 2 kilobytes
	SIZE_4KB   = 1 << 12 // 4 kilobytes
	SIZE_8KB   = 1 << 13 // 8 kilobytes
	SIZE_16KB  = 1 << 14 // 16 kilobytes
	SIZE_32KB  = 1 << 15 // 32 kilobytes
	SIZE_64KB  = 1 << 16 // 64 kilobytes
	SIZE_128KB = 1 << 17 // 128 kilobytes
	SIZE_256KB = 1 << 18 // 256 kilobytes
	SIZE_512KB = 1 << 19 // 512 kilobytes
	SIZE_1MB   = 1 << 20 // 1 megabyte (1024 kilobytes)
)

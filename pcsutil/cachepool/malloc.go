package cachepool

// RawMallocByteSlice allocates a new byte slice.
func RawMallocByteSlice(size int) []byte {
	bytesArray := make([]byte, size)
	return bytesArray
}

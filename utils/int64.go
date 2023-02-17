package utils

import "encoding/binary"

// Int64ToBytes converts a int64 into fixed length bytes for use in store keys.
func Int64ToBytes(id int64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(id))
	return bz
}

// Int64FromBytes converts some fixed length bytes back into a int64.
func Int64FromBytes(bz []byte) int64 {
	return int64(binary.BigEndian.Uint64(bz))
}

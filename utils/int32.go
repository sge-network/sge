package utils

import (
	"encoding/binary"

	"github.com/spf13/cast"
)

// Int32ToBytes converts a int32 into fixed length bytes for use in store keys.
func Int32ToBytes(id int32) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint32(bz, cast.ToUint32(id))
	return bz
}

// Int32FromBytes converts some fixed length bytes back into a int32.
func Int32FromBytes(bz []byte) int32 {
	return cast.ToInt32(binary.BigEndian.Uint32(bz))
}

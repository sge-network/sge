package utils

import "encoding/pem"

func NewPubKeyMemory(bs []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: bs})
}

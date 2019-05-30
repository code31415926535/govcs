package engine

import (
	"crypto/md5"
	"fmt"
)

func dataHash(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

package tools

import (
	"crypto/sha256"
	"fmt"
)

func FindStringInSlice(slice []string, elem string) (res bool) {
	for i := range slice {
		if slice[i] == elem {
			res = true
			return
		}
	}
	return
}

func Hash(s string) string {
	data := []byte(s)
	r := sha256.Sum256(data)
	return fmt.Sprintf("%x", r)
}

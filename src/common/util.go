package common

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
)

// ComputeMD5 computes MD5 from passed data
func ComputeMD5(data []byte) (string, error) {

	if len(data) == 0 {
		return "", errors.New("no data passed to compute MD5")
	}

	hasher := md5.New()

	_, err := hasher.Write(data)
	if err != nil {
		return "", fmt.Errorf("error while computing MD5: %s", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

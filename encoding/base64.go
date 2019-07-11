package encoding

import (
	"encoding/base64"

	"github.com/pkg/errors"
)

// Base64Decode ...
func Base64Decode(encoded string) (raw []byte, err error) {
	raw, err = base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		err = errors.Wrap(err, "unable to base64 decode the value provided")
	}

	return
}

// Base64Encode ...
func Base64Encode(raw []byte) string {
	return base64.StdEncoding.EncodeToString(raw)
}

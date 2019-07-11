package encoding

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

// Base16Encode ...
func Base16Encode(raw []byte) (encoded []byte) {
	_ = hex.Encode(encoded, raw)

	return
}

// Base16Decode ...
func Base16Decode(encoded string) (raw []byte, err error) {
	raw, err = hex.DecodeString(encoded)

	if err != nil {
		err = errors.Wrap(err, "unable to base16 decode the value provided")
	}

	return
}

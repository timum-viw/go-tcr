package qr

import (
    "github.com/skip2/go-qrcode"
    "timum-viw/go-tcr/hash"
    "strconv"
    "fmt"
    "errors"
)

func Generate(input string) ([]byte, error) {
	id64, err := strconv.ParseUint(input, 10, 32)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("input not an unsigned integer (%v)", err))
    }
    encodedString := hash.Generate(uint32(id64))
    // fmt.Println("Encoded String:", encodedString)

    png, err := qrcode.Encode(encodedString, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return png, nil
}
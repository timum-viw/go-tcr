package qr

import (
    "github.com/skip2/go-qrcode"
    "timum-viw/supercoop/hash"
    "strconv"
    "fmt"
    "errors"
)

func Generate(input string) (string, error) {
	id64, err := strconv.ParseUint(input, 10, 32)
    if err != nil {
        return "", errors.New(fmt.Sprintf("input not an unsigned integer (%v)", err))
    }
    encodedString := hash.Generate(uint32(id64))
    // fmt.Println("Encoded String:", encodedString)

    filename := fmt.Sprintf("qrcode_%s.png", input)
    err = qrcode.WriteFile(encodedString, qrcode.Medium, 256, filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}
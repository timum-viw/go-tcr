package main

import (
 "fmt"
 "log"
 "github.com/skip2/go-qrcode"
 "timum-viw/supercoop/hash"
)

func main() {
    encodedString := hash.Generate(12345)
    fmt.Println("Encoded String:", encodedString)

    filename := "qrcode.png"
    err := qrcode.WriteFile(encodedString, qrcode.Medium, 256, filename)
    if err != nil {
        log.Fatalf("Failed to generate QR code: %v", err)
    }

    fmt.Println("QR Code saved as", filename)
}
package main

import (
 "fmt"
 "log"
 "github.com/skip2/go-qrcode"
 "timum-viw/supercoop/hash"
)

// generateQRCode generates a QR code from the given string and saves it as an image file
func generateQRCode(data, filename string) error {
    err := qrcode.WriteFile(data, qrcode.Medium, 256, filename)
    if err != nil {
    return err
    }
    return nil
}

func main() {
    // Generate the encoded string for an integer (example: 12345)
    encodedString := hash.Generate(12345)
    fmt.Println("Encoded String:", encodedString)

    // Generate QR code
    filename := "qrcode.png"
    err := generateQRCode(encodedString, filename)
    if err != nil {
        log.Fatalf("Failed to generate QR code: %v", err)
    }

    fmt.Println("QR Code saved as", filename)
}
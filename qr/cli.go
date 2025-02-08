package main

import (
    "fmt"
    "log"
    "github.com/skip2/go-qrcode"
    "timum-viw/supercoop/hash"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("no input id given")
    }
    id64, err := strconv.ParseUint(os.Args[1], 10, 32)
    if err != nil {
        log.Fatalf("input not an unsigned integer (%v)", err)
    }
    encodedString := hash.Generate(uint32(id64))
    fmt.Println("Encoded String:", encodedString)

    filename := "qrcode.png"
    err = qrcode.WriteFile(encodedString, qrcode.Medium, 256, filename)
    if err != nil {
        log.Fatalf("Failed to generate QR code: %v", err)
    }

    fmt.Println("QR Code saved as", filename)
}
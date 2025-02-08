package main

import (
    "log"
    "os"
    "timum-viw/supercoop/qr"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("no input id given")
    }
    filename, err := qr.Generate(os.Args[1])
    if err != nil {
        log.Fatalf("Failed to generate QR code: %v", err)
    }

    log.Println("QR Code saved as", filename)
}
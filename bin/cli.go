package main

import (
    "log"
    "os"
    "fmt"
    "timum-viw/supercoop/qr"
    "io/ioutil"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("no input id given")
    }
    png, err := qr.Generate(os.Args[1])
    if err != nil {
        log.Fatalf("Failed to generate QR code: %v", err)
    }

    filename := fmt.Sprintf("qrcode_%s.png", os.Args[1])
    err = ioutil.WriteFile(filename, png, 0664)
	if err != nil {
		log.Fatalf("Failed to write qr code to file (%v)", err)
	}
    log.Println("QR Code saved as", filename)
}
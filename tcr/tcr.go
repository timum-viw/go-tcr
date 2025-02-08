package tcr

import (
	"log"
	i2c "github.com/d2r2/go-i2c"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"strings"
	"errors"
)


func Read() {
    i2c, err := i2c.NewI2C(0x0c, 1)
    if err != nil { log.Fatal(err) }
    // Free I2C connection on exit
    defer i2c.Close()
    var buf = make([]byte, 256)
    for {
		time.Sleep(2 * time.Second)
			_,err := i2c.ReadBytes(buf)
		if(err != nil) {
			log.Fatal(err)
		}
		// log.Printf("read %d bytes", n)
		length := binary.LittleEndian.Uint16(buf[:2])
		if length < 1 {
			continue
		}
		if length > 254 {
			log.Printf("invalid content length %d\n", length)
			continue
		}
		log.Printf("content len: %d\n", length)
		str := string(buf[2:2+length])
		log.Println(len(str))
		log.Println(str)
		n, err := validateAndExtractInteger(str)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("decoded: %d \n", n)
		}
    }
}


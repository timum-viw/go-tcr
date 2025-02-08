package tcr

import (
	"errors"
	"os"
	"fmt"
	"syscall"
	"encoding/binary"
)

const (
	I2C_SLAVE = 0x0703
	TCR_ADDR = 0x0c
)

type TCR struct {
	f *os.File
	buf []byte
}

func Open(bus int) (*TCR, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), I2C_SLAVE, uintptr(TCR_ADDR))
	if errno != 0 {
		f.Close()
		return nil, errors.New("syscall failed")
	}
    var buf = make([]byte, 256)
	v := &TCR{f: f, buf: buf}
	return v, nil
}

func (v *TCR) Close() error {
	return v.f.Close()
}

func (v *TCR) Read() (string, error) {
	n, err := v.f.Read(v.buf)
	if err != nil {
		return "", err
	}
	if n < 1 {
		return "", nil
	}

	length := binary.LittleEndian.Uint16(v.buf[:2])
	if length < 1 {
		return "", nil
	}
	if length > 254 {
		return "", errors.New("invalid content length returned")
	}
	
	str := string(v.buf[2:2+length])
	return str, nil
}


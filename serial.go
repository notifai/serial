/*
Goserial is a simple go package to allow you to read and write from
the serial port as a stream of bytes.

It aims to have the same API on all platforms, including windows.  As
an added bonus, the windows package does not use cgo, so you can cross
compile for windows from another platform.  Unfortunately goinstall
does not currently let you cross compile so you will have to do it
manually:

 GOOS=windows make clean install

Currently there is very little in the way of configurability.  You can
set the baud rate.  Then you can Read(), Write(), or Close() the
connection.  Read() will block until at least one byte is returned.
Write is the same.  There is currently no exposed way to set the
timeouts, though patches are welcome.

Currently all ports are opened with 8 data bits, 1 stop bit, no
parity, no hardware flow control, and no software flow control.  This
works fine for many real devices and many faux serial devices
including usb-to-serial converters and bluetooth serial ports.

You may Read() and Write() simulantiously on the same connection (from
different goroutines).

Example usage:

  package main

  import (
        "github.com/tarm/serial"
        "log"
  )

  func main() {
        c := &serial.Config{Name: "COM5", Baud: 115200}
        s, err := serial.OpenPort(c)
        if err != nil {
                log.Fatal(err)
        }

        n, err := s.Write([]byte("test"))
        if err != nil {
                log.Fatal(err)
        }

        buf := make([]byte, 128)
        n, err = s.Read(buf)
        if err != nil {
                log.Fatal(err)
        }
        log.Print("%q", buf[:n])
  }
*/
package serial

import (
	"errors"
	"io"
	"time"
)

type Port interface {
	io.ReadWriteCloser
	SetReadDeadline(time.Duration) error
	Flush() error
	Status() (uint, error)
	SetDTR(bool) error
	SetRTS(bool) error
	SetParity(Parity) error
}

var ErrNotSupported = errors.New("serial: not supported")

// ErrBadSize is returned if Size is not supported.
var ErrBadSize = errors.New("serial: unsupported serial data size")

// ErrBadStopBits is returned if the specified StopBits setting not supported.
var ErrBadStopBits = errors.New("serial: unsupported stop bit setting")

// ErrBadParity is returned if the parity is not supported.
var ErrBadParity = errors.New("serial: unsupported parity setting")

var ErrInvalidArg = errors.New("serial: invalid argument")

// OpenPort opens a serial port with the specified configuration
func OpenPort(c Config) (Port, error) {
	if c.Size == 0 {
		c.Size = DefaultSize
	}

	if c.Parity == 0 {
		c.Parity = ParityNone
	}

	if c.StopBits == 0 {
		c.StopBits = Stop1
	}

	c.timeout = MaxTimeout

	return openPort(c)
}

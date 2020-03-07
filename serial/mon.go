/*
 * This application can be used to experiment and test various serial port options
 */

package main

import (
	"io"
	"fmt"
	"sync"
 
	"github.com/jacobsa/go-serial/serial"
)
 
func Open(device string) (io.ReadWriteCloser) {
    options := serial.OpenOptions{
		PortName: device,
		BaudRate: 12800,
		DataBits: 8,
		StopBits: 1,
		MinimumReadSize: 4,
	}
	port, err := serial.Open(options)
	if err != nil {
		fmt.Println("serial.Open: %v", err)
	}
	return port
}

func Run(in io.ReadWriteCloser, out io.ReadWriteCloser) {
	for {
		buf := make([]byte, 1024)
		n, err := in.Read(buf)
		
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]
			for _, b := range buf {
				fmt.Println(b)
			} 
			out.Write(buf)
		}
	}
}

 func main() {
	var wait sync.WaitGroup
  	wait.Add(2)

	in  := Open("COM2")
	out := Open("COM1")

	//p = escpos.New(out)
	go Run(in, out)
	go Run(out, in)
	wait.Wait();
 }
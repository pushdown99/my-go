/*
 * This application can be used to experiment and test various serial port options
 */

package main

import (
	"io"
	"fmt"
	"sync"
	"time"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"io/ioutil"
 
	"github.com/jacobsa/go-serial/serial"
)
type JsonData struct {
	Data 		string
	Timestamp 	int64
}

func Open(device string) (io.ReadWriteCloser) {
    options := serial.OpenOptions{
		PortName: device,
		BaudRate: 19200,
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
		buf := make([]byte, 4096)
		n, err := in.Read(buf)
		
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf  = buf[:n]
			if n > 0 { 
				fmt.Println(hex.Dump(buf)) 
				d := JsonData { Data: hex.EncodeToString(buf), Timestamp: time.Now().Unix() }
				b, _ := json.Marshal(d)
				resp, err := http.Post("http://localhost:3000", "application/json", bytes.NewBuffer(b))
				if err != nil { 
					fmt.Println("Error: ", err) 
				} else { 
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						s, err := ioutil.ReadAll(resp.Body)
						if err == nil {
							b, _ = hex.DecodeString(string(s))
							fmt.Println(b) 
							out.Write(b)
						} else {
							fmt.Println("Error: ", err) 
						}
					}
				}
			}
			//out.Write(buf)
		}
	}
}

 func main() {
	var wait sync.WaitGroup
  	wait.Add(2)

	in  := Open("COM2")
	out := Open("COM1")

	go Run(in, out)
	go Run(out, in)
	wait.Wait();
 }
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const	Input string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const	Output string =  "NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm"

func rot13(b byte) byte {
	i := strings.Index(Input, string(b))
	if i < 0 {
		return b
	}
	return Output[i]
}

type rot13Reader struct {
	r io.Reader
}

func (orgReader rot13Reader) Read(buffer []byte) (bytesRead int, err error) {
	mybuffer := make([]byte, 8)
	read, err := orgReader.r.Read(mybuffer)
	if err == io.EOF {
		// fmt.Println("<EOF>")
		return 0, io.EOF
	}
	if err != nil {
		fmt.Println("oops we should panic here")
		return -1, io.EOF
	}
	for i, c := range mybuffer {
		buffer[i] = rot13(c)
	}
	return read, nil
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr! 42")
	r := rot13Reader{s}
	fmt.Printf("Try1: ")
	io.Copy(os.Stdout, &r)

	s2 := strings.NewReader("")
	r = rot13Reader{s2}
	fmt.Printf("\nTry2: ")
	io.Copy(os.Stdout, &r)
}

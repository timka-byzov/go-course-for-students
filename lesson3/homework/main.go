package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Options struct {
	From        string
	To          string
	Offset      int
	Limit       int
	BlockSize   int
	Convesrions []string
}

const (
	DBlockSize int = 2
)

const (
	EndLine        byte = 0xa
	CarriegeReturn byte = 0xd
	Empty          byte = 0x0
)

// type BlockReader interface{
// 	ReadBlock()
// }

type Reader interface {
	Read(bcount int, oofset int) ([]byte, error)
}

type Writer interface {
	Write([]byte)
}

type IOReader struct {
	BlockSize int
	Offset    int
	Limit     int
}

func (rd IOReader) Read(limit int, offset int) ([]byte, error) {
	var err error

	buffer := []byte{}
	bytesReaded := 0
	endLine := false

	for (bytesReaded < limit+offset || limit == -1) && !endLine {
		chunkBuffer := make([]byte, rd.BlockSize)
		n, err := os.Stdin.Read(chunkBuffer)
		if err != nil {
			break
		}

		bytesReaded += n

		for i := 0; i < n; i++ {
			if chunkBuffer[i] == EndLine || chunkBuffer[i] == Empty || chunkBuffer[i] == CarriegeReturn {
				endLine = true
				break
			}
			buffer = append(buffer, chunkBuffer[i])
		}
	}

	buffer = buffer[offset:]
	return buffer, err
}

// type FileReader struct {
// 	IOReader
// 	ifName string
// 	ofname string
// }

// func (fr FileReader) OffsetRead(offset int) {

// }

func ParseFlags() (*Options, error) {
	var opts Options
	var conversions string

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "Offset")
	flag.IntVar(&opts.Limit, "limit", -1, "Limit")
	flag.IntVar(&opts.BlockSize, "block-size", DBlockSize, "Block-size")
	flag.StringVar(&conversions, "conv", "", "Conversions")

	// todo: parse and validate all flags

	flag.Parse()
	opts.Convesrions = strings.Split(conversions, ", ")

	return &opts, nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	fmt.Println(opts, err)

	var IOReader Reader = IOReader{BlockSize: opts.BlockSize, Offset: opts.Offset, Limit: opts.Limit}
	buffer, err := IOReader.Read(5, 10)
	fmt.Println(string(buffer), err)

	// todo: implement the functional requirements described in read.me
}

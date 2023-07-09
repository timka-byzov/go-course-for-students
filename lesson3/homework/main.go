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

// type BlockReader interface{
// 	ReadBlock()
// }

type Reader interface {
	LimitRead(bcount int) ([]byte, error)
	UnlimitRead() ([]byte, error)
	OffsetRead(bcount int) error
}

type Writer interface {
	Write([]byte)
}

type IOReader struct {
	BlockSize int
	Offset    int
	Limit     int
}

func (rd IOReader) OffsetRead(bcount int) error {
	var err error

	for {
		n, err := os.Stdin.Read(make([]byte, rd.BlockSize))
		if err != nil {
			break
		}

		if n < rd.BlockSize {
			break
		}

	}
	// fmt.Println(err)
	return err
}

func (rd IOReader) LimitRead(bcount int) ([]byte, error) {
	var err error

	buffer := []byte{}
	for {
		chunkBuffer := make([]byte, rd.BlockSize)
		n, err := os.Stdin.Read(chunkBuffer)

		buffer = append(buffer, chunkBuffer...)

		if err != nil {
			break
		}

		if n < rd.BlockSize {
			break
		}

	}

	return buffer, err

}

func (rd IOReader) UnlimitRead() ([]byte, error) {
	return []byte{}, nil
}

type FileReader struct{}

func ParseFlags() (*Options, error) {
	var opts Options
	var conversions string

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "Offset")
	flag.IntVar(&opts.Limit, "limit", -1, "Limit")
	flag.IntVar(&opts.BlockSize, "block-size", 0, "Block-size")
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

	IOReader.OffsetRead(2)

	// todo: implement the functional requirements described in read.me
}

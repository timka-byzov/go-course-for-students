package main

import (
	"flag"
	"fmt"
	"io"
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
	Read() ([]byte, error)
}

type Writer interface {
	Write([]byte) error
}

type IOReader struct {
	opts *Options
	data []byte
}

type IOWriter struct {
	opts *Options
}

func (rd *IOReader) Read() ([]byte, error) {

	buffer := []byte{}
	bytesReaded := 0

	for rd.opts.Limit == -1 || bytesReaded < rd.opts.Limit+rd.opts.Offset {
		chunkBuffer := make([]byte, rd.opts.BlockSize)
		n, readErr := os.Stdin.Read(chunkBuffer)
		if readErr != nil && readErr != io.EOF {
			return nil, fmt.Errorf("ошибка чтения из консоли %w", readErr)
		}

		bytesReaded += n
		buffer = append(buffer, chunkBuffer[:n]...)

		if n < rd.opts.BlockSize {
			break
		}
	}

	buffer = buffer[rd.opts.Offset:]
	rd.data = buffer
	return buffer, nil
}

func (wr *IOWriter) Write(data []byte) error {
	_, err := os.Stdout.Write(data)

	if err != nil {
		return fmt.Errorf("ошибка при записи в файл %w", err)
	}
	return err

}

type FileReader struct {
	opts *Options
	data []byte
}

type FileWriter struct {
	opts *Options
}

func (fr *FileReader) Read() ([]byte, error) {

	exPath, _ := os.Getwd()
	fileName := exPath + "\\" + fr.opts.From

	inputFile, openErr := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if openErr != nil {
		return nil, fmt.Errorf("ошибка при открытии файла ввода %w", openErr)
	}

	defer inputFile.Close()

	buffer := []byte{}
	bytesReaded := 0

	for fr.opts.Limit == -1 || bytesReaded < fr.opts.Limit+fr.opts.Offset {

		chunkBuffer := make([]byte, fr.opts.BlockSize)

		n, readErr := inputFile.Read(chunkBuffer)
		if readErr != nil && readErr != io.EOF {
			return nil, fmt.Errorf("ошибка чтения из файла ввода %w", readErr)
		}

		bytesReaded += n
		for i := 0; i < n; i++ {
			buffer = append(buffer, chunkBuffer[i])
		} //TODO: заменить на срез

		if readErr == io.EOF {
			break
		}

	}
	buffer = buffer[fr.opts.Offset:]
	fr.data = buffer
	return buffer, nil
}

func (fw *FileWriter) Write(data []byte) error {
	filePath, _ := os.Getwd()
	fileName := filePath + "\\" + fw.opts.To

	file, openErr := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL, 0755)
	if openErr != nil {
		return fmt.Errorf("ошибка при открытии файла вывода %w", openErr)
	}
	defer file.Close()

	_, writeErr := file.Write(data)
	if writeErr != nil {
		return fmt.Errorf("ошибка при записи в файл вывода %w", writeErr)
	}

	return nil
}

func ParseFlags() (*Options, error) {
	var opts Options
	var conversions string

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "Offset")
	flag.IntVar(&opts.Limit, "limit", -1, "Limit")
	flag.IntVar(&opts.BlockSize, "block-size", DBlockSize, "Block-size")
	flag.StringVar(&conversions, "conv", "", "Conversions")

	// TODO: parse and validate all flags

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

	// fmt.Println(opts, err)

	// var IOReader Reader = IOReader{opts}
	// buffer, err := IOReader.Read(5, 10)
	// fmt.Println(string(buffer), err)

	// var fileReader Reader = FileReader{opts}
	// buffer, err := fileReader.Read(2, 0)
	// fmt.Println(string(buffer), err)

	var reader Reader
	if opts.From == "" {
		reader = &IOReader{opts: opts}
	} else {
		reader = &FileReader{opts: opts}
	}

	data, readErr := reader.Read()
	if readErr != nil {
		panic(readErr)
	}

	var writer Writer
	if opts.To == "" {
		writer = &IOWriter{opts: opts}
	} else {
		writer = &FileWriter{opts: opts}
	}

	writeErr := writer.Write(data)
	if writeErr != nil {
		panic(writeErr)
	}

}

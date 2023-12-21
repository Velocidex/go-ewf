package main

import (
	"io"
	"os"

	"github.com/Velocidex/go-ewf/parser"
	"github.com/alecthomas/kingpin"
)

var (
	cat_command = app.Command("cat", "Dump file contents")

	cat_command_files = cat_command.Arg(
		"files", "The image file to inspect",
	).Required().Strings()

	cat_command_out_file = cat_command.Flag(
		"output", "The file to write",
	).String()
)

func doCat() {
	var err error

	output := os.Stdout
	if *cat_command_out_file != "" {
		output, err = os.OpenFile(*cat_command_out_file,
			os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
		kingpin.FatalIfError(err, "Creating output file.")

		defer output.Close()
	}

	options := &parser.EWFOptions{
		LRUSize: 100,
	}

	files := []io.ReaderAt{}
	for _, filename := range *cat_command_files {
		fd, err := os.Open(filename)
		kingpin.FatalIfError(err, "Unable to open EWF File")
		files = append(files, fd)
	}

	volume, err := parser.OpenEWFFile(options, files...)
	kingpin.FatalIfError(err, "Unable to parse EWF File")

	buff := make([]byte, 1000)
	total_size := int(volume.ChunkSize * volume.NumberOfChunks)

	for i := 0; i < total_size; i += len(buff) {
		n, err := volume.ReadAt(buff, int64(i))
		kingpin.FatalIfError(err, "Unable to parse EWF File")
		output.Write(buff[:n])
	}
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case cat_command.FullCommand():
			doCat()

		default:
			return false
		}
		return true
	})
}

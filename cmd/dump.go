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
		"files", "The image files to inspect (should include all .E* files)",
	).Required().Strings()

	cat_command_skip = cat_command.Flag(
		"skip", "Bytes to skip").Int64()

	cat_command_count = cat_command.Flag(
		"count", "Total number of bytes to dump (0 mean to the end)").Int64()

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
	total_size := int64(volume.ChunkSize * volume.NumberOfChunks)

	if *cat_command_count > 0 {
		total_size = *cat_command_skip + *cat_command_count
	}

	for i := *cat_command_skip; i < total_size; i += int64(len(buff)) {
		n, err := volume.ReadAt(buff, i)
		kingpin.FatalIfError(err, "Unable to parse EWF File")
		if i+int64(n) > total_size {
			n = int(total_size - i)
		}

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

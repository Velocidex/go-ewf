package main

import (
	"io"
	"os"

	"github.com/Velocidex/go-ewf/parser"
	"github.com/alecthomas/kingpin"
)

var (
	parse_command = app.Command("parse", "Parse a file")

	parse_command_files = parse_command.Arg(
		"files", "The image files to inspect",
	).Required().Strings()
)

func doParse() {
	options := &parser.EWFOptions{
		LRUSize: 100,
	}

	files := []io.ReaderAt{}
	for _, filename := range *parse_command_files {
		fd, err := os.Open(filename)
		kingpin.FatalIfError(err, "Unable to open EWF File")
		files = append(files, fd)
	}

	volume, err := parser.OpenEWFFile(options, files...)
	kingpin.FatalIfError(err, "Unable to parse EWF File")

	volume.WriteDebug(os.Stdout)
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case parse_command.FullCommand():
			doParse()
		default:
			return false
		}
		return true
	})
}

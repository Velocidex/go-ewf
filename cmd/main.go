package main

import (
	"os"

	"github.com/Velocidex/go-ewf/parser"
	kingpin "github.com/alecthomas/kingpin/v2"
)

type CommandHandler func(command string) bool

var (
	app = kingpin.New("go-ewf",
		"A tool for inspecting ewf files.")

	verbose_flag = app.Flag(
		"verbose", "Show verbose information").Bool()

	command_handlers []CommandHandler
)

func main() {
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	if *verbose_flag {
		parser.SetDebug()
	}

	for _, command_handler := range command_handlers {
		if command_handler(command) {
			break
		}
	}
}

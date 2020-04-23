package main

import (
	"errors"
	"fmt"
	"github.com/cristianrz/aptlist/packages"
	"github.com/cristianrz/opts"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var program = path.Base(os.Args[0])

func main() {
	var (
		lines []string

		args = opts.Opts{
			'l': false,
			'f': "",
			'h': false,
			'v': false,
		}
	)

	err := args.Parse()
	if err != nil {
		s := fmt.Sprintf("%s: %s\n%s\n", program, err, usage())
		fatal(errors.New(s))
	}

	var (
		debug  = args['v'].(bool)
		filter = args['f'].(string)
		human  = args['h'].(bool)
		disk   = !args['l'].(bool)
	)

	if !debug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	lines, err = packages.Sizes(filter, disk, human)
	if err != nil {
		fatal(err)
	}

	for _, line := range lines {
		fields := strings.Split(line, ",")

		if disk && human {
			line = fields[2] + "\t" + fields[1]
		} else if disk {

			for fields[0][0] == '0' {
				if len(fields[0]) == 0 {
					fields[0] = "0"
					break
				}

				fields[0] = fields[0][1:]
			}

			line = fields[0] + "\t" + fields[1]

		} else {
			line = fields[0]
		}

		fmt.Println(line)
	}
}

func fatal(err error) {
	_, err = fmt.Fprintf(os.Stderr, "%s: %s\n", program, err)
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

func usage() string {
	s := fmt.Sprintf(`
Aptlist is a tool to show installed packages.
	
usage: %s [-dht] [-f pattern]

	-f  filter pattern
	-h  human readable size
	-l  don't show disk usage
	-r  reverse order

`, program)
	return s
}

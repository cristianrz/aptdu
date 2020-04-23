package main

import (
	"fmt"
	"github.com/cristianrz/aptlist/packages"
	"github.com/cristianrz/aptlist/print"
	"github.com/cristianrz/opts"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var program = path.Base(os.Args[0])

func main() {

	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	//f, err := os.Create("./aptlist.log")
	//if err != nil {
	//	panic(err)
	//}
	//log.SetOutput(f)

	args := opts.Opts{
		'd': false,
		'f': "",
		'h': false,
	}

	err := args.Parse()

	disk := args['d'].(bool)
	human := args['h'].(bool)

	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%s: %s\n", program, err)
		if err != nil {
			panic(err)
		}
		usage()
		os.Exit(1)
	}

	if disk {
		lines, err := packages.Sizes(human)
		if err != nil {
			fatal(err)
		}

		for _, line := range lines {
			for line[0] == '0' {
				line = line[1:]
			}

			fmt.Println(line)
		}
	} else {
		err := print.List()
		if err != nil {
			fatal(err)
		}
	}

	//err = f.Close()
	//if err != nil {
	//	panic(err)
	//}
}

func fatal(err error) {
	_, err = fmt.Fprintf(os.Stderr, "%s: %s\n", program, err)
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

func
usage() {
	_, err := fmt.Fprintf(os.Stderr, `
Aptlist is a tool to show installed packages.
	
usage: %s [-dht] [-f pattern]

	-d  shows disk usage
	-f  filter pattern
	-h  human readable size
	-r  reverse order

`, program)
	if err != nil {
		panic(err)
	}
}

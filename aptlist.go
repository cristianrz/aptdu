package main

import (
	"bufio"
	"fmt"
	"github.com/cristianrz/opts"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	program := path.Base(os.Args[0])

	f, err := os.Create("./aptlist.log")
	exitIf(program, err)

	log.SetOutput(f)

	args := make(opts.Opts)

	args['d'] = false
	args['f'] = ""
	args['h'] = false
	args['t'] = false

	err = args.Parse()
	exitIf(program, err)

	for k, v := range args {
		log.Println(string(k), v)
	}

	if args['d'] == false {
		err = noDisk()
		exitIf(program, err)
	} else {

		m, err := disk()
		exitIf(program, err)

		b := args['h'].(bool)
		err = printSorted(m, b)
		exitIf(program, err)
	}

	err = f.Close()
	exitIf(program, err)

}

func exitIf(p string, e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", p, e)
		os.Exit(1)
	}
}

func usage(p string) {
	fmt.Fprintf(os.Stderr, `Aptlist is a tool to show installed packages.
	
	usage: %s [-dht] [-f pattern]
	
		-d  shows disk usage
		-h  human readable size
		-t  reverse order
		-f  filter pattern
	
	`, p)
	os.Exit(1)
}

func noDisk() error {
	cmd := exec.Command("dpkg", "--get-selections")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := strings.Fields(scanner.Text())
		fmt.Println(m[0])
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}

//func usage() {
//
//}

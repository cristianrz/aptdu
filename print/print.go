package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/cristianrz/opts"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
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
	fmt.Fprintf(os.Stderr, "usage: %s [-dht] [-f pattern]\n", p)
	os.Exit(1)
}

func disk() (map[float64]string, error) {
	m := make(map[float64]string)

	cmd := exec.Command("dpkg-query", "-Wf", "${Installed-Size}\t${Package}\n")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		log.Println("converting", fields)

		key, err := strconv.ParseFloat(fields[0], 32)
		if err != nil {
			log.Println("failed converting '" + fields[0] + "' to int")
			continue
		}

		value := fields[1]
		m[key] = value
	}

	if err = cmd.Wait(); err != nil {
		return nil, err
	}

	return m, nil
}

func printSorted(m map[float64]string, human bool) error {
	keys := make([]float64, 0, len(m))
	for key, _ := range m {
		log.Println("appending", key)
		keys = append(keys, key)
	}

	sort.Float64s(keys)

	if human {
		for _, v := range keys {
			h, err := ToHuman(v)
			if err != nil {
				return err
			}
			fmt.Println(h, "\t", m[v])
		}
	} else {
		for _, v := range keys {
			fmt.Println(v, "\t", m[v])
		}
	}

	return nil
}

func ToHuman(f float64) (string, error) {
	var (
		c     = 0
		size  string
		units rune
	)

	log.Println("starting with c=", c, "and f=", f)

	for f > 1024 {
		log.Println(f, "is bigger than 1024")
		f = f / 1024
		c++
		log.Println("counter is now", c)
	}

	log.Println("exiting with c=", c, "and f=", f)

	switch c {
	case 0:
		units = 'K'
	case 1:
		units = 'M'
	case 2:
		units = 'G'
	case 3:
		units = 'T'
	default:
		return "", errors.New("found too high size")
	}

	s := fmt.Sprintf("%.0f", f) // s == "123.456000"
	size = s + string(units)
	return size, nil
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
//	basename := path.Base(os.Args[0])
//	fmt.Printf(`Aptlist is a tool to show installed packages.
//
//usage: %s [-dht] [-f pattern]
//
//	-d  shows disk usage
//	-h  human readable size
//	-t  reverse order
//	-f  filter pattern
//
//`, basename)
//
//}
